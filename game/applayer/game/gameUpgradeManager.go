package game

import (
	"errors"
	"fmt"
	gamedomain "github.com/app/game/applayer/game/domain"
)

type GameUpgradeManager struct {
	connection gamedomain.IGameEnvironment

	BoughtGameUpgradeIDs []string `json:"-" bson:boughtupgrades`
	UpgradesApplied      []gamedomain.IUpgrade
}

func (gu *GameUpgradeManager) AddUpgrade(u gamedomain.IUpgrade) error {
	if u.Eligible(gu.connection) {
		u.ApplyUpgrade(gu.connection)
		// CLEAN THIS
		gu.BoughtGameUpgradeIDs = append(gu.BoughtGameUpgradeIDs, u.GetUniqueID())
	} else {
		return errors.New("Game does not have requirements")
	}
	return nil
}

func (gu *GameUpgradeManager) RemoveUpgrade(u gamedomain.IUpgrade) {
	exists := false
	for _, v := range gu.BoughtGameUpgradeIDs {
		if v == u.GetUniqueID() {
			exists = true
		}
	}
	if exists {
		u.RemoveUpgrade()
	}
}

func (r *GameUpgradeManager) String() string {
	str := ""
	str += fmt.Sprintf("BoughtUpgrades:%s\n", r.BoughtGameUpgradeIDs)
	str += fmt.Sprintf("AppliedUpgrades:%s\n", r.UpgradesApplied)
	return str
}
