package domain

import (
//"fmt"
)

type IBProduct interface {
	GetType() string
	GetValue() float64

	SetType(string)
	SetValue(float64)

	Add(IBProduct) IBProduct
	Subtract(IBProduct) IBProduct

	Clone() IBProduct
}

type IBProducts []IBProduct

func (as *IBProducts) Compare(bs IBProducts) int {
	for _, a := range *as {
		found := false
		for _, b := range bs {
			if a.GetType() == b.GetType() {
				found = true
				if a.GetValue() > b.GetValue() {
					return -1
				}
			}
		}
		if !found {
			return -1
		}
	}

	return 1
}

func (a *IBProducts) Add(b IBProducts) {
	//fmt.Println("start")
	//fmt.Println(a)
	for _, bV := range b {
		elProcessed := false
		for _, aV := range *a {
			if aV.GetType() == bV.GetType() {
				aV.SetValue(bV.GetValue() + aV.GetValue())
				elProcessed = true
			}
		}
		if !elProcessed {
			*a = append(*a, bV)
		}
	}
	//fmt.Println(b)
	//fmt.Println(a)
}

func (a *IBProducts) Subtract(b IBProducts) {
	for _, bV := range b {
		elProcessed := false
		for _, aV := range *a {
			if aV.GetType() == bV.GetType() {
				aV.Subtract(bV)
				//aV.SetValue(aV.GetValue() - bV.GetValue())
				elProcessed = true
			}
		}
		if !elProcessed {
			*a = append(*a, bV)
		}
	}
}

func (a *IBProducts) Clone() interface{} {
	var toRet IBProducts
	for _, v := range *a {
		toRet = append(toRet, v.Clone())
	}
	return toRet
}
