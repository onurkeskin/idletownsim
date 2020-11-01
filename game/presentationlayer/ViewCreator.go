package presentationlayer

import (
	buildingdomain "github.com/app/game/applayer/building/domain"
	"github.com/app/game/applayer/buildinginstance"
	buildinginstancedomain "github.com/app/game/applayer/buildinginstance/domain"
	"github.com/app/game/applayer/effect"
	"github.com/app/game/applayer/game"
	"github.com/app/game/applayer/gamemap"
	gamemapdomain "github.com/app/game/applayer/gamemap/domain"
	generaldomain "github.com/app/game/applayer/general/domain"
	buildingview "github.com/app/game/presentationlayer/buildingview"
	buildingviewdomain "github.com/app/game/presentationlayer/buildingview/domain"
	"github.com/app/game/presentationlayer/effectview"
	effectviewdomain "github.com/app/game/presentationlayer/effectview/domain"
	resourceview "github.com/app/game/presentationlayer/resourceview"
	resourceviewdomain "github.com/app/game/presentationlayer/resourceview/domain"
	"strconv"
)

type ViewCreator struct {
	brep buildingviewdomain.IBuildingViewRepository
	rrep resourceviewdomain.IResourceViewRepository
}

func (resource *ViewCreator) ConstructGameEnvView(g game.GameEnvironment) (*GameView, error) {
	var toRet GameView
	toRet.GameID = g.GameID
	toRet.GameCreationTime = g.GameCreationTime
	toRet.LastIterationTime = g.LastIterationTime

	resView, err := resource.ConstructResourcesView(g.Resources)
	if err != nil {
		return nil, err
	}
	toRet.Resources = resView
	sview, err := resource.ConstructSpaceViews(g.GetGame().GetGameMap())
	if err != nil {
		return nil, err
	}
	toRet.Spaces = sview

	return &toRet, nil
}

func (resource *ViewCreator) ConstructGamesEnvView(gs game.GameEnvironments) (GamesView, error) {
	toRet := GamesView{}
	for _, g := range gs {
		toadd, err := resource.ConstructGameEnvView(g)
		if err != nil {
			return nil, err
		}
		toRet = append(toRet, *toadd)
	}
	return toRet, nil
}

func (resource *ViewCreator) ConstructResourcesView(rs generaldomain.IBProducts) (resourceviewdomain.IResourcesView, error) {
	toReturn := resourceview.ResourcesView{}
	for _, r := range rs {
		cons, err := resource.ConstructResourceView(r)
		if err != nil {
			return nil, err
		}
		toReturn = append(toReturn, (*cons.(*resourceview.ResourceView)))
	}
	return toReturn, nil
}

func (resource *ViewCreator) ConstructResourceView(r generaldomain.IBProduct) (resourceviewdomain.IResourceView, error) {
	rrepo, err := resource.rrep.GetResourceViewById(r.GetType())
	toReturn := rrepo.(*resourceview.ResourceView)
	if err != nil {
		return nil, err
	}
	toReturn.ResourceCount = r.GetValue()

	return toReturn, nil
}

func (resource *ViewCreator) ConstructSpaceView(_space gamemapdomain.ISpace) (*SpaceView, error) {
	toReturn := SpaceView{}
	toReturn.SpaceID = _space.GetID()
	toReturn.InMapID = _space.GetInMapID()
	ress, err := resource.ConstructResourcesView(_space.GetExpectedOutcome())
	if err != nil {
		return nil, err
	}

	toReturn.ExpectedResource = ress

	occupier := _space.GetOccupier()
	switch occupier := occupier.(type) {
	case buildinginstancedomain.IBuildingInstance:
		b, err := resource.ConstructBuildingInstanceView(occupier)
		if err != nil {
			return nil, err
		}

		toReturn.Occupier = b
	default:
	}

	//space := _space.(())

	return &toReturn, nil
}

func (resource *ViewCreator) ConstructSpaceViews(spaces gamemapdomain.IGameMap) (SpacesView, error) {
	toReturn := SpacesView{}
	var err error
	spaces.ForSpaces(
		func(s gamemapdomain.ISpace) {
			toadd, er := resource.ConstructSpaceView(s)
			if er != nil {
				err = er
				return
			}
			toReturn = append(toReturn, *toadd)
		})
	if err != nil {
		return nil, err
	}
	return toReturn, nil
}

func (resource *ViewCreator) ConstructBuildingView(buil buildingdomain.IBuilding) (buildingviewdomain.IBuildingView, error) {
	_toReturn, err := resource.brep.GetBuildingViewById(buil.GetID())
	toReturn := _toReturn.(*buildingview.BuildingView)
	if err != nil {
		return nil, err
	}
	rv, err := resource.ConstructResourcesView(buil.GetBaseProducts())
	if err != nil {
		return nil, err
	}
	toReturn.ExpectedResource = rv
	switch e := buil.GetSpaceEffect().(type) {
	case *gamemap.CenteredAreaSpaceEffect:
		ev, err := resource.ConstructCenteredAreaSpaceEffectView(e)
		if err != nil {
			return nil, err
		}
		toReturn.BuildingEffect = ev
	case *gamemap.SpaceEffect:
		ev, err := resource.ConstructSpaceEffectView(e)
		if err != nil {
			return nil, err
		}
		toReturn.BuildingEffect = ev
	}

	return toReturn, nil

}

