package maps

import (
	wrappers "github.com/app/goWrappers/openCvExportWrappers"
	"github.com/app/helpers/version"
	mapmodels "github.com/app/map/mapmodels"
	//"github.com/app/server/domain"
	"gopkg.in/mgo.v2/bson"
	//"strings"
	"time"
)

type Maps []Map

type Map struct {
	ID                        bson.ObjectId         `json:"id,omitempty" bson:"_id,omitempty"`
	MapIdentifierAdress       string                `json:"mapidentifieradress" bson:"mapidentifieradress"`
	MapCompleteAddress        string                `json:"mapadress" bson:"mapaddress"`
	MapRaw                    []byte                `json:"rawmap" bson:"rawmap"`
	MapFundamentalCoordinates mapmodels.LatLng      `json:"fundamentalcoordinates" bson:"fundamentalcoordinates"` // IN DEGREES
	ParsedRelations           wrappers.AllRelations `json:"parsedrelations" bson:"parsedrelations"`

	InstancesCreated int64           `json:"instancescreated" bson:"instancescreated"`
	MapVersion       version.Version `json:"version" bson:"version"`
	Status           string          `json:"status,omitempty" bson:"status"`
	LastModifiedDate time.Time       `json:"lastModifiedDate" bson:"lastModifiedDate"`
	CreatedDate      time.Time       `json:"createdDate,omitempty" bson:"createdDate"`
}

func (m *Map) GetID() string {
	return m.ID.Hex()
}

func (m *Map) IsValid() bool {
	//needs change
	return true
}

func (m *Map) GetMapFundementalCoordinates() mapmodels.LatLng {
	return m.MapFundamentalCoordinates
}

func (m *Map) GetMapBytes() []byte {
	return m.MapRaw
}

func (m *Map) GetMapRelations() wrappers.AllRelations {
	return m.ParsedRelations
}
