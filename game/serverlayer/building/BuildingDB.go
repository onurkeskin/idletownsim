package building

import (
	building "github.com/app/game/applayer/building"
	buildingdomain "github.com/app/game/applayer/building/domain"
	gamemapdomain "github.com/app/game/applayer/gamemap/domain"
	general "github.com/app/game/applayer/general"
	generaldomain "github.com/app/game/applayer/general/domain"
	"github.com/app/helpers/version"
	"gopkg.in/mgo.v2/bson"
	"time"
	//"github.com/app/game/serverlayer/building/domain"
	"github.com/app/game/serverlayer/dbmodels"
	resourcedomain "github.com/app/game/serverlayer/gameresources/domain"
)

type BuildingsDB []BuildingDB

type BuildingDB struct {
	*general.ObjectProperties
	ID  bson.ObjectId `json:"-" bson:"_id"`
	BID bson.ObjectId `json:"-" bson:"bid"`
	//UpgradeID bson.ObjectId `json:"uid,omitempty" bson:"uid,omitempty"`
	BaseValue    dbmodels.ValuesDB    `json:"value" bson:"value"`
	UpValScheme  dbmodels.ValSchemeDB `json:"upgradevaluescheme" bson:"upgradevaluescheme"`
	UpStatScheme dbmodels.ValSchemeDB `json:"upgraderesultscheme" bson:"upgraderesultscheme"`

	Yields         dbmodels.ValuesDB `json:"yields" bson:"yields"`
	ProductionTime int64             `json:"time" bson:"time"`

	SpaceEffect dbmodels.EffectDB `json:"buildingeffect" bson:"buildingeffect"`

	CreationDate time.Time       `json:"creationdate" bson:"creationdate"`
	Ver          version.Version `json:"version" bson:"version"`
}

func (builDB *BuildingDB) FormIBuilding(provider resourcedomain.IResourceRepository) buildingdomain.IBuilding {
	UpValScheme := builDB.UpValScheme.FormIValScheme()
	UpStatScheme := builDB.UpStatScheme.FormIValScheme()

	SpaceEffect, err := builDB.SpaceEffect.FormEffect()
	if err != nil {
		panic(err)
	}
	BaseValue := generaldomain.IBProducts{}
	for _, v := range builDB.BaseValue {
		res, err := provider.GetResourceById(v.ID.Hex())
		if err != nil {
			panic("resource not found")
		}
		product := v.FormIBProduct(*res)

		BaseValue = append(BaseValue, product)
	}
	Yields := generaldomain.IBProducts{}
	for _, v := range builDB.Yields {
		res, err := provider.GetResourceById(v.ID.Hex())
		if err != nil {
			panic("resource not found")
		}
		product := v.FormIBProduct(*res)

		Yields = append(Yields, product)
	}

	buil := building.NewBuilding(
		builDB.ID.Hex(),
		BaseValue,
		UpValScheme,
		UpStatScheme,
		Yields,
		builDB.ProductionTime,
	)
	if SpaceEffect != nil {
		castSE := SpaceEffect.(gamemapdomain.ISpaceEffect)
		buil.SetSpaceEffect(castSE)
	}
	return buil
}