func (resource *ViewCreator) ConstructBuildingInstanceView(buil buildinginstancedomain.IBuildingInstance) (buildingviewdomain.IBuildingView, error) {
	__toReturn, err := resource.brep.GetBuildingViewById(buil.GetParentID())
	_toReturn := __toReturn.(*buildingview.BuildingView)
	toReturn := &buildingview.BuildingInstanceView{
		ID:                  _toReturn.ID,
		BuildingName:        _toReturn.BuildingName,
		GlobalIdentifier:    _toReturn.GlobalIdentifier,
		BuildingDescription: _toReturn.BuildingDescription,
	}

	if err != nil {
		return nil, err
	}
	rv, err := resource.ConstructResourcesView(buil.GetBaseProducts())
	if err != nil {
		return nil, err
	}
	toReturn.ExpectedResource = rv
	switch e := buil.GetSpaceEffect().(type) {
	case *gamemap.CenteredAreaSpaceEffect:
		ev, err := resource.ConstructCenteredAreaSpaceEffectView(e)
		if err != nil {
			return nil, err
		}
		toReturn.BuildingEffect = ev
	case *gamemap.SpaceEffect:
		ev, err := resource.ConstructSpaceEffectView(e)
		if err != nil {
			return nil, err
		}
		toReturn.BuildingEffect = ev
	}
	toReturn.Level = strconv.Itoa(buil.GetLevel())

	for _, ae := range buil.GetProductMods() {
		aev, err := resource.ConstructBuildingProductionEffectView(ae.(*buildinginstance.BuildingProductionEffect))
		if err != nil {
			return nil, err
		}
		toReturn.AppliedEffects = append(toReturn.AppliedEffects, aev)
	}

	return toReturn, nil
}

func (resource *ViewCreator) ConstructBuildingProductionEffectView(effect *buildinginstance.BuildingProductionEffect) (effectviewdomain.IEffectView, error) {
	_toReturn, err := resource.ConstructGenerationEffectView(&effect.GenerationEffect)
	if err != nil {
		return nil, err
	}
	toReturn := _toReturn.(effectview.EffectView)
	toReturn.EffectType = "buildingproductioneffect"

	return toReturn, nil
}

func (resource *ViewCreator) ConstructSpaceEffectView(effect *gamemap.SpaceEffect) (effectviewdomain.IEffectView, error) {
	_toReturn, err := resource.ConstructGenerationEffectView(&effect.GenerationEffect)
	if err != nil {
		return nil, err
	}
	toReturn := _toReturn.(effectview.EffectView)
	toReturn.EffectType = "spaceeffect"

	return toReturn, nil
}

func (resource *ViewCreator) ConstructCenteredAreaSpaceEffectView(effect *gamemap.CenteredAreaSpaceEffect) (effectviewdomain.IEffectView, error) {
	_toReturn, err := resource.ConstructGenerationEffectView(&effect.GenerationEffect)
	if err != nil {
		return nil, err
	}
	toReturn := _toReturn.(effectview.EffectView)
	toReturn.EffectType = "centeredareaspaceeffect"
	areaArr := make([][]string, len(effect.AreaArr))
	//[len(effect.AreaArr)][len(effect.AreaArr[0])]string{}
	for i := 0; i < len(effect.AreaArr); i++ {
		areaArr[i] = make([]string, len(effect.AreaArr[0]))
	}

	for i := 0; i < len(effect.AreaArr); i++ {
		for j := 0; j < len(effect.AreaArr[0]); j++ {
			areaArr[i][j] = effect.AreaArr[i][j].String()
		}
	}
	toReturn.AreaArr = areaArr

	return toReturn, nil
}

func (resource *ViewCreator) ConstructGenerationEffectView(effect *effect.GenerationEffect) (effectviewdomain.IEffectView, error) {
	toReturn := effectview.EffectView{}

	geneff := effect
	affected := []string{}
	for _, e := range geneff.TargetResources {
		_r, err := resource.rrep.GetResourceViewById(e)
		if err != nil {
			panic(err)
			return nil, err
		}
		r := _r.(*resourceview.ResourceView)
		affected = append(affected, r.ResourceName)
	}
	toReturn.TargetResources = affected

	for i := 0; i < len(geneff.Scheme.Operators); i++ {
		op := geneff.Scheme.Operators[i]
		num := geneff.Scheme.Numbers[i]
		schemeadd := op + strconv.FormatFloat(num, 'E', -1, 64)
		toReturn.Scheme = append(toReturn.Scheme, schemeadd)
	}

	return toReturn, nil
}
