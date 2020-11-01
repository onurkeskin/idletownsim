package game

import (
	"errors"
	"fmt"

	gamedomain "github.com/app/game/applayer/game/domain"
	//gamemapdomain "github.com/app/game/applayer/gamemap/domain"
	general "github.com/app/game/applayer/general"
	generaldomain "github.com/app/game/applayer/general/domain"
)

type ComplishmentTarget struct {
	*general.ObjectProperties
}

type ComplishmentParams struct {
	InequalitySymbol string
	Value            int
}

type Complishment struct {
	ID     string               `json:"id,omitempty" bson:"_id,omitempty"`
	Target []ComplishmentTarget `json:"complishmenttargets" bson:"complishmenttargets"`
	Params []ComplishmentParams `json:"complishmentparams" bson:"complishmentparams"`
}

func (c Complishment) getTargetNum(i int) generaldomain.IObjectProperties {
	return (c.Target[i].ObjectProperties)
}

func (c *Complishment) GetTarget() []ComplishmentTarget {
	return c.Target
}

func (a *Complishment) ComplishmentParse() (func(g gamedomain.IGameEnvironment) bool, error) {
	switch a.getTargetNum(0).PropertiesSlice()[0].GetParamName() {
	case "building":
		return a.parseForBuildingType(), nil
	case "space":
		return a.parseForSpaceType(), nil
	case "product":
		return a.parseForProductType(), nil
	}
	return nil, errors.New("cant be parsed")
}

func InsideParams(Params []ComplishmentParams, Value int) bool {
	for _, v := range Params {
		switch v.InequalitySymbol {
		case "<":
			if Value >= v.Value {
				return false
			}
		case "<=":
			if Value > v.Value {
				return false
			}
		case ">":
			if Value <= v.Value {
				return false
			}
		case ">=":
			if Value < v.Value {
				return false
			}
		case "=":
			if Value != v.Value {
				return false
			}
		}
	}
	return true
}

func (a *Complishment) parseForBuildingType() func(g gamedomain.IGameEnvironment) bool {

	if len(a.Target) == 1 {
		BType := a.getTargetNum(0)
		return func(g gamedomain.IGameEnvironment) bool {
			bs := g.GetGame().GetBuildablesWithProperty(BType, generaldomain.CheckTypeAll)
			return InsideParams(a.Params, len(bs))
		}
	}

	if len(a.Target) == 2 {
		BP := a.getTargetNum(1)
		BType := a.getTargetNum(0)
		return func(g gamedomain.IGameEnvironment) bool {
			bs := g.GetGame().GetBuildablesWithProperty(BType, generaldomain.CheckTypeAll)
			count := 0
			for _, b := range bs {
				if b.SatisfiesType(BP, generaldomain.CheckTypeExact) {
					count++
				}
			}

			return InsideParams(a.Params, count)
		}
	}

	return nil
}

func (a *Complishment) parseForSpaceType() func(g gamedomain.IGameEnvironment) bool {
	return nil
}

func (a *Complishment) parseForProductType() func(g gamedomain.IGameEnvironment) bool {
	return nil
}

type ComplishmentTestFunc func(g gamedomain.IGameEnvironment) bool

func (c ComplishmentTestFunc) TestObject(g gamedomain.IGameEnvironment) bool {
	return c(g)
}

func (c *Complishment) String() string {
	str := ""
	str += fmt.Sprintf("comp Id:%s, Comps:%v, Params:%v", c.ID, c.Target, c.Params) //TODO IMPLEMENT
	return str
}

func (c ComplishmentTarget) String() string {
	str := ""
	str += fmt.Sprintf("%v", c.ObjectProperties) //TODO IMPLEMENT
	return str
}
