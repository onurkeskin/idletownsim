package domain

import (
	"github.com/app/helpers/version"
	"time"
)

type IGameMap interface {
	GetMapDate() time.Time
	SetMapDate(nd time.Time)
	GetMapVersion() version.Version
	SetMapVersion(Version version.Version)

	HasSpace(id string) bool
	AddSpace(s ISpace) error
	RemoveSpace(sid string) error

	PlaceBuildingInstance(sid string, buil ITileBuildableElement) error
	RemoveBuildingInstance(sid string) error
	SetBuildingInstance(sid string, buil ITileBuildableElement) error

	ForBuildables(func(b ITileBuildableElement))
	ForSpaces(f func(s ISpace))
}
