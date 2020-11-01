package game

import (
	"errors"
	"fmt"
	buildinginstancedomain "github.com/app/game/applayer/buildinginstance/domain"

	effectdomain "github.com/app/game/applayer/effect/domain"
	gamedomain "github.com/app/game/applayer/game/domain"
	general "github.com/app/game/applayer/general"
	generaldomain "github.com/app/game/applayer/general/domain"
)

type UpgradeTarget struct {
	*general.ObjectProperties
}

func (c UpgradeTarget) String() string {
	str := ""
	str += fmt.Sprintf("%v", c.ObjectProperties) //TODO IMPLEMENT
	return str
}

type Upgrades []Upgrade

type Upgrade struct {
	ID                 string                   `json:"id,omitempty" bson:"_id,omitempty"`
	UpgradeID          string                   `json:"uid" bson:"uid"`
	UpgradeRequirement *Requirement             `json:"req" bson:"req"`
	Price              generaldomain.IBProducts `json:"price" bson:"price"`
	UpTargets          []UpgradeTarget          `json:"targets" bson:"targets"`
	Effects            []effectdomain.IEffect   `json:"effects" bson:"effects"`
}

func NewUpgrade(
	UpgradeID string,
	Price generaldomain.IBProducts,
	UpgradeRequirement *Requirement,
	UpTargets []UpgradeTarget,
	Effects []effectdomain.IEffect) *Upgrade {

	u := Upgrade{
		UpgradeID:          UpgradeID,
		Price:              Price,
		UpgradeRequirement: UpgradeRequirement,
		UpTargets:          UpTargets,
		Effects:            Effects}
	for _, v := range Effects {
		v.SetIssuer(&u)
	}
	return &u
}

func (u *Upgrade) GetDeployedEffects() []effectdomain.IEffect {
	return u.Effects
}

func (u *Upgrade) GetUniqueID() string {
	return u.UpgradeID
}
func (u *Upgrade) GetBasePrice() generaldomain.IBProducts {
	return u.Price
}
func (u *Upgrade) Eligible(env gamedomain.IGameEnvironment) bool {
	return u.UpgradeRequirement.Satisfy(env)
}

func (u *Upgrade) ApplyUpgrade(env gamedomain.IGameEnvironment) bool {
	if u.Eligible(env) {
		for _, v := range u.UpTargets {
			ctype := v.ObjectProperties.PropertiesSlice()[0].GetParamName()
			cval := v.ObjectProperties.PropertiesSlice()[0].GetParamValue()
			switch ctype {
			case "building":
				if v.ObjectProperties.Count() == 1 {
					bid := general.NewObjectProperties()
					bid.AddProperty(ctype, cval)
					fmt.Println(bid)
					for _, eff := range u.Effects {
						//fmt.Println("start")
						env.GetGame().AddBuildingProductModifier(bid, generaldomain.CheckTypeAll, eff.(buildinginstancedomain.IBuildingProductionEffect))
					}
					/*
						for _, sp := range env.Gm.gMap.Spaces {
							el := sp.GetOccupier()
							if el != nil && el.GetParentID() == cval {
								for _, eff := range u.Effects {
									eff.ApplyEffectGlobal(el)
									//somehow add effect to applied Effects
								}
							}
						}
					*/
				}
			case "space":
			}
		}
	}

	return true
}

func (u *Upgrade) RemoveUpgrade() error {
	//fmt.Println("removing")
	if u.Effects == nil {
		return errors.New("Instance does not have space effect")
	}

	for _, v := range u.Effects {
		v.RemoveEffect()
	}

	return nil
}
