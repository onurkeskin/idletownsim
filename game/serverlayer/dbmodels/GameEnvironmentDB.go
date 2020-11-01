package dbmodels

import (
	gamedomain "github.com/app/game/applayer/game/domain"
	//gamemapdomain "github.com/app/game/applayer/gamemap/domain"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type GameEnvironmentsDB []GameEnvironmentDB

type GameEnvironmentDB struct {
	GameID          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	BelongingUserID bson.ObjectId `json:"uid,omitempty" bson:"_uid,omitempty"`
	BuiltOnMapID    bson.ObjectId `json:"mid,omitempty" bson:"_mid,omitempty"`

	GameCreationTime  time.Time `json:"gamect" bson:"gamect"`
	LastIterationTime time.Time `json:"gamelitime" bson:"gamelitime"`

	Man GameUpgradeManagerDB `json:"upgrades" bson:upgrademanager`

	Resources ValuesDB `json:"resources" bson:"resources"`
	Gm        GameDB   `json:"game" bson:game`
}

func FormFromIGameEnvironment(gameEnv gamedomain.IGameEnvironment) (*GameEnvironmentDB, error) {
	products, err := FormFromIBProducts(gameEnv.GetResources())
	if err != nil {
		return nil, err
	}
	toRet := GameEnvironmentDB{
		GameID:            bson.ObjectIdHex(gameEnv.GetGameID()),
		BelongingUserID:   bson.ObjectIdHex(gameEnv.GetBelongingUserID()),
		BuiltOnMapID:      bson.ObjectIdHex(gameEnv.GetMapID()),
		GameCreationTime:  gameEnv.GetCreationTime(),
		LastIterationTime: gameEnv.GetLastIterationTime(),
		Man:               GameUpgradeManagerDB{BoughtGameUpgradeIDs: gameEnv.GetBoughtUpgrades()},
		Resources:         products,
		/*GameDB{
			GMap: GameMapDB{
				Version: gameEnv.GetGame().GetGameMap().GetMapVersion(),
				MapDate: gameEnv.GetGame().GetGameMap().GetMapDate(),
			},
		},*/
	}
	return &toRet, nil
}

/*
func FormFromIGame(game gamedomain.IGame) (*GameDB, error) {
	toRet := GameDB{
		GMap:FormFromIGameMap()
	}
}

func FormFromIGameMap(gamemap gamemapdomain.IGameMap) (*GameMapDB, error) {
	toRet := GameMapDB{
		ID:gamemap.
	}
}
*/
