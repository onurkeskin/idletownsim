package building

import (
	buildingdomain "github.com/app/game/applayer/building/domain"
	gamemapdomain "github.com/app/game/applayer/gamemap/domain"

	general "github.com/app/game/applayer/general"
	generaldomain "github.com/app/game/applayer/general/domain"
)

type Buildings []Building

type Building struct {
	ID                        string `json:"id,omitempty"" bson:"_id,omitempty""`
	*general.ObjectProperties `json:"props" bson:"-"`

	BasePrice    generaldomain.IBProducts `json:"price" bson:"price"`
	UpValScheme  generaldomain.IValScheme `json:"upgradepricescheme" bson:"upgradepricescheme"`
	UpStatScheme generaldomain.IValScheme `json:"upgraderesultscheme" bson:"upgraderesultscheme"`

	Yields         generaldomain.IBProducts `json:"yields" bson:"yields"`
	ProductionTime int64                    `json:"time" bson:"time"`

	SpaceEffect gamemapdomain.ISpaceEffect `json:"spaceeffect" bson:"spaceeffect"`
}

func NewBuilding(
	ID string,
	BasePrice generaldomain.IBProducts,
	upValScheme generaldomain.IValScheme,
	UpStatScheme generaldomain.IValScheme,
	Yields generaldomain.IBProducts,
	ProductionTime int64) buildingdomain.IBuilding {
	toRet := Building{
		ObjectProperties: general.NewObjectProperties(),
		ID:               ID,
		BasePrice:        BasePrice,
		UpValScheme:      upValScheme,
		UpStatScheme:     UpStatScheme,
		Yields:           Yields,
		ProductionTime:   ProductionTime,
	}
	toRet.AddProperty("building", ID)
	return &toRet
}

func (b *Building) GetID() string {
	return b.ID
}

func (b *Building) GetBasePrice() generaldomain.IBProducts {
	return b.BasePrice
}

func (b *Building) GetUpValScheme() generaldomain.IValScheme {
	return b.UpValScheme
}

func (b *Building) GetUpStatScheme() generaldomain.IValScheme {
	return b.UpStatScheme
}

func (b *Building) GetBaseProducts() generaldomain.IBProducts {
	return b.Yields
}

func (b *Building) GetSpaceEffect() gamemapdomain.ISpaceEffect {
	return b.SpaceEffect
}

func (b *Building) SetSpaceEffect(eff gamemapdomain.ISpaceEffect) bool {
	b.SpaceEffect = eff
	return true
}

func (b *Building) GetProductionIntervalNano() int64 {
	return b.ProductionTime
}

/*
func (b Building) GetVicinity(direction int) domain.IBuilding {
	if dir > 3 || dir < 0 {
		return nil
	}

	return b.vicinity[dir]
}

func (b Building) SetVicinity(b domain.IBuilding, dir int) bool {
	if dir > 3 || dir < 0 {
		return false
	 }
	b.vicinity[dir] = b
	return true
}
*/
