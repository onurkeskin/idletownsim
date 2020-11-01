package game

import (
	"encoding/json"
	"errors"
	"fmt"
	buildingdomain "github.com/app/game/applayer/building/domain"
	"github.com/app/game/applayer/buildinginstance"
	buildinginstancedomain "github.com/app/game/applayer/buildinginstance/domain"
	gamedomain "github.com/app/game/applayer/game/domain"
	generaldomain "github.com/app/game/applayer/general/domain"
	//upgrade "github.com/app/game/applayer/upgrade"

	"time"
)

type GameEnvironments []GameEnvironment

type GameEnvironment struct {
	GameID          string `json:"id,omitempty" bson:"_id,omitempty"`
	BelongingUserID string `json:"uid,omitempty" bson:"_uid,omitempty"`
	BuiltOnMapID    string `json:"mid,omitempty" bson:"_mid,omitempty"`

	GameCreationTime  time.Time `json:"gamect" bson:"gamect"`
	LastIterationTime time.Time `json:"gamelitime" bson:"gamelitime"`

	Resources generaldomain.IBProducts `json:"resources" bson:"resources"`
	Gm        gamedomain.IGame         `json:"game" bson:"game"`

	Man GameUpgradeManager `json:"-" bson:upgrademanager`
	//gameComplishment complishments
}

func NewGameEnvironment(
	GameID string,
	BelongingUserID string,
	BuiltOnMapID string,
	GameCreationTime time.Time,
	LastIterationTime time.Time,
	Resources generaldomain.IBProducts,
	BoughtGameUpgradeIDs []string,
	game gamedomain.IGame) gamedomain.IGameEnvironment {

	toRet := GameEnvironment{
		GameID:            GameID,
		BelongingUserID:   BelongingUserID,
		BuiltOnMapID:      BuiltOnMapID,
		Resources:         Resources,
		Gm:                game,
		GameCreationTime:  GameCreationTime,
		LastIterationTime: LastIterationTime,
	}
	man := GameUpgradeManager{BoughtGameUpgradeIDs: BoughtGameUpgradeIDs, connection: &toRet}
	toRet.Man = man
	return &toRet
}

func (d *GameEnvironment) MarshalJSON() ([]byte, error) {
	type toJsonGE GameEnvironment
	return json.Marshal(&struct {
		*toJsonGE
		GameCreationTime  string `json:"gamect" bson:"gamect"`
		LastIterationTime string `json:"gamelitime" bson:"gamelitime"`
	}{
		toJsonGE:          (*toJsonGE)(d),
		LastIterationTime: d.LastIterationTime.Format("2006-01-02T15:04:05.999999Z"),
		GameCreationTime:  d.GameCreationTime.Format("2006-01-02T15:04:05.999999Z"),
	})
}

func (g *GameEnvironment) GetGameID() string {
	return g.GameID
}
func (g *GameEnvironment) GetBelongingUserID() string {
	return g.BelongingUserID
}
func (g *GameEnvironment) GetMapID() string {
	return g.BuiltOnMapID
}

func (g *GameEnvironment) BuildingBuyable(building buildingdomain.IBuilding) bool {
	price := building.GetBasePrice()
	if price.Compare(g.Resources) > 0 {
		return false
	}
	return true
}

func (g *GameEnvironment) UpgradeBuyable(upgrade gamedomain.IUpgrade) error {
	price := upgrade.GetBasePrice()
	if price.Compare(g.Resources) > 0 {
		return errors.New("Not Enough Resources")
	}

	ok := upgrade.Eligible(g)
	if !ok {
		return errors.New("Requirements not Done")
	}

	return nil
}

func (g *GameEnvironment) BuyBuilding(building buildingdomain.IBuilding, sid string, options buildinginstance.FormFromBuildingOptions) (buildinginstancedomain.IBuildingInstance, error) {
	if g.BuildingBuyable(building) {
		return nil, errors.New("Not Enough Resources")
	}
	g.Resources.Subtract(building.GetBasePrice())

	bi := buildinginstance.FormFromBuilding(building, options)
	err := g.GetGame().PlaceBuildingInstance(sid, bi)
	if err != nil {
		return nil, err
	}
	return bi, err
}

func (g *GameEnvironment) BuyUpgrade(u gamedomain.IUpgrade) error {
	err := g.UpgradeBuyable(u)
	if err != nil {
		return err
	}
	g.Resources.Subtract(u.GetBasePrice())
	return g.Man.AddUpgrade(u)
}

func (g *GameEnvironment) AddUpgrade(u gamedomain.IUpgrade) error {
	return g.Man.AddUpgrade(u)
}

func (g *GameEnvironment) RemoveUpgrade(u gamedomain.IUpgrade) {
	g.Man.RemoveUpgrade(u)
}

func (g *GameEnvironment) PlayFor(t time.Duration) {
	med := g.Gm.TickFor(t)
	g.LastIterationTime = g.LastIterationTime.Add(t)
	g.Resources.Add(med)
}

func (g *GameEnvironment) GetResources() generaldomain.IBProducts {
	return g.Resources
}

func (g *GameEnvironment) GetGameStats() {

}

func (g *GameEnvironment) GetGame() gamedomain.IGame {
	return g.Gm
}

func (g *GameEnvironment) GetLastIterationTime() time.Time {
	return g.LastIterationTime
}
func (g *GameEnvironment) SetLastIterationTime(a time.Time) {
	g.LastIterationTime = a
}

func (g *GameEnvironment) GetCreationTime() time.Time {
	return g.GameCreationTime
}

func (g *GameEnvironment) SetCreationTime(a time.Time) {
	g.GameCreationTime = a
}

func (g *GameEnvironment) GetBoughtUpgrades() []string {
	return g.Man.BoughtGameUpgradeIDs
}

func (r *GameEnvironment) String() string {
	str := ""
	str += fmt.Sprintf("GameID:%s\n", r.GameID)
	str += fmt.Sprintf("BelongingUserID:%s\n", r.BelongingUserID)
	str += fmt.Sprintf("BuiltOnMapID:%s\n", r.BuiltOnMapID)
	str += fmt.Sprintf("Resources:%s\n", r.Resources)
	str += fmt.Sprintf("Game:%s\n", r.Gm)
	str += fmt.Sprintf("UpgradeManager:%v\n", r.Man)
	str += fmt.Sprintf("CreationTime:%v\n", r.GameCreationTime)
	return str
}
