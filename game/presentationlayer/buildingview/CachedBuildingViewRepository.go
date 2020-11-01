package buildingview

import (
	"errors"
	"fmt"
	buildingviewdomain "github.com/app/game/presentationlayer/buildingview/domain"
	"github.com/app/server/domain"
	"github.com/patrickmn/go-cache"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func NewBuildingViewRepositoryFactory() buildingviewdomain.IBuildingViewRepositoryFactory {
	return &BuildingViewRepositoryFactory{}
}

type BuildingViewRepositoryFactory struct {
	instancedCache *cache.Cache
}

func (factory *BuildingViewRepositoryFactory) New(db domain.IDatabase) buildingviewdomain.IBuildingViewRepository {
	cacheByBuildingViewID := factory.instancedCache
	if cacheByBuildingViewID == nil {
		cacheByBuildingViewID = cache.New(5*time.Minute, 30*time.Second)
		factory.instancedCache = cacheByBuildingViewID
	}

	return &CachedBuildingViewRepository{
		dbrepo:                BuildingViewRepository{db},
		cacheByBuildingViewID: cacheByBuildingViewID,
	}
}

type CachedBuildingViewRepository struct {
	dbrepo                BuildingViewRepository
	cacheByBuildingViewID *cache.Cache
}

func (repo *CachedBuildingViewRepository) CreateBuildingView(_b buildingviewdomain.IBuildingView) error {
	return repo.dbrepo.CreateBuildingView(_b)
}

// GetUsers Get list of users
func (repo *CachedBuildingViewRepository) GetBuildingsView() buildingviewdomain.IBuildingsView {
	_fromdb := repo.dbrepo.GetBuildingsView()
	fromdb := _fromdb.(BuildingsView)
	for _, v := range fromdb {
		repo.cacheByBuildingViewID.Set(v.ID.Hex(), &v, cache.DefaultExpiration)
	}
	return fromdb
}

// GetUser Get user specified by the id
func (repo *CachedBuildingViewRepository) GetBuildingViewById(id string) (buildingviewdomain.IBuildingView, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	cached, ok := repo.cacheByBuildingViewID.Get(id)
	if ok {
		return cached.(*BuildingView), nil
	}
	_b, err := repo.dbrepo.GetBuildingViewById(id)
	if err != nil {
		return _b, err
	}
	b := _b.(*BuildingView)
	repo.cacheByBuildingViewID.Set(b.ID.Hex(), b, cache.DefaultExpiration)
	return b, err
}
