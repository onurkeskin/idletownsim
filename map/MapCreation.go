package maps

import (
	"errors"
	mapdomain "github.com/app/map/domain"
	"github.com/app/map/mapmodels"
	"github.com/onurk/externalmath"
	"math"
	"time"
)

const RadiusOfEarthKm float64 = 6371
const radiusOfCircleKm float64 = 1

func CreateMapForPosition(mapProv mapdomain.IMapProvider, repo mapdomain.IMapRepository, pos mapmodels.LatLng) (mapdomain.IMap, error) {
	_, err := GetClosestMapByPosition(repo, pos)
	//fmt.Println(found)
	if err != nil {
		retrieved, err := mapProv.GetMapByLatLang(pos)
		newcreatedmap := retrieved.(*Map)
		if err != nil {
			return nil, err
		}

		newcreatedmap.MapVersion.Ver = "1.0.0"
		newcreatedmap.MapVersion.VerDateAfter = time.Now()
		newcreatedmap.CreatedDate = time.Now()
		err = repo.CreateMap(newcreatedmap)
		if err != nil {
			return newcreatedmap, err
		}
		return newcreatedmap, nil
	}

	return nil, errors.New("A close map is already there")
}

func GetClosestMapByPosition(repo mapdomain.IMapRepository, pos mapmodels.LatLng) (mapdomain.IMap, error) {
	possibilities := getBorderingBoxFromMap(repo, pos)

	found := false
	var BestDist float64 = math.MaxFloat64
	var curBest mapdomain.IMap
	for _, v := range possibilities {
		cLat := v.MapFundamentalCoordinates.Lat
		cLng := v.MapFundamentalCoordinates.Lng
		//dist := math.Acos(math.Sin(cLat))*math.Sin(convertDegreesToRad(cLng)) + math.Cos(cLat)*math.Cos(convertDegreesToRad(cLat)*math.Cos(convertDegreesToRad(cLng)))
		dist := getDistanceFromLatLonInKm(cLat, cLng, pos.Lat, pos.Lng)
		if dist < BestDist && dist < radiusOfCircleKm {
			curBest = &v
			BestDist = dist
			found = true
		}
	}
	if found {
		return curBest, nil
	}

	return nil, errors.New("No Close Maps Found")
}

func GetClosestMapsByPosition(repo mapdomain.IMapRepository, pos mapmodels.LatLng, distanceKM float64) (Maps, error) {
	maxLat := pos.Lat + convertRadToDegrees(distanceKM/RadiusOfEarthKm)
	minLat := pos.Lat - convertRadToDegrees(distanceKM/RadiusOfEarthKm)
	maxLng := pos.Lng + convertRadToDegrees(math.Asin(distanceKM/RadiusOfEarthKm)/math.Cos(convertDegreesToRad(pos.Lat)))
	minLng := pos.Lng - convertRadToDegrees(math.Asin(distanceKM/RadiusOfEarthKm)/math.Cos(convertDegreesToRad(pos.Lat)))
	toRet := repo.FilterMaps(mapmodels.LatLng{Lat: minLat, Lng: minLng}, mapmodels.LatLng{Lat: maxLat, Lng: maxLng}, "", math.MaxInt32, "")
	//fmt.Println(toRet)
	return *toRet.(*Maps), nil
}

func getBorderingBoxFromMap(repo mapdomain.IMapRepository, pos mapmodels.LatLng) Maps {
	maxLat := pos.Lat + convertRadToDegrees(radiusOfCircleKm/RadiusOfEarthKm)
	minLat := pos.Lat - convertRadToDegrees(radiusOfCircleKm/RadiusOfEarthKm)
	maxLng := pos.Lng + convertRadToDegrees(math.Asin(radiusOfCircleKm/RadiusOfEarthKm)/math.Cos(convertDegreesToRad(pos.Lat)))
	minLng := pos.Lng - convertRadToDegrees(math.Asin(radiusOfCircleKm/RadiusOfEarthKm)/math.Cos(convertDegreesToRad(pos.Lat)))
	toRet := repo.FilterMaps(mapmodels.LatLng{Lat: minLat, Lng: minLng}, mapmodels.LatLng{Lat: maxLat, Lng: maxLng}, "", math.MaxInt32, "")
	//fmt.Println(toRet)
	return *toRet.(*Maps)
}

func getDistanceFromLatLonInKm(lat1, lon1, lat2, lon2 float64) float64 {
	var dLat = convertDegreesToRad(lat2 - lat1) // convertDegreesToRad below
	var dLon = convertDegreesToRad(lon2 - lon1)
	var a = math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(convertDegreesToRad(lat1))*math.Cos(convertDegreesToRad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	var d = RadiusOfEarthKm * c // Distance in km
	return d
}

func LatLonNormalizer(latln mapmodels.LatLng) mapmodels.LatLng {
	externalmath.Round(latln.Lat, 0.00005)

	toReturn := mapmodels.LatLng{}

	return toReturn
}

func convertRadToDegrees(rad float64) float64 {
	return rad * 180 / math.Pi
}

func convertDegreesToRad(degree float64) float64 {
	return degree * math.Pi / 180
}
