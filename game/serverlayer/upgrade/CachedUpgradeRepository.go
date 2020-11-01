package upgrade

import (
	"errors"
	"fmt"

	dbmodels "github.com/app/game/serverlayer/dbmodels"
	upgradedbdomain "github.com/app/game/serverlayer/upgrade/domain"

	"github.com/app/server/domain"
	"gopkg.in/mgo.v2/bson"
	"time"

	"github.com/patrickmn/go-cache"
)

func NewUpgradeRepositoryFactory() upgradedbdomain.IUpgradeRepositoryFactory {
	return &UpgradeRepositoryFactory{}
}

type UpgradeRepositoryFactory struct {
	instancedCache *cache.Cache
}

func (factory *UpgradeRepositoryFactory) New(db domain.IDatabase) upgradedbdomain.IUpgradeRepository {
	cacheByUpgradeID := factory.instancedCache
	if cacheByUpgradeID == nil {
		cacheByUpgradeID = cache.New(5*time.Minute, 30*time.Second)
		factory.instancedCache = cacheByUpgradeID
	}

	return &CachedUpgradeRepository{
		dbrepo:           UpgradeRepository{db},
		cacheByUpgradeID: cacheByUpgradeID,
	}
}

type CachedUpgradeRepository struct {
	dbrepo           UpgradeRepository
	cacheByUpgradeID *cache.Cache
}

// CreateUpgrade Insert new upgrade document into the database
func (repo *CachedUpgradeRepository) CreateUpgrade(_upgrade *dbmodels.UpgradeDB) error {
	return repo.dbrepo.CreateUpgrade(_upgrade)
}

// GetUsers Get list of users
func (repo *CachedUpgradeRepository) GetUpgrades() dbmodels.UpgradesDB {
	fromdb := repo.dbrepo.GetUpgrades()
	for _, v := range fromdb {
		repo.cacheByUpgradeID.Set(v.ID.Hex(), &v, cache.DefaultExpiration)
	}
	return fromdb
}

// GetUser Get user specified by the id
func (repo *CachedUpgradeRepository) GetUpgradeById(id string) (*dbmodels.UpgradeDB, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	cached, ok := repo.cacheByUpgradeID.Get(id)
	if ok {
		return cached.(*dbmodels.UpgradeDB), nil
	}
	b, err := repo.dbrepo.GetUpgradeById(id)
	if err != nil {
		return b, err
	}
	repo.cacheByUpgradeID.Set(b.ID.Hex(), b, cache.DefaultExpiration)
	return b, err
}
