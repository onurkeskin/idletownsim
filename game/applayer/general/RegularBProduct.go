package general

import (
	"fmt"
	"github.com/app/game/applayer/general/domain"
)

type RegularBProduct struct {
	A string  `json:"restype" bson:"restype"`
	B float64 `json:"resvalue" bson:"resvalue"`
}

func NewRegularBProduct(
	a string,
	b float64) *RegularBProduct {
	return &RegularBProduct{A: a, B: b}
}

func (a *RegularBProduct) SetType(typ string) {
	a.A = typ
}

func (a *RegularBProduct) SetValue(val float64) {
	a.B = val
}

func (a *RegularBProduct) GetType() string {
	return a.A
}

func (a *RegularBProduct) GetValue() float64 {
	return a.B
}

func (a *RegularBProduct) Add(b domain.IBProduct) domain.IBProduct {
	if a.GetType() == b.GetType() {
		a.SetValue(a.GetValue() + b.GetValue())
	}
	return a
}

func (a *RegularBProduct) Subtract(b domain.IBProduct) domain.IBProduct {
	//fmt.Println(a.GetValue(), b.GetValue())
	if a.GetType() == b.GetType() {
		a.SetValue(a.GetValue() - b.GetValue())
	}
	return a
}

func (a *RegularBProduct) Clone() domain.IBProduct {
	return NewRegularBProduct(a.A, a.B)
}

func (a *RegularBProduct) String() string {
	return fmt.Sprintf("type:[%s],count:[%f] ", a.A, a.B)
}
