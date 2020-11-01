package game

import (
	gamemap "github.com/app/game/applayer/gamemap"
	dbmodels "github.com/app/game/serverlayer/dbmodels"
	wrappers "github.com/app/goWrappers/openCvExportWrappers"
	maps "github.com/app/map"
	//mapdomain "github.com/app/map/domain"
	"github.com/app/map/mapmodels"
	//"github.com/onurk/externalmath"
	"gopkg.in/mgo.v2/bson"
	//"math"
	//"time"
)

//TODO CHANGE ALGORITHMS FOR SPACE AND GAME CREATION
func (res *Resource) CreateGameForPosition(pos mapmodels.LatLng) (*dbmodels.GameEnvironmentDB, *dbmodels.SpacesDB, error) {
	maprepo := res.MapRepository(nil)

	var uMap *maps.Map
	m, err := maps.GetClosestMapByPosition(maprepo, pos)
	if err != nil {
		m, err = maps.CreateMapForPosition(res.MapProvider, maprepo, pos)
		if err != nil {
			return nil, nil, err
		}
		uMap = m.(*maps.Map)
	} else {
		uMap = m.(*maps.Map)
	}
	toCreateSpaces := formSpacesFromRelations(m.GetMapRelations())

	toCreate := dbmodels.GameEnvironmentDB{
		BuiltOnMapID: uMap.ID,
		Gm: dbmodels.GameDB{
			dbmodels.GameMapDB{
				MapDate: uMap.LastModifiedDate,
				Version: uMap.MapVersion},
		},
	}

	return &toCreate, &toCreateSpaces, nil
}
func (res *Resource) CreateGameForMapID(mapid string) (*dbmodels.GameEnvironmentDB, *dbmodels.SpacesDB, error) {
	maprepo := res.MapRepository(nil)

	m, err := maprepo.GetMapById(mapid)
	if err != nil {
		return nil, nil, err
	}

	toCreateSpaces := formSpacesFromRelations(m.GetMapRelations())
	uMap := m.(*maps.Map)

	toCreate := dbmodels.GameEnvironmentDB{
		BuiltOnMapID: uMap.ID,
		Gm: dbmodels.GameDB{
			dbmodels.GameMapDB{
				ID:      uMap.ID,
				MapDate: uMap.LastModifiedDate,
				Version: uMap.MapVersion},
		},
	}

	return &toCreate, &toCreateSpaces, nil
}

/*
func (res *Resource) CreateGameForPosition(pos mapmodels.LatLng) (*dbmodels.GameEnvironmentDB, *dbmodels.SpacesDB, error) {
}
*/

func formSpacesFromRelations(rels wrappers.AllRelations) dbmodels.SpacesDB {
	toRet := dbmodels.SpacesDB{}

	for _, v := range rels {
		toAdd := dbmodels.SpaceDB{
			ID:      bson.NewObjectId(),
			InMapID: v.SelfID,
		}

		for _, east := range v.East {
			toAdd.Around[gamemap.East] = append(toAdd.Around[gamemap.East], east.FreePositionID)
		}
		for _, west := range v.West {
			toAdd.Around[gamemap.West] = append(toAdd.Around[gamemap.West], west.FreePositionID)
		}
		for _, north := range v.North {
			toAdd.Around[gamemap.North] = append(toAdd.Around[gamemap.North], north.FreePositionID)
		}
		for _, south := range v.South {
			toAdd.Around[gamemap.South] = append(toAdd.Around[gamemap.South], south.FreePositionID)
		}
		for _, northeast := range v.NorthEast {
			toAdd.Around[gamemap.Northeast] = append(toAdd.Around[gamemap.Northeast], northeast.FreePositionID)
		}
		for _, northwest := range v.NorthWest {
			toAdd.Around[gamemap.Northwest] = append(toAdd.Around[gamemap.Northwest], northwest.FreePositionID)
		}
		for _, southwest := range v.SouthWest {
			toAdd.Around[gamemap.Southwest] = append(toAdd.Around[gamemap.Southwest], southwest.FreePositionID)
		}
		for _, southeast := range v.SouthEast {
			toAdd.Around[gamemap.Southeast] = append(toAdd.Around[gamemap.Southeast], southeast.FreePositionID)
		}

		toRet = append(toRet, toAdd)
	}

	return toRet
}

/*
func findEPosInMap(find wrappers.EPosition, findIn map[wrappers.EPosition]dbmodels.SpaceDB) *dbmodels.SpaceDB {
	for k, v := range findIn {
		if k.P1 == find.P1 && k.P2 == find.P2 {
			return &v
		}
	}
	return nil
}
*/
