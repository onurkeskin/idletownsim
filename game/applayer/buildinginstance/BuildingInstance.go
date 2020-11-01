package buildinginstance

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	buildingdomain "github.com/app/game/applayer/building/domain"
	buildinginstancedomain "github.com/app/game/applayer/buildinginstance/domain"
	effectdomain "github.com/app/game/applayer/effect/domain"
	gamemapdomain "github.com/app/game/applayer/gamemap/domain"
	Generaldomain "github.com/app/game/applayer/general/domain"
	Generatordomain "github.com/app/game/applayer/generator/domain"

	General "github.com/app/game/applayer/general"
	Generator "github.com/app/game/applayer/generator"
)

type BuildingInstance struct {
	*General.ObjectProperties `json:"props" bson:"-"`
	UniqueId                  string                   `json:"-" bson:"_id"`
	ParentId                  string                   `json:"-" bson:"ParentId"`
	Level                     int                      `json:"bLevel" bson:"bLevel"`
	BaseProducts              Generaldomain.IBProducts `json:"bproduct" bson:"-"`
	ProductionIntervalNano    int64                    `json:"bproductiontime" bson:"-"`
	LastProductionTime        time.Time                `json:"prodtime" bson:"prodtime"`
	BuiltTime                 time.Time                `json:"builttime" bson:"builttime"`
	// to change the modifiers from game to instnace

	SpaceEffect         gamemapdomain.ISpaceEffect                         `json:"spaceeffect" bson:"-"`
	DeployedEffects     []effectdomain.IEffect                             `json:"-" bson:"-"`
	BIPChangedListeners []buildinginstancedomain.BIPropertyChangedListener `json:"-" bson:"-"`
	Gen                 Generatordomain.Generator                          `json:"-" bson:"-"`
}

func (d *BuildingInstance) MarshalJSON() ([]byte, error) {
	type toJsonBI BuildingInstance
	return json.Marshal(&struct {
		*toJsonBI
		LastProductionTime string `json:"prodtime"`
		BuiltTime          string `json:"builttime"`
	}{
		toJsonBI:           (*toJsonBI)(d),
		LastProductionTime: d.LastProductionTime.Format("2006-01-02T15:04:05.999999Z"),
		BuiltTime:          d.BuiltTime.Format("2006-01-02T15:04:05.999999Z"),
	})
}

func NewBuildingInstance(
	UniqueId string,
	ParentId string,
	Level int,
	BaseProducts Generaldomain.IBProducts,
	ProductionIntervalNano int64,
	BuiltTime time.Time,
	SpaceEffect gamemapdomain.ISpaceEffect,
	DeployedEffects []effectdomain.IEffect,
) buildinginstancedomain.IBuildingInstance {
	toRet := BuildingInstance{
		ObjectProperties:       General.NewObjectProperties(),
		UniqueId:               UniqueId,
		ParentId:               ParentId,
		Level:                  Level,
		BaseProducts:           BaseProducts,
		BuiltTime:              BuiltTime,
		ProductionIntervalNano: ProductionIntervalNano,
		SpaceEffect:            SpaceEffect,
		DeployedEffects:        DeployedEffects,
		Gen:                    Generator.NewBasicGenerator(),
		BIPChangedListeners:    []buildinginstancedomain.BIPropertyChangedListener{}}
	toRet.LastProductionTime = BuiltTime

	return &toRet
}

type FormFromBuildingOptions struct {
	ID            string
	Buildinglevel int
}

func FormFromBuilding(b buildingdomain.IBuilding, options FormFromBuildingOptions) buildinginstancedomain.IBuildingInstance {
	var bIProds Generaldomain.IBProducts
	var bIEff gamemapdomain.ISpaceEffect
	bEff := b.GetSpaceEffect()
	if bEff != nil {
		bIEff = bEff.Clone().(gamemapdomain.ISpaceEffect)
	}
	bProds := b.GetBaseProducts()
	if bProds != nil {
		bIProds = bProds.Clone().(Generaldomain.IBProducts)
	}
	interval := b.GetProductionIntervalNano()

	toRet := NewBuildingInstance(options.ID, b.GetID(), options.Buildinglevel, bIProds, interval, time.Time{}, bIEff, []effectdomain.IEffect{})
	if toRet.GetSpaceEffect() != nil {
		toRet.GetSpaceEffect().SetIssuer(toRet)
	}
	for _, v := range b.PropertiesSlice() {
		toRet.AddProperty(v.GetParamName(), v.GetParamValue())
	}
	return toRet
}

func (b *BuildingInstance) GetUniqueID() string {
	return b.UniqueId
}

func (b *BuildingInstance) GetParentID() string {
	return b.ParentId
}

func (b *BuildingInstance) GetLevel() int {
	return b.Level
}

func (b *BuildingInstance) GetGenerator() Generatordomain.Generator {
	return b.Gen
}

func (b *BuildingInstance) GetBaseProducts() Generaldomain.IBProducts {
	return b.BaseProducts
}

func (b *BuildingInstance) PredictWorkTime(t int64) int64 {
	toProduce := t / b.ProductionIntervalNano
	return toProduce
}

