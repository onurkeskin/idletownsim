package game

import (
	gamedomain "github.com/app/game/applayer/game/domain"
)

type IGameManager interface {
	AddGame(game gamedomain.IGameEnvironment)
	GetGameByGameID(id string) (gamedomain.IGameEnvironment, bool)
	GetGamesByUserID(id string) (gamedomain.IGameEnvironments, bool)
	RemoveGameByGameID(id string)
	RemoveGamesByUserID(id string)
}
