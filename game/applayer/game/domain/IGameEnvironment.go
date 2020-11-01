package domain

import (
	buildingdomain "github.com/app/game/applayer/building/domain"
	"github.com/app/game/applayer/buildinginstance"
	buildinginstancedomain "github.com/app/game/applayer/buildinginstance/domain"
	generaldomain "github.com/app/game/applayer/general/domain"
	"time"
)

type IGameEnvironment interface {
	GetGameID() string
	GetBelongingUserID() string
	GetMapID() string

	BuildingBuyable(building buildingdomain.IBuilding) bool
	BuyBuilding(building buildingdomain.IBuilding, sid string, options buildinginstance.FormFromBuildingOptions) (buildinginstancedomain.IBuildingInstance, error)

	BuyUpgrade(u IUpgrade) error
	AddUpgrade(u IUpgrade) error
	//IsEligible(u IUpgrade) error
	RemoveUpgrade(u IUpgrade)

	PlayFor(t time.Duration)

	GetGame() IGame
	GetBoughtUpgrades() []string

	GetResources() generaldomain.IBProducts
	GetGameStats()

	GetLastIterationTime() time.Time
	SetLastIterationTime(a time.Time)
	GetCreationTime() time.Time
	SetCreationTime(a time.Time)
}
type IGameEnvironments interface{}
