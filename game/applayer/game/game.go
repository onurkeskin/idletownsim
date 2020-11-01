package game

import (
	"fmt"
	"time"

	buildinginstancedomain "github.com/app/game/applayer/buildinginstance/domain"
	gamedomain "github.com/app/game/applayer/game/domain"
	//gamemap "github.com/app/game/applayer/gamemap"
	gamemapdomain "github.com/app/game/applayer/gamemap/domain"
	generaldomain "github.com/app/game/applayer/general/domain"
	generator "github.com/app/game/applayer/generator"
	generatordomain "github.com/app/game/applayer/generator/domain"
)

type Game struct {
	//buildingModifiers map[string][]domain.IEffect
	//spaceModifiers           []gamemapdomain.ISpaceEffect
	//buildingProductModifiers []buildinginstancedomain.IBuildingProductionEffect

	GMap gamemapdomain.IGameMap `json:"gmap" bson:"gmap"`

	gen generatordomain.Generator
}

func NewGame(GMap gamemapdomain.IGameMap) gamedomain.IGame {
	return &Game{
		//spaceModifiers:           []gamemapdomain.ISpaceEffect{},
		//buildingProductModifiers: []buildinginstancedomain.IBuildingProductionEffect{},
		GMap: GMap,
		gen:  generator.NewBasicGenerator(),
	}
}

func (g *Game) PlaceBuildingInstance(sid string, buil buildinginstancedomain.IBuildingInstance) error {
	err := g.GMap.PlaceBuildingInstance(sid, buil)
	if err != nil {
		return err
	}

	return nil
}

func (g *Game) RemoveBuildingInstance(sid string) error {
	return g.GMap.RemoveBuildingInstance(sid)
}

func (g *Game) GetGameMap() gamemapdomain.IGameMap {
	return g.GMap
}

func (g *Game) AddSpace(s gamemapdomain.ISpace) error {
	return g.GMap.AddSpace(s)
}

func (g *Game) RemoveSpace(sid string) error {
	return g.GMap.RemoveSpace(sid)
}

func (g *Game) AddBuildingProductModifier(props generaldomain.IObjectProperties, chkType generaldomain.CheckType, eff buildinginstancedomain.IBuildingProductionEffect) {
	//g.buildingProductModifiers = append(g.buildingProductModifiers, eff)

	g.GMap.ForBuildables(func(b gamemapdomain.ITileBuildableElement) {
		v, ok := b.(buildinginstancedomain.IBuildingInstance)
		if !ok {
			return
		}
		if props.SatisfiesType(b, chkType) {
			eff.ApplyEffect(v)
		}
	})
}

func (g *Game) RemoveBuildingProductModifier(eff buildinginstancedomain.IBuildingProductionEffect) bool {
	/*for in, v := range g.buildingProductModifiers {
		if v == eff {
			g.buildingProductModifiers = append(g.buildingProductModifiers[:in], g.buildingProductModifiers[in+1:]...)
		}
	}*/

	eff.RemoveEffect()

	return false
}

func (g *Game) AddSpaceModifier(props generaldomain.IObjectProperties, chkType generaldomain.CheckType, eff gamemapdomain.ISpaceEffect) {
	//g.spaceModifiers = append(g.spaceModifiers, eff)

	g.GMap.ForSpaces(func(s gamemapdomain.ISpace) {
		if props.SatisfiesType(s, chkType) {
			eff.ApplyEffect(s)
		}
	})
}

func (g *Game) GetBuildablesWithProperty(props generaldomain.IObjectProperties, chkType generaldomain.CheckType) []gamemapdomain.ITileBuildableElement {
	toRet := []gamemapdomain.ITileBuildableElement{}
	g.GMap.ForBuildables(func(b gamemapdomain.ITileBuildableElement) {
		/*
			fmt.Println(b.PropertiesSlice())
			fmt.Println(props)
			fmt.Println("then")
		*/
		if props.SatisfiesType(b, chkType) {
			toRet = append(toRet, b)
		}
	})
	return toRet
}
func (g *Game) GetSpacesWithProperty(props generaldomain.IObjectProperties, chkType generaldomain.CheckType) []gamemapdomain.ISpace {
	toRet := []gamemapdomain.ISpace{}
	g.GMap.ForSpaces(func(s gamemapdomain.ISpace) {
		if props.SatisfiesType(s, chkType) {
			toRet = append(toRet, s)
		}
	})

	return toRet
}

func (g *Game) RemoveSpaceModifier(eff gamemapdomain.ISpaceEffect) bool {
	//if (g.spaceModifiers)
	eff.RemoveEffect()
	return false
}

func (g *Game) RemoveSpaceModifierFrom(eff gamemapdomain.ISpaceEffect, s gamemapdomain.ISpace) bool {
	//if (g.spaceModifiers)
	eff.RemoveEffectFrom(s)
	return false
}

/*
func (g *Game) CountSpaceWith(f func(gamemapdomain.ISpace) bool) int {
	toRet := 0
	for _, v := range g.GMap.Spaces { // TODO CHANGE THIS
		if f(v) {
			toRet++
		}
	}
	return toRet
}
*/
func (g *Game) TickFor(tick time.Duration) generaldomain.IBProducts {
	total := generaldomain.IBProducts{}
	g.GMap.ForSpaces(func(s gamemapdomain.ISpace) {
		curSpace := s
		if curSpace.GetOccupier() != nil {
			v, ok := curSpace.GetOccupier().(buildinginstancedomain.IBuildingInstance)
			if ok {
				builTotal := generaldomain.IBProducts{}
				returnedresources := v.DoWork(tick)
				for _, net := range returnedresources {
					builTotal.Add(curSpace.DoWork(net))
				}
				//fmt.Println(builTotal)
				total.Add(builTotal)
			}
		}
	})
	return total
}

func (g *Game) String() string {
	str := ""
	//str += fmt.Sprintf("\n---------------------Current Space Modifiers---------------------\n%v\n", g.spaceModifiers)
	//str += fmt.Sprintf("---------------------------------------------------------------\n")
	//str += fmt.Sprintf("---------------------Current Building Modifiers---------------------\n%v\n", g.buildingProductModifiers)
	//str += fmt.Sprintf("---------------------------------------------------------------\n")
	str += fmt.Sprintf("---------------------Current Map---------------------\n%s\n", g.GMap)
	str += fmt.Sprintf("---------------------------------------------------------------\n")
	str += fmt.Sprintf("---------------------Current Generator---------------------\n%s\n", g.gen)
	str += fmt.Sprintf("---------------------------------------------------------------\n")
	return str
}

func StringContains(a []string, word string) bool {
	for _, el := range a {
		//fmt.Println(el, " ", word)
		if el == word {
			return true
		}
	}
	return false
}
