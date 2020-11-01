package dbmodels

import (
	"github.com/app/helpers/version"
	//"github.com/app/game/applayer/general"
	//"github.com/app/game/applayer/upgrade"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type GameDB struct {
	GMap GameMapDB `json:"gmap" bson:"gmap"`
}

type GameMapDB struct {
	ID bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	//Spaces  SpacesDB      `json:"spaces" bson:"spaces"`
	MapDate time.Time       `json:"mapdate" bson:"mapdate"`
	Version version.Version `json:"version" bson:"version"`
}

type SpacesDB []SpaceDB

type SpaceDB struct {
	ID             bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	BelongingMapID bson.ObjectId `json:"mapid,omitempty" bson:"mapid,omitempty"`

	InMapID string                  `json:"inmapid" bson:"inmapid"`
	Element *TileBuildableElementDB `json:"occupier" bson:"occupier"`
	Around  [8][]string             `json:"around" bson:"around"`
}

/*
type TileBuildableElementDB struct {
	BuildableType    string          `json:"type,omitempty" bson:"type,omitempty"`
	BuildableVersion version.Version `json:"version,omitempty" bson:"version,omitempty"`
	BuildableJSON    []byte          `json:"buildable,omitempty" bson:"buildable,omitempty"`
}
*/

type GameUpgradeManagerDB struct {
	BoughtGameUpgradeIDs []string `json:"boughtupgrades" bson:boughtupgrades`
}

type ValSchemeDB struct {
	SchemeType string `json:"type,omitempty" bson:"type,omitempty"`
	SchemeJson []byte `json:"scheme,omitempty" bson:"scheme,omitempty"`
}

type ValuesDB []ValueDB

type ValueDB struct {
	ID            bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	ResourceCount float64       `json:"count,omitempty" bson:"count,omitempty"`
}

type ResourcesDB []ResourceDB
type ResourceDB struct {
	ID           bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	ResType      string        `json:"type,omitempty" bson:"type,omitempty"`
	CreationTime time.Time     `json:"-" bson:"createdonon"`
}

type EffectsDB []EffectDB
type EffectDB struct {
	EffectType string `json:"type,omitempty" bson:"type,omitempty"`
	EffectJson []byte `json:"effect,omitempty" bson:"effect,omitempty"`
}
