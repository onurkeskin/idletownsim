package game

import (
	"fmt"
	gamedomain "github.com/app/game/applayer/game/domain"
)

type Requirement struct {
	ID                        string                 `json:"id,omitempty" bson:"_id,omitempty"`
	RequiredComplishments     []Complishment         `json:"complishmentss,omitempty" bson:"complishmentss,omitempty"`
	RequiredComplishmentFuncs []ComplishmentTestFunc `json:"-" bson:"-"`
}

func NewRequirement(
	ID string,
	RequiredComplishments []Complishment) *Requirement {
	reqs := []Complishment{}
	funs := []ComplishmentTestFunc{}
	for _, _v := range RequiredComplishments {
		v, err := _v.ComplishmentParse()
		if err == nil {
			funs = append(funs, v)
			reqs = append(reqs, _v)
		}
	}
	return &Requirement{
		ID: ID,
		RequiredComplishments:     reqs,
		RequiredComplishmentFuncs: funs,
	}
}

func (r *Requirement) Satisfy(g gamedomain.IGameEnvironment) bool {
	for _, v := range r.RequiredComplishmentFuncs {
		res := v.TestObject(g)
		if !res {
			return false
		}
	}
	return true
}

func (r *Requirement) String() string {
	str := ""
	str += fmt.Sprintf("Req Id:%s, Comps:%s", r.ID, r.RequiredComplishments) //TODO IMPLEMENT
	return str
}
