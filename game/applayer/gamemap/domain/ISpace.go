package domain

import (
	generaldomain "github.com/app/game/applayer/general/domain"
)

type ISpace interface {
	generaldomain.IObjectProperties
	GetID() string
	GetInMapID() string

	GetOccupier() ITileBuildableElement

	AddElement(ITileBuildableElement) error
	RemoveElement() error
	SetElement(Element ITileBuildableElement) error

	SendElementChangedEvent(element ITileBuildableElement)
	AddElementChangeListener(listener ElementChangedListener) error
	RemoveElementChangeListener(listener ElementChangedListener) error

	AddSpaceMod(ISpaceEffect)
	RemoveSpaceMod(ISpaceEffect) bool
	ResetSpaceMods()

	GetNeighboorsAt(dir int) ([]ISpace, error)
	AddNeighboorTo(ns ISpace, dir int) error

	GetExpectedOutcome() generaldomain.IBProducts
	DoWork(total generaldomain.IBProducts) generaldomain.IBProducts
}

type ElementChangedListener interface {
	OnSpaceElementChange(element ITileBuildableElement)
}
