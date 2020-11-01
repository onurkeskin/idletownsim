package gamemap

import (
	"errors"
	"fmt"

	buildinginstancedomain "github.com/app/game/applayer/buildinginstance/domain"
	gamemapdomain "github.com/app/game/applayer/gamemap/domain"

	general "github.com/app/game/applayer/general"
	generaldomain "github.com/app/game/applayer/general/domain"

	generator "github.com/app/game/applayer/generator"
	generatordomain "github.com/app/game/applayer/generator/domain"
)

const ( // iota is reset to 0
	Northeast = iota
	North     = iota
	Northwest = iota
	West      = iota
	Southwest = iota
	South     = iota
	Southeast = iota
	East      = iota
)

type Space struct {
	*general.ObjectProperties

	SpaceId string `json:"id,omitempty" bson:"_id,omitempty"`
	InMapID string `json:"inmapid" bson:"inmapid"`

	Element gamemapdomain.ITileBuildableElement `json:"element" bson:"element"`

	gen generatordomain.Generator `json:"-" bson:"-"`

	around [8][]gamemapdomain.ISpace `json:"-" bson:"-"`

	elChangedListeners []gamemapdomain.ElementChangedListener `json:"-" bson:"-"`
}

func NewSpace(
	SpaceId string,
	InMapID string,
	Element gamemapdomain.ITileBuildableElement) gamemapdomain.ISpace {

	var toRet *Space
	toRet = &Space{
		ObjectProperties:   general.NewObjectProperties(),
		SpaceId:            SpaceId,
		InMapID:            InMapID,
		gen:                generator.NewBasicGenerator(),
		around:             [8][]gamemapdomain.ISpace{},
		elChangedListeners: []gamemapdomain.ElementChangedListener{}}

	toRet.AddElement(Element)
	return toRet
}

func (s *Space) GetOccupier() gamemapdomain.ITileBuildableElement {
	return s.Element
}

func (s *Space) SetElement(Element gamemapdomain.ITileBuildableElement) error {
	s.Element = Element
	if Element.GetSpaceEffect() == nil {
		return nil
	}
	s.SendElementChangedEvent(Element)
	return nil
}

func (s *Space) AddElement(Element gamemapdomain.ITileBuildableElement) error {
	if Element == nil {
		//fmt.Println("Element nill")
		return errors.New(fmt.Sprintf("Cant add nil Element"))
	}
	if s.Element != nil {
		return errors.New(fmt.Sprintf("Already have a building on space%s", s.GetID()))
	}

	s.Element = Element
	if Element.GetSpaceEffect() == nil {
		return nil
	}
	Element.ApplySpaceEffect(s)

	s.SendElementChangedEvent(Element)
	return nil
}

func (s *Space) RemoveElement() error {
	if s.Element == nil {
		return errors.New(fmt.Sprintf("Already Empty Space%s", s.GetID()))
	}
	s.Element.RemoveSpaceEffect()
	s.Element = nil
	return nil
}

func (s *Space) AddNeighboorTo(ns gamemapdomain.ISpace, dir int) error {
	if s == ns {
		return errors.New(fmt.Sprintf("Cant have self neighboor relation"))
	}
	if dir > East {
		return errors.New(fmt.Sprintf("Direction out of bounds.(Must Be Between %d and %d)", Northeast, East))
	}
	if s.around[dir] != nil {
		return errors.New(fmt.Sprintf("Direction Not initialized:%d", dir))
	}
	if s.neighboorsDoesContain(ns) {
		return errors.New(fmt.Sprintf("Space cant be added twice"))
	}

	s.around[dir] = append(s.around[dir], ns)
	return nil
}

func (s *Space) GetNeighboorsAt(dir int) ([]gamemapdomain.ISpace, error) {
	if dir > East {
		return nil, errors.New(fmt.Sprintf("Direction out of bounds.(Must Be Between %d and %d)", Northeast, East))
	}
	if s.around[dir] == nil {
		return nil, errors.New(fmt.Sprintf("No Element at pos:%d", dir))
	}
	return s.around[dir], nil
}

func (s *Space) AddSpaceMod(eff gamemapdomain.ISpaceEffect) {
	s.gen.AddProductionMod(eff)
}
func (s *Space) RemoveSpaceMod(eff gamemapdomain.ISpaceEffect) bool {
	s.gen.RemoveProductionMod(eff)
	return true
}

func (s *Space) ResetSpaceMods() {
	s.gen.Reset()
}

