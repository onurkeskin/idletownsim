package buildinginstance

import (
	"fmt"
	buildinginstancedomain "github.com/app/game/applayer/buildinginstance/domain"
	effect "github.com/app/game/applayer/effect"
	effectdomain "github.com/app/game/applayer/effect/domain"
	general "github.com/app/game/applayer/general"
	generaldomain "github.com/app/game/applayer/general/domain"
)

type BuildingProductionEffect struct {
	effect.GenerationEffect `json:"generationeffect" bson:"generationeffect"`

	AppliedTo []*buildinginstancedomain.IBuildingInstance `json:"-" bson:"-"`
}

func NewBuildingProductionEffect(
	ID string,
	priority int64,
	targetResources []string,
	scheme *general.MathValScheme,
	issuer effectdomain.IEffectIssuer) *BuildingProductionEffect {

	return &BuildingProductionEffect{
		GenerationEffect: effect.GenerationEffect{
			Effect: effect.Effect{
				ID:       ID,
				Priority: priority,
				Issuer:   issuer},
			Scheme:          scheme,
			TargetResources: targetResources,
		},
		AppliedTo: []*buildinginstancedomain.IBuildingInstance{}}
}

func (eff *BuildingProductionEffect) ApplyEffectGlobal(anything interface{}) {
	v, ok := anything.(buildinginstancedomain.IBuildingInstance)
	if !ok {
		return
	}

	eff.ApplyEffect(v)
}

func (eff *BuildingProductionEffect) ApplyEffect(b buildinginstancedomain.IBuildingInstance) {
	b.AddProductMods(eff)
	eff.AppliedTo = append(eff.AppliedTo, &b)
	return
}
func (eff *BuildingProductionEffect) RemoveEffect() {
	for _, b := range eff.AppliedTo {
		(*b).RemoveProductMods(eff)
	}
	eff.AppliedTo = nil
}

func (eff *BuildingProductionEffect) ReapplyEffect() {
	toApply := eff.AppliedTo
	eff.RemoveEffect()

	for _, s := range toApply {
		eff.ApplyEffect(*s)
	}
}

func (eff *BuildingProductionEffect) RemoveEffectFrom(b buildinginstancedomain.IBuildingInstance) {
	for _, v := range eff.AppliedTo {
		if *v == b {
			b.RemoveProductMods(eff)
		}
	}
	return
}

func (eff *BuildingProductionEffect) GetTargetResources() []string {
	return eff.TargetResources
}

func (eff *BuildingProductionEffect) GetScheme() generaldomain.IValScheme {
	return eff.Scheme
}

func (g *BuildingProductionEffect) String() string {
	str := ""
	str += fmt.Sprintf("GenerationEffect:%s ", g.GenerationEffect.String())
	return str
}
