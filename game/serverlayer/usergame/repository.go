package game

import (
	"errors"
	"fmt"
	dbmodels "github.com/app/game/serverlayer/dbmodels"
	"time"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"log"
	//"time"

	"github.com/app/server/domain"
)

const GameCollection string = "games"

type GameRepository struct {
	DB domain.IDatabase
}

// CreateGame Insert new Game document into the database
func (repo *GameRepository) CreateGame(Game *dbmodels.GameEnvironmentDB) error {
	Game.GameID = bson.NewObjectId()
	return repo.DB.Insert(GameCollection, Game)
}

// GetUsers Get list of games
func (repo *GameRepository) GetGames() dbmodels.GameEnvironmentsDB {
	Games := dbmodels.GameEnvironmentsDB{}
	err := repo.DB.FindAll(GameCollection, nil, Games, 50, "")
	if err != nil {
		return Games
	}
	return Games
}

// GetUsers Get list of users games
func (repo *GameRepository) GetGamesByUserId(id string) (dbmodels.GameEnvironmentsDB, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	Games := dbmodels.GameEnvironmentsDB{}
	err := repo.DB.FindAll(GameCollection, domain.Query{"_uid": bson.ObjectIdHex(id)}, &Games, 50, "")
	return Games, err
}

// GetUsers Get list of users games
func (repo *GameRepository) GetGamesByMapId(id string) (dbmodels.GameEnvironmentsDB, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	Games := dbmodels.GameEnvironmentsDB{}
	err := repo.DB.FindAll(GameCollection, domain.Query{"_mid": bson.ObjectIdHex(id)}, &Games, 50, "")
	return Games, err
}

// GetUser Get user specified by the id
func (repo *GameRepository) GetGameByGameId(id string) (*dbmodels.GameEnvironmentDB, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	var Game dbmodels.GameEnvironmentDB
	err := repo.DB.FindOne(GameCollection, domain.Query{"_id": bson.ObjectIdHex(id)}, &Game)
	return &Game, err
}

// GetUser Get user specified by the id
func (repo *GameRepository) GameExistsByUserVMapId(userid string, mapid string) (bool, error) {
	// TO DO : NEED TO MODIFT THIS
	if !bson.IsObjectIdHex(userid) || !bson.IsObjectIdHex(mapid) {
		return false, errors.New("Invalid ObjectId")
	}

	return repo.DB.Exists(GameCollection, domain.Query{"_uid": bson.ObjectIdHex(userid), "_mid": bson.ObjectIdHex(mapid)}), nil
}

func (repo *GameRepository) UpdateGame(id string, gUpdate *dbmodels.GameEnvironmentDB) (*dbmodels.GameEnvironmentDB, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	update := domain.Query{}
	if gUpdate.BelongingUserID.Hex() != "" {
		update["_uid"] = gUpdate.BelongingUserID
	}
	if gUpdate.BuiltOnMapID.Hex() != "" {
		update["_mid"] = gUpdate.BuiltOnMapID
	}
	if gUpdate.LastIterationTime != (time.Time{}) {
		update["gamelitime"] = gUpdate.LastIterationTime
	}
	if gUpdate.Man.BoughtGameUpgradeIDs != nil {
		update["upgrademanager"] = gUpdate.Man
	}
	if gUpdate.Resources != nil {
		update["resources"] = gUpdate.Resources
	}
	if gUpdate.Gm != (dbmodels.GameDB{}) {
		update["game"] = gUpdate.Gm
	}

	query := domain.Query{"_id": bson.ObjectIdHex(id)}
	change := domain.Change{
		Update:    domain.Query{"$set": update},
		ReturnNew: true,
	}

	var changedGameEnv dbmodels.GameEnvironmentDB
	err := repo.DB.Update(GameCollection, query, change, &changedGameEnv)
	//fmt.Println(err)
	//fmt.Println(changedSpace)
	return &changedGameEnv, err
}
