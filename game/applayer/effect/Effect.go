package effect

import (
	"fmt"
	"github.com/app/game/applayer/effect/domain"
)

type Effect struct {
	ID       string               `json:"id,omitempty" bson:"_id,omitempty"`
	Priority int64                `json:"priority" bson:"priority"`
	Issuer   domain.IEffectIssuer `json:"-" bson:"-"`
}

func (eff *Effect) GetUniqueID() string {
	return eff.ID
}

func (eff *Effect) ApplyEffectGlobal(anything interface{}) {
	return
}

func (eff *Effect) RemoveEffect() {
	return
}

func (eff *Effect) ReapplyEffect() {
	return
}

func (eff *Effect) GetPriority() int64 {
	return eff.Priority
}

func (eff *Effect) GetIssuer() domain.IEffectIssuer {
	return eff.Issuer
}

func (eff *Effect) SetIssuer(iss domain.IEffectIssuer) {
	eff.Issuer = iss
	return
}

func (g *Effect) String() string {
	str := ""
	str += fmt.Sprintf("id:%v ", g.ID)
	str += fmt.Sprintf("priority:%d ", g.Priority)
	str += fmt.Sprintf("issuer:%T", g.Issuer)
	return str
}
