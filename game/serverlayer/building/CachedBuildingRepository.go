package building

import (
	"errors"
	"fmt"
	buildingdbdomain "github.com/app/game/serverlayer/building/domain"
	"github.com/app/helpers/version"
	"github.com/app/server/domain"
	"github.com/patrickmn/go-cache"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func NewBuildingRepositoryFactory() buildingdbdomain.IBuildingRepositoryFactory {
	return &BuildingRepositoryFactory{}
}

type BuildingRepositoryFactory struct {
	instancedCache *cache.Cache
}

func (factory *BuildingRepositoryFactory) New(db domain.IDatabase) buildingdbdomain.IBuildingRepository {
	cacheByBuildingID := factory.instancedCache
	if cacheByBuildingID == nil {
		cacheByBuildingID = cache.New(5*time.Minute, 30*time.Second)
		factory.instancedCache = cacheByBuildingID
	}

	return &CachedBuildingRepository{
		dbrepo:            BuildingRepository{db},
		cacheByBuildingID: cacheByBuildingID,
	}
}

type CachedBuildingRepository struct {
	dbrepo            BuildingRepository
	cacheByBuildingID *cache.Cache
}

// Createbuilding Insert new building document into the database
func (repo *CachedBuildingRepository) CreateBuilding(_building buildingdbdomain.IBuildingDB) error {
	return repo.dbrepo.CreateBuilding(_building)
}

// GetUsers Get list of users
func (repo *CachedBuildingRepository) GetBuildings() buildingdbdomain.IBuildingsDB {
	_fromdb := repo.dbrepo.GetBuildings()
	fromdb := _fromdb.(*BuildingsDB)
	for _, v := range *fromdb {
		repo.cacheByBuildingID.Set(v.ID.Hex(), &v, cache.DefaultExpiration)
	}

	return fromdb
}

// GetUser Get user specified by the id
func (repo *CachedBuildingRepository) GetBuildingById(id string) (buildingdbdomain.IBuildingDB, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	cached, ok := repo.cacheByBuildingID.Get(id)
	if ok {
		return cached.(buildingdbdomain.IBuildingDB), nil
	}
	_b, err := repo.dbrepo.GetBuildingById(id)
	if err != nil {
		return _b, err
	}
	b := _b.(*BuildingDB)
	repo.cacheByBuildingID.Set(b.ID.Hex(), b, cache.DefaultExpiration)
	return b, err
}

func (repo *CachedBuildingRepository) GetBuildingsByBuildingId(bid string) (buildingdbdomain.IBuildingsDB, error) {
	bs, err := repo.dbrepo.GetBuildingsByBuildingId(bid)
	if err != nil {
		return bs, err
	}
	for _, v := range *bs.(*BuildingsDB) {
		repo.cacheByBuildingID.Set(v.ID.Hex(), v, cache.DefaultExpiration)
	}
	return bs, err
}

func (repo *CachedBuildingRepository) GetBuildingByBuildingIdVersion(bid string, ver version.Version) (buildingdbdomain.IBuildingDB, error) {
	if !bson.IsObjectIdHex(bid) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", bid))
	}
	_b, err := repo.dbrepo.GetBuildingByBuildingIdVersion(bid, ver)
	if err != nil {
		return _b, err
	}

	b := _b.(*BuildingDB)
	repo.cacheByBuildingID.Set(b.ID.Hex(), b, cache.DefaultExpiration)
	return b, err
}
