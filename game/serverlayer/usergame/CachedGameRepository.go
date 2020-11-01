package game

import (
	"errors"
	"fmt"

	dbmodels "github.com/app/game/serverlayer/dbmodels"
	gamedbdomain "github.com/app/game/serverlayer/usergame/domain"
	"github.com/patrickmn/go-cache"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"log"
	"time"

	"github.com/app/server/domain"
)

func NewGameRepositoryFactory() gamedbdomain.IGameRepositoryFactory {
	return &GameRepositoryFactory{}
}

type GameRepositoryFactory struct {
	instancedCacheByUser *cache.Cache
	instancedCacheByGame *cache.Cache
}

type CachedGameRepository struct {
	dbrepo        GameRepository
	cacheByUserID *cache.Cache
	cacheByGameID *cache.Cache
}

func (factory *GameRepositoryFactory) New(db domain.IDatabase) gamedbdomain.IGameRepository {

	cacheByUserID := factory.instancedCacheByUser
	if cacheByUserID == nil {
		cacheByUserID = cache.New(5*time.Minute, 30*time.Second)
		factory.instancedCacheByUser = cacheByUserID
	}
	cacheByGameID := factory.instancedCacheByGame
	if cacheByGameID == nil {
		cacheByGameID = cache.New(5*time.Minute, 30*time.Second)
		factory.instancedCacheByUser = cacheByGameID
	}

	return &CachedGameRepository{
		dbrepo:        GameRepository{db},
		cacheByUserID: cacheByUserID,
		cacheByGameID: cacheByGameID,
	}
}

func (repo *CachedGameRepository) CreateGame(_game *dbmodels.GameEnvironmentDB) error {
	err := repo.dbrepo.CreateGame(_game)
	if err != nil {
		return err
	}
	repo.cacheByUserID.Set(_game.GameID.Hex(), &_game, cache.DefaultExpiration)
	return err
}

func (repo *CachedGameRepository) GetGamesByUserId(id string) (dbmodels.GameEnvironmentsDB, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	val, ok := repo.cacheByUserID.Get(id)
	if ok {
		repo.cacheByUserID.Set(id, val, cache.DefaultExpiration)
		return val.(dbmodels.GameEnvironmentsDB), nil
	} else {
		games, err := repo.dbrepo.GetGamesByUserId(id)
		if err != nil {
			return games, err
		}
		repo.cacheByUserID.Set(id, games, cache.DefaultExpiration)
		return games, err
	}
}

func (repo *CachedGameRepository) GetGamesByMapId(id string) (dbmodels.GameEnvironmentsDB, error) {

	return repo.dbrepo.GetGamesByMapId(id)
}

func (repo *CachedGameRepository) GetGameByGameId(id string) (*dbmodels.GameEnvironmentDB, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	val, ok := repo.cacheByGameID.Get(id)
	if ok {
		repo.cacheByGameID.Set(id, val, cache.DefaultExpiration)
		return val.(*dbmodels.GameEnvironmentDB), nil
	} else {
		games, err := repo.dbrepo.GetGameByGameId(id)
		if err != nil {
			return games, err
		}
		repo.cacheByGameID.Set(id, games, cache.DefaultExpiration)
		return games, err
	}
}

func (repo *CachedGameRepository) GetGames() dbmodels.GameEnvironmentsDB {
	return repo.dbrepo.GetGames()
}

func (repo *CachedGameRepository) UpdateGame(id string, gUpdate *dbmodels.GameEnvironmentDB) (*dbmodels.GameEnvironmentDB, error) {
	return repo.dbrepo.UpdateGame(id, gUpdate)
}

func (repo *CachedGameRepository) GameExistsByUserVMapId(userid string, mapid string) (bool, error) {
	return repo.dbrepo.GameExistsByUserVMapId(userid, mapid)
}
