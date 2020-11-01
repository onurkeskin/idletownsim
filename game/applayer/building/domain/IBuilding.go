package domain

import (
	gamemapdomain "github.com/app/game/applayer/gamemap/domain"
	generaldomain "github.com/app/game/applayer/general/domain"
)

type IBuilding interface {
	generaldomain.IObjectProperties
	GetID() string

	GetBasePrice() generaldomain.IBProducts

	GetUpValScheme() generaldomain.IValScheme
	GetUpStatScheme() generaldomain.IValScheme

	GetBaseProducts() generaldomain.IBProducts
	GetProductionIntervalNano() int64

	GetSpaceEffect() gamemapdomain.ISpaceEffect
	SetSpaceEffect(eff gamemapdomain.ISpaceEffect) bool
	/*
		GetVicinity(direction int)
		SetVicinity(b IBuilding, dir int)
	*/

}

type IBuildings interface{}
