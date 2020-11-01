package domain

import (
	effectdomain "github.com/app/game/applayer/effect/domain"
	generaldomain "github.com/app/game/applayer/general/domain"
)

type IUpgrade interface {
	effectdomain.IEffectIssuer

	GetBasePrice() generaldomain.IBProducts
	GetUniqueID() string
	ApplyUpgrade(env IGameEnvironment) bool
	RemoveUpgrade() error

	Eligible(env IGameEnvironment) bool
}
type IUpgrades interface{}
