package maps

import (
	wrappers "github.com/app/goWrappers/openCvExportWrappers"
	"github.com/app/map/mapmodels"
)

type IMaps interface{}

type IMap interface {
	GetID() string

	IsValid() bool

	GetMapFundementalCoordinates() mapmodels.LatLng

	GetMapBytes() []byte
	GetMapRelations() wrappers.AllRelations
}
