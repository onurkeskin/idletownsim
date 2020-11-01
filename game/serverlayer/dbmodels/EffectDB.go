package dbmodels

import (
	"encoding/json"
	"errors"
	"github.com/app/game/applayer/buildinginstance"
	"github.com/app/game/applayer/effect/domain"
	"github.com/app/game/applayer/gamemap"
)

func (v EffectDB) FormEffect() (domain.IEffect, error) {
	//fmt.Println(v)
	//FORM UPSCHEME
	switch v.EffectType {
	case "buildingproductioneffect":
		var scheme buildinginstance.BuildingProductionEffect
		err := json.Unmarshal(v.EffectJson, &scheme)
		if err != nil {
			return nil, err
		}
		return &scheme, nil
	case "spaceeffect":
		var scheme gamemap.SpaceEffect
		err := json.Unmarshal(v.EffectJson, &scheme)
		if err != nil {
			return nil, err
		}
		return &scheme, nil
	case "centeredareaspaceeffect":
		var scheme gamemap.CenteredAreaSpaceEffect
		err := json.Unmarshal(v.EffectJson, &scheme)
		if err != nil {
			return nil, err
		}
		//pretty.Println(scheme)
		return &scheme, nil
	default:
		return nil, errors.New("Unkown effect type")
	}
	return nil, errors.New("Unkown effect type")
}

func (v EffectsDB) FormEffects() ([]domain.IEffect, error) {
	toRet := []domain.IEffect{}
	for _, _eff := range v {
		eff, err := _eff.FormEffect()
		if err != nil {
			return nil, err
		}
		toRet = append(toRet, eff)
	}
	return toRet, nil
}

func FormFromEffect(v domain.IEffect) (*EffectDB, error) {
	switch v := v.(type) {
	case *buildinginstance.BuildingProductionEffect:
		jsonPart, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		toRet := EffectDB{
			EffectType: "buildingproductioneffect",
			EffectJson: jsonPart,
		}
		return &toRet, nil
	case *gamemap.SpaceEffect:
		jsonPart, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		toRet := EffectDB{
			EffectType: "spaceeffect",
			EffectJson: jsonPart,
		}
		return &toRet, nil
	case *gamemap.CenteredAreaSpaceEffect:
		jsonPart, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		toRet := EffectDB{
			EffectType: "centeredareaspaceeffect",
			EffectJson: jsonPart,
		}
		return &toRet, nil
	default:
		return nil, errors.New("Cant understand the scheme type")
	}
	return nil, errors.New("Cant understand the scheme type")
}
