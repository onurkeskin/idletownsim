package buildingview

import (
	effectview "github.com/app/game/presentationlayer/effectview/domain"
	resourceview "github.com/app/game/presentationlayer/resourceview/domain"
	"gopkg.in/mgo.v2/bson"
)

type BuildingsView []BuildingView

type BuildingView struct {
	ID                  bson.ObjectId              `json:"id,omitempty" bson:"_id,omitempty"`
	BuildingName        string                     `bson:"buildingname" json:"buildingname"`
	GlobalIdentifier    string                     `bson:"globalidentifier" json:"globalidentifier"`
	BuildingDescription string                     `bson:"buildingdescription" json:"buildingdescription"`
	ExpectedResource    resourceview.IResourceView `bson:"-" json:"expectedresource"`
	BuildingEffect      effectview.IEffectView     `bson:"-" json:"buildingeffect"`
	//UpValScheme         generalview.ValSchemeView  `json:"upgradepricescheme" bson:"-"`
	//UpStatScheme        generalview.ValSchemeView  `json:"upgraderesultscheme" bson:"-"`
}

type BuildingInstancesView []BuildingInstanceView

type BuildingInstanceView struct {
	ID                  bson.ObjectId              `json:"id,omitempty" bson:"_id,omitempty"`
	BuildingName        string                     `bson:"buildingname" json:"buildingname"`
	GlobalIdentifier    string                     `bson:"globalidentifier" json:"globalidentifier"`
	BuildingDescription string                     `bson:"buildingdescription" json:"buildingdescription"`
	Level               string                     `bson:"-" json:"level"`
	ExpectedResource    resourceview.IResourceView `bson:"-" json:"expectedresource"`
	BuildingEffect      effectview.IEffectView     `bson:"-" json:"buildingeffect"`
	AppliedEffects      []effectview.IEffectView   `bson:"-" json:"effectson"`
}
