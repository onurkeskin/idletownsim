package resourceview

import (
	"gopkg.in/mgo.v2/bson"
)

type ResourcesView []ResourceView

type ResourceView struct {
	ID               bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	GlobalIdentifier string        `json:"globalidentifier" bson:"globalidentifier"`
	ResourceName     string        `json:"resourcename" bson:"resourcename"`
	ResourceCount    float64       `json:"resourcecount" bson:"-"`
}
