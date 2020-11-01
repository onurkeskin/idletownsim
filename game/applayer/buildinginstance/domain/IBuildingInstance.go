package domain

import (
	effectdomain "github.com/app/game/applayer/effect/domain"
	gamemapdomain "github.com/app/game/applayer/gamemap/domain"
	generaldomain "github.com/app/game/applayer/general/domain"
	"time"
)

type IBuildingInstance interface {
	effectdomain.IEffectIssuer
	gamemapdomain.ITileBuildableElement
	GetUniqueID() string
	GetParentID() string
	//	GetUniqueID() string
	GetLevel() int

	PredictWorkTime(t int64) int64
	DoWork(d time.Duration) []generaldomain.IBProducts
	GetExpectedOutcome() generaldomain.IBProducts

	GetBaseProducts() generaldomain.IBProducts
	GetBaseProductionIntervalNano() int64

	GetLastProductionTimeUnix() time.Time
	SetLastProductionTimeUnix(time.Time)
	GetBuiltTimeUnix() time.Time
	SetBuiltTimeUnix(a time.Time)

	GetProductMods() []IBuildingProductionEffect
	AddProductMods(IBuildingProductionEffect)
	RemoveProductMods(IBuildingProductionEffect)
	ResetProductMods()

	BIPropertyChangedEvent()
	AddBIPropertyChangeListener(listener BIPropertyChangedListener) error
	RemoveBIPropertyChangeListener(listener BIPropertyChangedListener) error
}

type BIPropertyChangedListener interface {
	OnBIPropertyChange(IBuildingInstance)
}
