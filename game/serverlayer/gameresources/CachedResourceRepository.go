package resource

import (
	"errors"
	"fmt"
	dbmodels "github.com/app/game/serverlayer/dbmodels"
	resourcedbdomain "github.com/app/game/serverlayer/gameresources/domain"
	"github.com/app/server/domain"
	"github.com/patrickmn/go-cache"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func NewResourceRepositoryFactory() resourcedbdomain.IResourceRepositoryFactory {
	return &ResourceRepositoryFactory{}
}

type ResourceRepositoryFactory struct {
	instancedCache *cache.Cache
}

func (factory *ResourceRepositoryFactory) New(db domain.IDatabase) resourcedbdomain.IResourceRepository {
	cacheByResourceID := factory.instancedCache
	if cacheByResourceID == nil {
		cacheByResourceID = cache.New(5*time.Minute, 30*time.Second)
		factory.instancedCache = cacheByResourceID
	}

	return &CachedResourceRepository{
		dbrepo:            ResourceRepository{db},
		cacheByResourceID: cacheByResourceID,
	}
}

type CachedResourceRepository struct {
	dbrepo            ResourceRepository
	cacheByResourceID *cache.Cache
}

// Createbuilding Insert new building document into the database
func (repo *CachedResourceRepository) CreateResource(res *dbmodels.ResourceDB) error {
	return repo.dbrepo.CreateResource(res)
}

// GetUsers Get list of users
func (repo *CachedResourceRepository) GetResources() dbmodels.ResourcesDB {
	fromdb := repo.dbrepo.GetResources()
	for _, v := range fromdb {
		repo.cacheByResourceID.Set(v.ID.Hex(), &v, cache.DefaultExpiration)
	}
	return fromdb
}

// GetUser Get user specified by the id
func (repo *CachedResourceRepository) GetResourceById(id string) (*dbmodels.ResourceDB, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	cached, ok := repo.cacheByResourceID.Get(id)
	if ok {
		return cached.(*dbmodels.ResourceDB), nil
	}
	b, err := repo.dbrepo.GetResourceById(id)
	if err != nil {
		return b, err
	}
	repo.cacheByResourceID.Set(b.ID.Hex(), b, cache.DefaultExpiration)
	return b, err
}