func (b *BuildingInstance) DoWork(d time.Duration) []Generaldomain.IBProducts {
	t := d.Nanoseconds()
	toRet := []Generaldomain.IBProducts{}
	toProduce := t / b.ProductionIntervalNano
	leftOver := t % b.ProductionIntervalNano
	//fmt.Printf("toProduce:%v\n", toProduce)
	//fmt.Printf("t:%v\n", t)
	//fmt.Printf("interval:%v\n\n", b.ProductionIntervalNano)
	b.LastProductionTime = b.LastProductionTime.Add(time.Unix(0, t).Sub(time.Unix(0, leftOver)))
	var i int64 = 0
	for i < toProduce {
		toRet = append(toRet, b.GetGenerator().Generate(b.GetBaseProducts()))
		i++
	}
	return toRet
}

func (b *BuildingInstance) GetExpectedOutcome() Generaldomain.IBProducts {
	total := Generaldomain.IBProducts{}
	total = b.Gen.Generate(total)
	return total
}

func (b *BuildingInstance) AddSpaceEffect(e gamemapdomain.ISpaceEffect) {
	b.SpaceEffect = e

	if b.GetSpaceEffect() != nil {
		b.GetSpaceEffect().SetIssuer(b)
	}
}

func (b *BuildingInstance) GetSpaceEffect() gamemapdomain.ISpaceEffect {
	return b.SpaceEffect
}

func (b *BuildingInstance) ApplySpaceEffect(s gamemapdomain.ISpace) error {
	b.SpaceEffect.ApplyEffect(s)
	b.DeployedEffects = append(b.DeployedEffects, b.SpaceEffect)
	return nil
}

func (b *BuildingInstance) RemoveSpaceEffect() error {
	if b.SpaceEffect == nil {
		return errors.New("Instance does not have space effect")
	}
	tileEff := b.SpaceEffect

	for in, v := range b.DeployedEffects {
		if v == tileEff {
			b.DeployedEffects = append(b.DeployedEffects[:in], b.DeployedEffects[in+1:]...)
		}
	}
	(tileEff).RemoveEffect()
	return nil
}

func (b *BuildingInstance) GetBaseProductionIntervalNano() int64 {
	return b.ProductionIntervalNano
}

func (b *BuildingInstance) GetLastProductionTimeUnix() time.Time {
	return b.LastProductionTime
}

func (b *BuildingInstance) SetLastProductionTimeUnix(a time.Time) {
	b.LastProductionTime = a
	return
}

func (b *BuildingInstance) GetBuiltTimeUnix() time.Time {
	return b.BuiltTime
}

func (b *BuildingInstance) SetBuiltTimeUnix(a time.Time) {
	b.BuiltTime = a
	if a.Unix() > b.LastProductionTime.Unix() {
		b.SetLastProductionTimeUnix(a)
	}
	return
}

func (b *BuildingInstance) GetProductMods() []buildinginstancedomain.IBuildingProductionEffect {
	toReturn := []buildinginstancedomain.IBuildingProductionEffect{}
	gens := b.Gen.GetProductionMods()
	for _, v := range gens {
		toReturn = append(toReturn, v.(buildinginstancedomain.IBuildingProductionEffect))
	}
	return toReturn
}

func (b *BuildingInstance) AddProductMods(eff buildinginstancedomain.IBuildingProductionEffect) {
	b.Gen.AddProductionMod(eff)
}

func (b *BuildingInstance) RemoveProductMods(eff buildinginstancedomain.IBuildingProductionEffect) {
	b.Gen.RemoveProductionMod(eff)
}

func (b *BuildingInstance) ResetProductMods() {
	b.Gen.Reset()
}

//PROPERTYCHANGED
func (b *BuildingInstance) BIPropertyChangedEvent() {
	for _, v := range b.BIPChangedListeners {
		v.OnBIPropertyChange(b)
	}
}

func (b *BuildingInstance) AddBIPropertyChangeListener(listener buildinginstancedomain.BIPropertyChangedListener) error {
	for _, v := range b.BIPChangedListeners {
		if v == listener {
			return errors.New(fmt.Sprintf("Already has listener:%v", listener))
		}
	}
	b.BIPChangedListeners = append(b.BIPChangedListeners, listener)
	return nil
}

func (b *BuildingInstance) RemoveBIPropertyChangeListener(listener buildinginstancedomain.BIPropertyChangedListener) error {
	for in, v := range b.BIPChangedListeners {
		if v == listener {
			b.BIPChangedListeners = append(b.BIPChangedListeners[:in], b.BIPChangedListeners[in+1:]...)
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Listener not found:%v", listener))
}

//IEFFECTISSUER INTERFACE
func (b *BuildingInstance) GetDeployedEffects() []effectdomain.IEffect {
	return b.DeployedEffects
}

func (g *BuildingInstance) String() string {
	str := ""
	str += fmt.Sprintf("Building Of Type:%s\n", g.ParentId)
	str += fmt.Sprintf("Properties:%s\n", g.ObjectProperties)
	str += fmt.Sprintf("Level:%d\n", g.Level)
	str += fmt.Sprintf("Base Products:%s\n", g.BaseProducts)
	str += fmt.Sprintf("Production Interval:%d\n", g.ProductionIntervalNano)
	str += fmt.Sprintf("Last Production Time:%d\n", g.LastProductionTime)
	str += fmt.Sprintf("Generator:%s", g.Gen)
	return str
}
