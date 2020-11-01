package generator

import (
	"fmt"

	generaldomain "github.com/app/game/applayer/general/domain"
	generatordomain "github.com/app/game/applayer/generator/domain"
	dtstructers "github.com/datastructures/datastructures"
)

type BasicGenerator struct {
	additiveProductions generaldomain.IBProducts
	bProductMods        *dtstructers.PQueueList
	//bProductMods []domain.IGeneratorable
}

func NewBasicGenerator() *BasicGenerator {

	return &BasicGenerator{generaldomain.IBProducts{}, dtstructers.NewPQueueList(dtstructers.MAXPQ)}
}

func (a *BasicGenerator) baseProductsUpdated() {

}

func (a *BasicGenerator) SetProduct(p generaldomain.IBProduct) {
	for _, el := range a.additiveProductions {
		if el.GetType() == p.GetType() {
			el.SetValue(p.GetValue())
			return
		}
	}
	a.additiveProductions = append(a.additiveProductions, p)
	return
}

func (a *BasicGenerator) GetProduct(t string) generaldomain.IBProduct {
	for _, el := range a.additiveProductions {
		if el.GetType() == t {
			return el
		}
	}
	return nil
}
func (a *BasicGenerator) RemoveProduct(t string) bool {
	for i, el := range a.additiveProductions {
		if el.GetType() == t {
			a.additiveProductions = append(a.additiveProductions[:i], a.additiveProductions[i+1:]...)
		}
	}
	return true
}

func (a *BasicGenerator) AddIncrease(p generaldomain.IBProduct) {
	for _, el := range a.additiveProductions {
		if el.GetType() == p.GetType() {
			el.SetValue(el.GetValue() + p.GetValue())
			return
		}
	}
	a.additiveProductions = append(a.additiveProductions, p)
	return
}

func (a *BasicGenerator) Reset() {
	//	a.mockDecrease = domain.IBProducts{}
	a.bProductMods = dtstructers.NewPQueueList(dtstructers.MAXPQ)
}

func (a *BasicGenerator) GetProductionMods() []generatordomain.IGeneratorable {
	toReturn := []generatordomain.IGeneratorable{}
	for v, next := a.bProductMods.Iterate()(); next != nil; v, next = next() {
		eff := v.(generatordomain.IGeneratorable)
		toReturn = append(toReturn, eff)
	}
	return toReturn
}

func (a *BasicGenerator) AddProductionMod(eff generatordomain.IGeneratorable) {
	//a.bProductMods = append(a.bProductMods, eff)
	a.bProductMods.Push(eff, eff.GetPriority())
	return
}

func (a *BasicGenerator) RemoveProductionMod(eff generatordomain.IGeneratorable) bool {
	comparer := func(a, b interface{}) bool {
		return a == b
	}

	a.bProductMods.Remove(eff, comparer)
	//fmt.Printf("%p\n", asda)
	return false
}

func (a *BasicGenerator) Generate(baseProducts generaldomain.IBProducts) generaldomain.IBProducts {
	base := baseProducts.Clone().(generaldomain.IBProducts)
	base.Add(a.additiveProductions)
	for v, next := a.bProductMods.Iterate()(); next != nil; v, next = next() {
		eff := v.(generatordomain.IGeneratorable)
		for _, val1 := range eff.GetTargetResources() {
			for _, bproduct := range base {
				if bproduct.GetType() == val1 {
					val, err := eff.GetScheme().CalculateValue(bproduct.GetValue())
					if err {
						// TODO
					} else {
						bproduct.SetValue(val.(float64))
					}
				}
			}
		}
	}
	return base
}

func (a *BasicGenerator) String() string {
	str := ""
	str += fmt.Sprintf("Current Additive Products:%s\n", a.additiveProductions)
	str += fmt.Sprintf("Current Multiplier Effects:%s", a.bProductMods)
	return str
}
