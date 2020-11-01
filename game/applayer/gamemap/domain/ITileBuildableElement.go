package domain

import (
	generaldomain "github.com/app/game/applayer/general/domain"
)

type ITileBuildableElement interface {
	generaldomain.IObjectProperties

	AddSpaceEffect(e ISpaceEffect)
	GetSpaceEffect() ISpaceEffect
	ApplySpaceEffect(s ISpace) error
	RemoveSpaceEffect() error
}
