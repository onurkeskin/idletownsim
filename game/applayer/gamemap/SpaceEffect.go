package gamemap

import (
	"fmt"

	effect "github.com/app/game/applayer/effect"
	effectdomain "github.com/app/game/applayer/effect/domain"
	gamemapdomain "github.com/app/game/applayer/gamemap/domain"
	general "github.com/app/game/applayer/general"
	generaldomain "github.com/app/game/applayer/general/domain"
)

type SpaceEffect struct {
	effect.GenerationEffect `json:"generationeffect" bson:"generationeffect"`

	AppliedTo []*gamemapdomain.ISpace `json:"-" bson:"-"`
}

func NewSpaceEffect(
	ID string,
	Priority int64,
	TargetResources []string,
	Scheme *general.MathValScheme,
	Issuer effectdomain.IEffectIssuer) gamemapdomain.ISpaceEffect {
	return &SpaceEffect{
		GenerationEffect: effect.GenerationEffect{
			Effect: effect.Effect{
				ID:       ID,
				Priority: Priority,
				Issuer:   Issuer},
			Scheme:          Scheme,
			TargetResources: TargetResources,
		},
		AppliedTo: []*gamemapdomain.ISpace{}}
}

func (eff *SpaceEffect) Clone() interface{} {
	newSliceTargetRes := make([]string, len(eff.TargetResources))
	copy(newSliceTargetRes, eff.TargetResources)
	return NewSpaceEffect(eff.ID, eff.Priority, newSliceTargetRes, eff.Scheme, eff.Issuer)
}

func (eff *SpaceEffect) ApplyEffectGlobal(anything interface{}) {
	v, ok := anything.(gamemapdomain.ISpace)
	if !ok {
		return
	}

	eff.ApplyEffect(v)
}

func (eff *SpaceEffect) ApplyEffect(s gamemapdomain.ISpace) {
	s.AddSpaceMod(eff)
	eff.AppliedTo = append(eff.AppliedTo, &s)
	return
}

func (eff *SpaceEffect) RemoveEffect() {
	for _, s := range eff.AppliedTo {
		(*s).RemoveSpaceMod(eff)
	}
	eff.AppliedTo = []*gamemapdomain.ISpace{}
}

func (eff *SpaceEffect) ReapplyEffect() {
	toApply := eff.AppliedTo
	eff.RemoveEffect()

	for _, s := range toApply {
		eff.ApplyEffect(*s)
	}
}

func (eff *SpaceEffect) RemoveEffectFrom(space gamemapdomain.ISpace) {
	for in, s := range eff.AppliedTo {
		if *s == space {
			(*s).RemoveSpaceMod(eff)
			eff.AppliedTo = append(eff.AppliedTo[:in], eff.AppliedTo[in+1:]...)
		}
	}
}

func (eff *SpaceEffect) GetTargetResources() []string {
	return eff.TargetResources
}

func (eff *SpaceEffect) GetScheme() generaldomain.IValScheme {
	return eff.Scheme
}

func (g *SpaceEffect) String() string {
	str := ""
	str += fmt.Sprintf("GenerationEffect:%s ", g.GenerationEffect.String())
	return str
}
