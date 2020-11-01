package dbmodels

import (
	"github.com/app/helpers/version"
	"time"

	//"fmt"
	"errors"
	"github.com/app/game/applayer/buildinginstance"
	"github.com/app/game/applayer/gamemap/domain"
)

type TileBuildableElementDB struct {
	/*
		BuildableVersion version.Version `json:"version,omitempty" bson:"version,omitempty"`

		BuildableType string                              `json:"btype,omitempty" bson:"btype,omitempty"`
		Buildable     gamemapdomain.ITileBuildableElement `json:"buildable,omitempty" bson:"buildable,omitempty"`
		EffectRaw bson.M `json:"-" bson:"effectraw,omitempty"`
	*/
	UniqueId                 string          `bson:"_id"`
	ParentId                 string          `bson:"ParentId"`
	BuildingGlobalIdentifier string          `bson:"buildingglobalidentifier"`
	Level                    int             `bson:"bLevel"`
	LastProductionTime       time.Time       `bson:"prodtime"`
	BuiltTime                time.Time       `bson:"builttime"`
	SpaceEffect              EffectDB        `json:"buildingeffect" bson:"buildingeffect"`
	BuildableVersion         version.Version `json:"version,omitempty" bson:"version,omitempty"`
}

func (buildable TileBuildableElementDB) FormITileBuildableElement() domain.ITileBuildableElement {
	uid := buildable.UniqueId
	pid := buildable.ParentId
	Level := buildable.Level
	LastProductionTime := buildable.LastProductionTime
	BuiltTime := buildable.BuiltTime

	SpaceEffect, err := buildable.SpaceEffect.FormEffect()
	if err != nil {
		panic(err)
	}

	buil := buildinginstance.NewBuildingInstance(
		uid,
		pid,
		Level,
		nil,
		0,
		BuiltTime,
		nil,
		nil,
	)
	buil.SetLastProductionTimeUnix(LastProductionTime)

	if SpaceEffect != nil {
		castSE := SpaceEffect.(domain.ISpaceEffect)
		buil.AddSpaceEffect(castSE)
	}
	return buil
}

func FormFromITileBuildableElement(v domain.ITileBuildableElement) (*TileBuildableElementDB, error) {
	toRet := TileBuildableElementDB{}

	switch v := v.(type) {
	case *buildinginstance.BuildingInstance:
		toRet.UniqueId = v.UniqueId
		toRet.ParentId = v.ParentId
		toRet.Level = v.Level
		toRet.LastProductionTime = v.LastProductionTime
		toRet.BuiltTime = v.BuiltTime
	default:
		return nil, errors.New("Cant understand the scheme type")
	}

	if v.GetSpaceEffect() == nil {
		return &toRet, nil
	}

	eDb, err := FormFromEffect(v.GetSpaceEffect())
	if err != nil {
		panic(err)
	}
	toRet.SpaceEffect = *eDb
	return &toRet, nil
}

/*
// GetBSON implements bson.Getter.
func (v *TileBuildableElementDB) GetBSON() (interface{}, error) {
	return v, nil
}
*/
// SetBSON implements bson.Setter.
/*
func (v *TileBuildableElementDB) SetBSON(raw bson.Raw) error {
	if raw.Data == nil || len(raw.Data) == 0 {
		//fmt.Println("err")
		return nil
	}

	decoded := new(struct {
		BuildableType    string          `json:"btype,omitempty" bson:"btype,omitempty"`
		EffectRaw        bson.M          `json:"-" bson:"effectraw,omitempty"`
		BuildableVersion version.Version `json:"version,omitempty" bson:"version,omitempty"`
	})
	//fmt.Println("start")
	bsonErr := raw.Unmarshal(decoded)

	if bsonErr == nil {
		v.BuildableType = decoded.BuildableType
		v.BuildableVersion = decoded.BuildableVersion
		v.EffectRaw = decoded.EffectRaw
	} else {
		//fmt.Println("err")
		return bsonErr
	}

	switch v.BuildableType {
	case "buildinginstance":
		decodeBuildable := new(struct {
			Buildable *buildinginstance.BuildingInstance `json:"buildable,omitempty" bson:"buildable,omitempty"`
		})
		bsonErr = raw.Unmarshal(decodeBuildable)
		if bsonErr == nil {
			v.Buildable = decodeBuildable.Buildable
			//fmt.Println("ret")
		} else {
			//fmt.Println("err")
			return bsonErr
		}
	default:
		//fmt.Println("err")
		return errors.New("unrecognized buildable type")
	}

	for k, bs := range v.EffectRaw {
		switch k {
		case "centeredareaspaceeffect":
			var effect gamemap.CenteredAreaSpaceEffect
			err := bson.Unmarshal(GetBytes(bs), &effect)
			if err != nil {
				panic(err)
			}
			v.Buildable.AddSpaceEffect(&effect)
		case "spaceeffect":
			var effect gamemap.SpaceEffect
			err := bson.Unmarshal(GetBytes(bs), &effect)
			if err != nil {
				panic(err)
			}
			v.Buildable.AddSpaceEffect(&effect)
		}
	}
	return nil
}

func GetBytes(key interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		panic(err)
		return nil
	}
	return buf.Bytes()
}
*/
