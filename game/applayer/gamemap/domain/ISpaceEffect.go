package domain

import (
	effectdomain "github.com/app/game/applayer/effect/domain"
	generaldomain "github.com/app/game/applayer/general/domain"
	generatordomain "github.com/app/game/applayer/generator/domain"
)

type ISpaceEffect interface {
	ApplyEffect(s ISpace)
	RemoveEffectFrom(s ISpace)

	effectdomain.IEffect
	generatordomain.IGeneratorable
	generaldomain.ICloneable
}
