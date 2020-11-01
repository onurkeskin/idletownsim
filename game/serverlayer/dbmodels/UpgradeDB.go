package dbmodels

import (
	//"github.com/app/game/applayer/game"
	upgrade "github.com/app/game/applayer/upgrade"
	"github.com/app/helpers/version"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type UpgradesDB []UpgradeDB
type UpgradeDB struct {
	ID                 bson.ObjectId           `json:"id,omitempty" bson:"_id,omitempty"`
	UpgradeID          bson.ObjectId           `json:"uid,omitempty" bson:"uid,omitempty"`
	UpgradeRequirement *upgrade.Requirement    `json:"req" bson:"req"`
	Value              ValuesDB                `json:"Value" bson:"Value"`
	UpTargets          []upgrade.UpgradeTarget `json:"targets" bson:"targets"`
	Effects            EffectsDB               `json:"effects" bson:"effects"`

	Ver          version.Version `json:"version" bson:"version"`
	CreationDate time.Time       `json:"creationdate" bson:"creationdate"`
}

func (v UpgradeDB) FormUpgrade(rs ResourcesDB) (*upgrade.Upgrade, error) {
	effs, err := v.Effects.FormEffects()
	if err != nil {
		return nil, err
	}
	toRet := upgrade.Upgrade{
		ID:                 v.ID.Hex(),
		UpgradeID:          v.UpgradeID.Hex(),
		UpgradeRequirement: v.UpgradeRequirement,
		Price:              v.Value.FormIBProducts(rs),
		UpTargets:          v.UpTargets,
		Effects:            effs,
	}
	return &toRet, nil
}

/*
func FormFromUpgrade(v upgrade.Upgrade) (*UpgradeDB, error) {

	switch v := v.(type) {
	case *buildinginstance.BuildingInstance:

		toRet := UpgradeDB{
			BuildableType: "buildinginstance",
			Buildable:     v,
		}
		return &toRet, nil
	default:
		return nil, errors.New("Cant understand the scheme type")
	}
	return nil, errors.New("Cant understand the scheme type")
}
*/
