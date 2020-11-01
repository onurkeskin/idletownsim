package gamemap

import (
	"errors"
	"fmt"
	"github.com/app/helpers/version"
	"time"

	//buildinginstancedomain "github.com/app/game/applayer/buildinginstance/domain"
	gamemapdomain "github.com/app/game/applayer/gamemap/domain"
)

type GameMap struct {
	Spaces []gamemapdomain.ISpace `json:"spaces" bson:"spaces"`
	Ver    version.Version        `json:"version" bson:"version"`
}

func NewGameMap(s []gamemapdomain.ISpace, Version version.Version) gamemapdomain.IGameMap {
	return &GameMap{
		Spaces: s,
		Ver:    Version,
	}
}

func (g *GameMap) GetMapDate() time.Time {
	return g.Ver.VerDateAfter
}
func (g *GameMap) SetMapDate(nd time.Time) {
	g.Ver.VerDateAfter = nd
}

func (g *GameMap) GetMapVersion() version.Version {
	return g.Ver
}
func (g *GameMap) SetMapVersion(Version version.Version) {
	g.Ver = Version
}

func (g *GameMap) HasSpace(id string) bool {
	for _, v := range g.Spaces {
		//fmt.Printf("%v, %v\n", v.GetID(), id)
		if v.GetID() == id {
			return true
		}
	}
	return false
}

func (g *GameMap) AddSpace(s gamemapdomain.ISpace) error {
	if g.Spaces == nil {
		return errors.New("Map Spaces not initialized yet")
	}

	for _, space := range g.Spaces {
		if space.GetID() == s.GetID() {
			return errors.New(fmt.Sprintf("Already have a space with id=%s", space.GetID()))
		}
	}

	g.Spaces = append(g.Spaces, s)
	return nil
}

func (g *GameMap) RemoveSpace(sid string) error {
	if g.Spaces == nil {
		return errors.New("Map Spaces not initialized yet")
	}

	for in, space := range g.Spaces {
		if space.GetID() == sid {
			g.Spaces[in] = nil
		}
	}
	return errors.New(fmt.Sprintf("Unknown space for id=%s", sid))
}

func (g *GameMap) ForBuildables(f func(b gamemapdomain.ITileBuildableElement)) {
	for _, v := range g.Spaces {
		if v.GetOccupier() != nil {
			f(v.GetOccupier())
		}
	}
}

func (g *GameMap) ForSpaces(f func(s gamemapdomain.ISpace)) {
	for _, v := range g.Spaces {
		f(v)
	}
}
func (g *GameMap) SetBuildingInstance(sid string, buil gamemapdomain.ITileBuildableElement) error {
	if g.Spaces == nil {
		return errors.New("Map Spaces not initialized yet")
	}

	for _, space := range g.Spaces {
		if space.GetID() == sid {
			return space.SetElement(buil)
		}
	}
	return errors.New(fmt.Sprintf("Unknown space for id=%s", sid))
}

func (g *GameMap) PlaceBuildingInstance(sid string, buil gamemapdomain.ITileBuildableElement) error {
	if g.Spaces == nil {
		return errors.New("Map Spaces not initialized yet")
	}

	for _, space := range g.Spaces {
		if space.GetID() == sid {
			return space.AddElement(buil)
		}
	}
	return errors.New(fmt.Sprintf("Unknown space for id=%s", sid))
}
func (g *GameMap) RemoveBuildingInstance(sid string) error {
	for _, val := range g.Spaces {
		if val.GetID() == sid {
			return val.RemoveElement()
		}
	}

	return errors.New(fmt.Sprintf("Space Not Found=%s", sid))
}

func (g *GameMap) String() string {
	str := ""
	str += fmt.Sprintf("Version:%v \n", g.Ver)
	str += fmt.Sprintf("------------------Spaces------------------\n")
	str += fmt.Sprintf("Total Spaces:%d\n", len(g.Spaces))
	for _, v := range g.Spaces {
		str += fmt.Sprintf("Space:%v \n", v)
	}

	return str
}