func (s *Space) SendElementChangedEvent(Element gamemapdomain.ITileBuildableElement) {
	for _, v := range s.elChangedListeners {
		v.OnSpaceElementChange(Element)
	}
}

func (s *Space) AddElementChangeListener(listener gamemapdomain.ElementChangedListener) error {
	for _, v := range s.elChangedListeners {
		if v == listener {
			return errors.New(fmt.Sprintf("Already has listener:%v", listener))
		}
	}
	s.elChangedListeners = append(s.elChangedListeners, listener)
	return nil
}

func (s *Space) RemoveElementChangeListener(listener gamemapdomain.ElementChangedListener) error {
	for in, v := range s.elChangedListeners {
		if v == listener {
			s.elChangedListeners = append(s.elChangedListeners[:in], s.elChangedListeners[in+1:]...)
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Listener not found:%v", listener))
}

func (s *Space) DoWork(total generaldomain.IBProducts) generaldomain.IBProducts {
	/*	total := domain.IBProducts{}
		if s.Element != nil {
			v, ok := s.Element.(domain.IBuildingInstance)
			if ok {
				builTotal := domain.IBProducts{}
				i := int64(0)
				lastProd := v.GetLastProductionTimeUnix()
				interval := v.GetBaseProductionInterval()
				ProductionToDo := (curTime - lastProd) / interval
				Productionleft := (curTime - lastProd) % interval
				for i = 0; i < ProductionToDo; i++ {
					builTotal.Add(v.DoWork())
				}
				v.SetLastProductionTimeUnix(curTime - Productionleft)

				total.Add(s.gen.Generate(builTotal))
			}
		}*/ // MAYBE USEFUL IF IDEA CHANGES
	toRet := s.gen.Generate(total)
	return toRet
}

func (s *Space) GetExpectedOutcome() generaldomain.IBProducts {
	total := generaldomain.IBProducts{}
	if s.Element != nil {
		v, ok := s.Element.(buildinginstancedomain.IBuildingInstance)
		if ok {
			res := v.GetExpectedOutcome()
			total.Add(res)
		}
	}
	total.Add(s.gen.Generate(total))

	return total
}

func (g *Space) GetID() string {
	return g.SpaceId
}

func (g *Space) GetInMapID() string {
	return g.InMapID
}

func (s *Space) neighboorsDoesContain(ns gamemapdomain.ISpace) bool {
	for _, v := range s.around {
		for _, s := range v {
			if s.GetID() == ns.GetID() {
				return true
			}
		}
	}
	return false
}

func (s *Space) String() string {
	str := ""
	str += fmt.Sprintf("Space Id:%s\n", s.SpaceId)
	str += fmt.Sprintf("Properties:%s\n", s.ObjectProperties)
	str += fmt.Sprintf("------------------Occupier------------------\n%s\n", s.Element)
	str += fmt.Sprintf("------------------Generator------------------\n%s\n", s.gen)
	str += fmt.Sprintf("------------------Neighboors------------------\n")
	totalNeighboors := 0
	for i := Northeast; i <= East; i++ {
		totalNeighboors = totalNeighboors + len(s.around[i])
	}
	str += fmt.Sprintf("\nTotal Neighboors:%d", totalNeighboors)

	for i := Northeast; i <= East; i++ {
		curN := s.around[i]
		if curN != nil {
			switch i {
			case Northeast:
				str += fmt.Sprintf("\nNortheast:")
				for _, v := range curN {
					str += fmt.Sprintf("%p", v)
				}
			case North:
				str += fmt.Sprintf("\nNorth:")
				for _, v := range curN {
					str += fmt.Sprintf("%p", v)
				}
			case Northwest:
				str += fmt.Sprintf("\nNorthwest:")
				for _, v := range curN {
					str += fmt.Sprintf("%p", v)
				}
			case West:
				str += fmt.Sprintf("\nWest:")
				for _, v := range curN {
					str += fmt.Sprintf("%p", v)
				}
			case Southwest:
				str += fmt.Sprintf("\nSouthwest:")
				for _, v := range curN {
					str += fmt.Sprintf("%p", v)
				}
			case South:
				str += fmt.Sprintf("\nSouth:")
				for _, v := range curN {
					str += fmt.Sprintf("%p", v)
				}
			case Southeast:
				str += fmt.Sprintf("\nSoutheast:")
				for _, v := range curN {
					str += fmt.Sprintf("%p", v)
				}
			case East:
				str += fmt.Sprintf("\nEast:")
				for _, v := range curN {
					str += fmt.Sprintf("%p", v)
				}
			}
		}
	}
	str += fmt.Sprintf("\n")
	return str
}
