package game

import (
	dbmodels "github.com/app/game/serverlayer/dbmodels"
	"github.com/app/server/domain"
)

type IGameRepositoryFactory interface {
	New(db domain.IDatabase) IGameRepository
}

type IGameRepository interface {
	CreateGame(_game *dbmodels.GameEnvironmentDB) error
	GetGamesByUserId(id string) (dbmodels.GameEnvironmentsDB, error)
	GetGamesByMapId(id string) (dbmodels.GameEnvironmentsDB, error)
	GetGameByGameId(id string) (*dbmodels.GameEnvironmentDB, error)
	GetGames() dbmodels.GameEnvironmentsDB
	UpdateGame(id string, gUpdate *dbmodels.GameEnvironmentDB) (*dbmodels.GameEnvironmentDB, error)
	GameExistsByUserVMapId(userid string, mapid string) (bool, error)
}
