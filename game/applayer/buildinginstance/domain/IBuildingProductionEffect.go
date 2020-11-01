package domain

import (
	effectdomain "github.com/app/game/applayer/effect/domain"
	generatordomain "github.com/app/game/applayer/generator/domain"
)

type IBuildingProductionEffect interface {
	ApplyEffect(b IBuildingInstance)
	RemoveEffectFrom(b IBuildingInstance)

	generatordomain.IGeneratorable
	effectdomain.IEffect
}
