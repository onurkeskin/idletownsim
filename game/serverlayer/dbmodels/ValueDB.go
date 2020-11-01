package dbmodels

import (
	"errors"
	"github.com/app/game/applayer/general"
	"github.com/app/game/applayer/general/domain"
	"gopkg.in/mgo.v2/bson"
)

func (v ValueDB) FormIBProduct(r ResourceDB) domain.IBProduct {
	//FORM UPSCHEME
	switch r.ResType {
	case "regularbproduct":
		return general.NewRegularBProduct(v.ID.Hex(), v.ResourceCount)
	default:
		return nil
	}
}

func FormFromIBProduct(v domain.IBProduct) (*ValueDB, error) {
	switch v := v.(type) {
	case *general.RegularBProduct:
		toRet := ValueDB{
			bson.ObjectIdHex(v.GetType()),
			v.GetValue(),
		}
		return &toRet, nil
	default:
		return nil, errors.New("Cant understand the scheme type")
	}
	return nil, errors.New("Cant understand the scheme type")
}

func (vs ValuesDB) FormIBProducts(rs ResourcesDB) domain.IBProducts {
	//FORM UPSCHEME
	toRet := domain.IBProducts{}
	for _, v := range vs {
		for _, r := range rs {
			if r.ID == v.ID {
				bp := v.FormIBProduct(r)
				toRet = append(toRet, bp)
			}
		}
	}
	return toRet
}

func FormFromIBProducts(ibs domain.IBProducts) (ValuesDB, error) {
	toRet := ValuesDB{}
	for _, ib := range ibs {
		v, err := FormFromIBProduct(ib)
		if err != nil {
			return nil, err
		}
		toRet = append(toRet, *v)
	}
	return toRet, nil
}

/*
func (vs ValuesDB) FormIBProducts(rs ResourcesDB) domain.IBProducts {
	//FORM UPSCHEME
	toRet := domain.IBProducts{}
	for _, v := range vs {
		for _, r := range rs {
			if v.ID == r.ID {
				switch r.ResourceType {
				case "regularbproduct":
					toRet = append(toRet, general.NewRegularBProduct(v.ID.Hex(), v.ResourceCount))
				default:
					return nil
				}
			}
		}
	}
	return toRet
}
*/
/*
func (v ValuesDB) FormFromIBProducts(v domain.IBProducts) (ValuesDB, error) {
	//FORM UPSCHEME
	switch r.ResourceType {
	case "regularbproduct":
		return general.NewRegularBProduct(v.ID.Hex(), v.ResourceCount)
	default:
		return nil
	}
}
*/
