package general

import (
	"github.com/app/game/applayer/general/domain"

	"errors"
	"fmt"
)

type ObjectProperties []ObjectProperty

type ObjectProperty struct {
	ParamName  string `json:"pname" bson:"pname"`
	ParamValue string `json:"pvalue" bson:"pvalue"`
}

func NewObjectProperty(
	ParamName,
	ParamValue string) *ObjectProperty {
	return &ObjectProperty{
		ParamName:  ParamName,
		ParamValue: ParamValue,
	}
}

func (o *ObjectProperty) GetParamName() string {
	return o.ParamName
}

func (o *ObjectProperty) GetParamValue() string {
	return o.ParamValue
}

func NewObjectProperties() *ObjectProperties {
	return &ObjectProperties{}
}

func (o *ObjectProperties) PropertiesSlice() []domain.IObjectProperty {
	toRet := []domain.IObjectProperty{}
	for _, v := range *o {
		toRet = append(toRet, &v)
	}
	return toRet
}

func (o *ObjectProperties) GetProperty(paramName string) ([]string, error) {
	toRet := []string{}
	foundOne := false
	for _, v := range *o {
		if v.GetParamName() == paramName {
			toRet = append(toRet, v.GetParamValue())
			foundOne = true
		}
	}
	if !foundOne {
		return nil, errors.New(fmt.Sprintf("Property not found with name %s", paramName))
	}

	return toRet, nil
}
func (o *ObjectProperties) Count() int {
	return len(*o)
}

func (o *ObjectProperties) AddProperty(paramName, paramValue string) error {
	for _, v := range *o {
		if v.GetParamName() == paramName && v.GetParamValue() == paramValue {
			return errors.New(fmt.Sprintf("Already has parameter %s with value %s", v.GetParamName(), v.GetParamValue()))
		}
	}
	*o = append(*o, *NewObjectProperty(paramName, paramValue))

	return nil
}
func (o *ObjectProperties) RemoveProperty(paramName, paramValue string) error {
	for in, v := range *o {
		if v.GetParamName() == paramName && v.GetParamValue() == paramValue {
			*o = append((*o)[:in], (*o)[in+1:]...)
		}
	}

	return errors.New(fmt.Sprintf("Doesnt have parameter %s with value %s", paramName, paramValue))
}

func (o *ObjectProperties) SatisfiesType(check domain.IObjectProperties, chkType domain.CheckType) bool {
	if len(*o) == 0 {
		return false
	}

	if chkType == domain.CheckTypeExact {
		if len(*o) != check.Count() {
			return false
		}
	}

	for _, mainProp := range *o {
		gotEqualParamName := false
		for _, checkProp := range check.PropertiesSlice() {

			if mainProp.GetParamName() == checkProp.GetParamName() {
				gotEqualParamName = true

				if mainProp.GetParamValue() == checkProp.GetParamValue() {
					if chkType == domain.CheckTypeAny {
						return true
					}
				} else {
					if chkType == domain.CheckTypeAll || chkType == domain.CheckTypeExact {
						return false
					}
				}

			}
		}
		if !gotEqualParamName {
			if chkType == domain.CheckTypeAll || chkType == domain.CheckTypeExact {
				return false
			}
		}
	}

	return true
}

func (o *ObjectProperty) String() string {
	str := ""
	str += fmt.Sprintf("Param Name:%s, Param Value:%s", o.ParamName, o.ParamValue) //TODO IMPLEMENT
	return str
}
