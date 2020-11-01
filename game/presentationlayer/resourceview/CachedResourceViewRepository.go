package resourceview

import (
	"errors"
	"fmt"

	ResourceViewdomain "github.com/app/game/presentationlayer/resourceview/domain"

	"github.com/app/server/domain"
	"gopkg.in/mgo.v2/bson"
	"time"

	"github.com/patrickmn/go-cache"
)

func NewResourceViewRepositoryFactory() ResourceViewdomain.IResourceViewRepositoryFactory {
	return &ResourceViewRepositoryFactory{}
}

type ResourceViewRepositoryFactory struct {
	instancedCache *cache.Cache
}

func (factory *ResourceViewRepositoryFactory) New(db domain.IDatabase) ResourceViewdomain.IResourceViewRepository {
	cacheByResourceViewID := factory.instancedCache
	if cacheByResourceViewID == nil {
		cacheByResourceViewID = cache.New(5*time.Minute, 30*time.Second)
		factory.instancedCache = cacheByResourceViewID
	}

	return &CachedResourceViewRepository{
		dbrepo:                ResourceViewRepository{db},
		cacheByResourceViewID: cacheByResourceViewID,
	}
}

type CachedResourceViewRepository struct {
	dbrepo                ResourceViewRepository
	cacheByResourceViewID *cache.Cache
}

// CreateResourceView Insert new ResourceView document into the database
func (repo *CachedResourceViewRepository) CreateResourceView(_ResourceView ResourceViewdomain.IResourceView) error {
	return repo.dbrepo.CreateResourceView(_ResourceView)
}

// GetUsers Get list of users
func (repo *CachedResourceViewRepository) GetResourceViews() ResourceViewdomain.IResourcesView {
	_fromdb := repo.dbrepo.GetResourceViews()
	fromdb := _fromdb.(ResourcesView)
	for _, v := range fromdb {
		repo.cacheByResourceViewID.Set(v.ID.Hex(), &v, cache.DefaultExpiration)
	}
	return fromdb
}

// GetUser Get user specified by the id
func (repo *CachedResourceViewRepository) GetResourceViewById(id string) (ResourceViewdomain.IResourceView, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	cached, ok := repo.cacheByResourceViewID.Get(id)
	if ok {
		return cached.(*ResourceView), nil
	}
	_b, err := repo.dbrepo.GetResourceViewById(id)
	if err != nil {
		return _b, err
	}
	b := _b.(*ResourceView)
	repo.cacheByResourceViewID.Set(b.ID.Hex(), b, cache.DefaultExpiration)
	return b, err
}
