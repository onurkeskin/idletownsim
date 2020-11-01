package domain

import (
	"github.com/app/game/applayer/general/domain"
)

type Generator interface {
	Generate(baseProducts domain.IBProducts) domain.IBProducts

	SetProduct(domain.IBProduct)
	GetProduct(string) domain.IBProduct
	RemoveProduct(string) bool

	AddIncrease(domain.IBProduct)
	//AddDecrease(IBProduct)

	GetProductionMods() []IGeneratorable
	AddProductionMod(IGeneratorable)
	RemoveProductionMod(IGeneratorable) bool
	Reset()
}
