package maps

import (
	"github.com/app/map/mapmodels"
)

type IMapProvider interface {
	GetMapByLatLang(latlng mapmodels.LatLng) (IMap, error)
}
