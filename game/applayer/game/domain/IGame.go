package domain

import (
	buildinginstancedomain "github.com/app/game/applayer/buildinginstance/domain"
	gamemapdomain "github.com/app/game/applayer/gamemap/domain"
	generaldomain "github.com/app/game/applayer/general/domain"
	"time"
)

type IGame interface {
	PlaceBuildingInstance(sid string, buil buildinginstancedomain.IBuildingInstance) error
	RemoveBuildingInstance(sid string) error

	AddSpace(s gamemapdomain.ISpace) error
	RemoveSpace(sid string) error
	GetGameMap() gamemapdomain.IGameMap

	AddBuildingProductModifier(props generaldomain.IObjectProperties, chkType generaldomain.CheckType, eff buildinginstancedomain.IBuildingProductionEffect)
	RemoveBuildingProductModifier(eff buildinginstancedomain.IBuildingProductionEffect) bool

	AddSpaceModifier(props generaldomain.IObjectProperties, chkType generaldomain.CheckType, eff gamemapdomain.ISpaceEffect)
	RemoveSpaceModifier(eff gamemapdomain.ISpaceEffect) bool
	RemoveSpaceModifierFrom(eff gamemapdomain.ISpaceEffect, s gamemapdomain.ISpace) bool

	TickFor(tick time.Duration) generaldomain.IBProducts
	//CountSpaceWith(f func(gamemapdomain.ISpace) bool) int
	GetBuildablesWithProperty(props generaldomain.IObjectProperties, chk generaldomain.CheckType) []gamemapdomain.ITileBuildableElement
	GetSpacesWithProperty(props generaldomain.IObjectProperties, chk generaldomain.CheckType) []gamemapdomain.ISpace
}
