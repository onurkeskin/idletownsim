package maps

import (
	"errors"
	"fmt"
	mapdomain "github.com/app/map/domain"
	"github.com/app/map/mapmodels"
	"github.com/patrickmn/go-cache"
	"gopkg.in/mgo.v2/bson"
	//"log"
	"github.com/app/server/domain"
	"time"
)

func NewMapRepositoryFactory() mapdomain.IMapRepositoryFactory {
	return &MapRepositoryFactory{}
}

type MapRepositoryFactory struct {
	instancedCache *cache.Cache
}

func (factory *MapRepositoryFactory) New(db domain.IDatabase) mapdomain.IMapRepository {
	cacheByMapID := factory.instancedCache
	if cacheByMapID == nil {
		cacheByMapID = cache.New(5*time.Minute, 10*time.Minute)
		factory.instancedCache = cacheByMapID
	}

	return &CachedMapRepository{
		dbrepo:       MapRepository{db},
		cacheByMapID: cacheByMapID}
}

type CachedMapRepository struct {
	dbrepo       MapRepository
	cacheByMapID *cache.Cache
}

func (repo *CachedMapRepository) CreateMap(_map mapdomain.IMap) error {
	err := repo.dbrepo.CreateMap(_map)
	if err != nil {
		return err
	}

	repo.cacheByMapID.Set(_map.GetID(), _map, cache.DefaultExpiration)
	return err
}

func (repo *CachedMapRepository) CountMaps(field string, query string) int {
	return repo.dbrepo.CountMaps(field, query)
}

func (repo *CachedMapRepository) GetMapById(id string) (mapdomain.IMap, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	cachedVer, ok := repo.cacheByMapID.Get(id)
	if ok {
		return cachedVer.(mapdomain.IMap), nil
	}
	toRet, err := repo.dbrepo.GetMapById(id)
	if err != nil {
		return toRet, err
	}
	repo.cacheByMapID.Set(toRet.GetID(), toRet, cache.DefaultExpiration)
	return toRet, err
}

func (repo *CachedMapRepository) UpdateMap(id string, _inMap mapdomain.IMap) (mapdomain.IMap, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}
	updated, err := repo.dbrepo.UpdateMap(id, _inMap)
	if err != nil {
		return updated, err
	}
	repo.cacheByMapID.Replace(updated.GetID(), updated, cache.DefaultExpiration)
	return updated, err
}

func (repo *CachedMapRepository) DeleteMap(id string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}
	err := repo.dbrepo.DeleteMap(id)
	if err != nil {
		return err
	}
	repo.cacheByMapID.Delete(id)
	return err
}

func (repo *CachedMapRepository) DeleteAllMaps() error {
	//TODO
	return repo.DeleteAllMaps()
}

func (repo *CachedMapRepository) FilterMaps(greater mapmodels.LatLng, smaller mapmodels.LatLng, lastID string, limit int, sort string) mapdomain.IMaps {
	return repo.dbrepo.FilterMaps(greater, smaller, lastID, limit, sort)
}
