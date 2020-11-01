package dbmodels

import (
	"encoding/json"
	"errors"
	"github.com/app/game/applayer/general"
	"github.com/app/game/applayer/general/domain"
)

func (v ValSchemeDB) FormIValScheme() domain.IValScheme {
	//FORM UPSCHEME
	switch v.SchemeType {
	case "mathvalscheme":
		var scheme general.MathValScheme
		json.Unmarshal(v.SchemeJson, &scheme)
		return &scheme
	default:
		return nil
	}
	return nil
}

func FormFromIValScheme(v domain.IValScheme) (*ValSchemeDB, error) {
	switch v := v.(type) {
	case *general.MathValScheme:
		jsonPart, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		toRet := ValSchemeDB{
			SchemeType: "mathvalscheme",
			SchemeJson: jsonPart,
		}
		return &toRet, nil
	default:
		return nil, errors.New("Cant understand the scheme type")
	}
	return nil, errors.New("Cant understand the scheme type")
}
