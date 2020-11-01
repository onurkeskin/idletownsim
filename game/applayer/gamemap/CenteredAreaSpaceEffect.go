package gamemap

import (
	"errors"
	"fmt"

	"github.com/app/game/applayer/datastructures"
	effectdomain "github.com/app/game/applayer/effect/domain"
	gamemapdomain "github.com/app/game/applayer/gamemap/domain"
	general "github.com/app/game/applayer/general"
	//generaldomain "github.com/app/game/applayer/general/domain"
)

type CenteredSEElementType int

const ( // iota is reset to 0
	Unknown            CenteredSEElementType = 0
	SpaceCenterElement                       = 1
	SpaceEmpty                               = 2
	SpaceElement                             = 3
)

func (c CenteredSEElementType) String() string {
	switch c {
	case SpaceCenterElement:
		return "Center"
	case SpaceEmpty:
		return "Empty"
	case SpaceElement:
		return "Element"
	default:
		return "Unknown"
	}
}

type CenteredAreaSpaceEffect struct {
	SpaceEffect `json:"spaceeffect" bson:"spaceeffect"`
	AreaArr     [][]CenteredSEElementType `json:"spaceeffectarea" bson:"spaceeffectarea"`
}

func NewCenteredAreaSpaceEffect(
	ID string,
	Priority int64,
	AreaArr [][]CenteredSEElementType,
	targetResources []string,
	scheme *general.MathValScheme,
	issuer effectdomain.IEffectIssuer) gamemapdomain.ISpaceEffect {
	if verifyAreaArr(AreaArr) {
		return &CenteredAreaSpaceEffect{}
	}

	return &CenteredAreaSpaceEffect{
		SpaceEffect: *NewSpaceEffect(
			ID,
			Priority,
			targetResources,
			scheme,
			issuer).(*SpaceEffect),
		AreaArr: AreaArr,
	}
}

func (eff *CenteredAreaSpaceEffect) Clone() interface{} {
	return &CenteredAreaSpaceEffect{
		SpaceEffect: *eff.SpaceEffect.Clone().(*SpaceEffect),
		AreaArr:     eff.AreaArr}
}

func (eff *CenteredAreaSpaceEffect) ApplyEffectGlobal(anything interface{}) {
	v, ok := anything.(gamemapdomain.ISpace)
	if !ok {
		return
	}

	eff.ApplyEffect(v)
}

func (eff *CenteredAreaSpaceEffect) ApplyEffect(s gamemapdomain.ISpace) {
	centerX, centerY, err := findCenterElementPos(eff.AreaArr)
	//fmt.Println("willerr?")
	if err != nil {
		return
	}
	//fmt.Println("noerr?")
	schemeMap := findAndMarkReachables(s, centerX, centerY)
	//fmt.Println(schemeMap)
	for i := 0; i < len(eff.AreaArr); i++ {
		for j := 0; j < len(eff.AreaArr); j++ {
			// 1)FIND ALL REACHABLES
			// 2) MARK THEM WITH THEIR POSITIONS RELATIVE
			// 3) DETERMINE ALL
			if eff.AreaArr[i][j] != SpaceEmpty { // MAYBE ROTATING THINGS AROUND IS BETTER LOL)
				for key, val := range schemeMap {
					if val[0] == i && val[1] == j {
						(key).AddSpaceMod(eff)
						if eff.AreaArr[i][j] == SpaceCenterElement {
							eff.AppliedTo = append(eff.AppliedTo, &key)
						}
					}
				}
			}
		}
	}

	//fmt.Println("")
	//fmt.Println(schemeMap)
	//fmt.Println("----------end--------")
	return
}

func (eff *CenteredAreaSpaceEffect) RemoveEffect() {
	for _, s := range eff.AppliedTo {
		centerX, centerY, err := findCenterElementPos(eff.AreaArr)
		//fmt.Println("willerr?")
		if err != nil {
			return
		}
		//fmt.Println("noerr?")
		schemeMap := findAndMarkReachables(*s, centerX, centerY)
		//fmt.Println(schemeMap)
		for i := 0; i < len(eff.AreaArr); i++ {
			for j := 0; j < len(eff.AreaArr); j++ {
				// 1)FIND ALL REACHABLES
				// 2) MARK THEM WITH THEIR POSITIONS RELATIVE
				// 3) DETERMINE ALL
				if eff.AreaArr[i][j] != SpaceEmpty { // MAYBE ROTATING THINGS AROUND IS BETTER LOL)
					for key, val := range schemeMap {
						if val[0] == i && val[1] == j {
							//fmt.Println(fmt.Sprintf("%d, %d", i, j))
							(key).RemoveSpaceMod(eff)
						}
					}
				}
			}
		}
	}
}

func (eff *CenteredAreaSpaceEffect) ReapplyEffect() {
	toApply := eff.AppliedTo
	eff.RemoveEffect()

	for _, s := range toApply {
		eff.ApplyEffect(*s)
	}
}

func (eff *CenteredAreaSpaceEffect) RemoveEffectFrom(space gamemapdomain.ISpace) {
	for _, s := range eff.AppliedTo {
		if *s == space {

			centerX, centerY, err := findCenterElementPos(eff.AreaArr)
			//fmt.Println("willerr?")
			if err != nil {
				return
			}
			//fmt.Println("noerr?")
			schemeMap := findAndMarkReachables(*s, centerX, centerY)
			//fmt.Println(schemeMap)
			for i := 0; i < len(eff.AreaArr); i++ {
				for j := 0; j < len(eff.AreaArr); j++ {
					// 1)FIND ALL REACHABLES
					// 2) MARK THEM WITH THEIR POSITIONS RELATIVE
					// 3) DETERMINE ALL
					if eff.AreaArr[i][j] != SpaceEmpty { // MAYBE ROTATING THINGS AROUND IS BETTER LOL)
						for key, val := range schemeMap {
							if val[0] == i && val[1] == j {
								//fmt.Println(fmt.Sprintf("%d, %d", i, j))
								(key).RemoveSpaceMod(eff)
							}
						}
					}
				}
			}
		}
	}
}

func (eff *CenteredAreaSpaceEffect) String() string {
	str := ""
	str += fmt.Sprintf("GenerationEffect:%s ", eff.GenerationEffect.String())
	str += fmt.Sprintf("Aoe:%v ", eff.AreaArr)
	return str
}

func findAndMarkReachables(toProcess gamemapdomain.ISpace, i int, j int) map[gamemapdomain.ISpace][2]int {
	processed := make(map[gamemapdomain.ISpace][2]int)
	q := datastructures.NewQueue(1)
	q.Push(&datastructures.Node{toProcess})
	processed[toProcess] = [2]int{i, j}
	for q.Count() > 0 {
		cur := q.Pop().Value.(gamemapdomain.ISpace)
		curX := processed[cur][0]
		curY := processed[cur][1]
		for i := Northeast; i <= East; i++ {
			curNs, _ := (cur).GetNeighboorsAt(i)
			for _, curN := range curNs {
				if curN != nil {
					switch i {
					case Northeast:
						_, ok := processed[curN]
						if !ok {
							processed[curN] = [2]int{curX - 1, curY + 1}
							q.Push(&datastructures.Node{curN})
						}
					case North:
						_, ok := processed[curN]
						if !ok {
							processed[curN] = [2]int{curX - 1, curY}
							q.Push(&datastructures.Node{curN})
						}
					case Northwest:
						_, ok := processed[curN]
						if !ok {
							processed[curN] = [2]int{curX - 1, curY - 1}
							q.Push(&datastructures.Node{curN})
						}
					case West:
						_, ok := processed[curN]
						if !ok {
							processed[curN] = [2]int{curX, curY - 1}
							q.Push(&datastructures.Node{curN})
						}
					case Southwest:
						_, ok := processed[curN]
						if !ok {
							processed[curN] = [2]int{curX + 1, curY - 1}
							q.Push(&datastructures.Node{curN})
						}
					case South:
						_, ok := processed[curN]
						if !ok {
							processed[curN] = [2]int{curX + 1, curY}
							q.Push(&datastructures.Node{curN})
						}
					case Southeast:
						_, ok := processed[curN]
						if !ok {
							processed[curN] = [2]int{curX + 1, curY + 1}
							q.Push(&datastructures.Node{curN})
						}
					case East:
						_, ok := processed[curN]
						if !ok {
							processed[curN] = [2]int{curX, curY + 1}
							q.Push(&datastructures.Node{curN})
						}
					}
				}
			}
		}
	}

	return processed
}

func verifyAreaArr(arr [][]CenteredSEElementType) bool {
	centerElementCount := 0

	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr[0]); j++ {
			switch arr[i][j] {
			case SpaceCenterElement:
				centerElementCount++
				if centerElementCount > 1 {
					return true
				}
			case SpaceElement:

			case SpaceEmpty:

			default:

			}
		}
	}
	return false
}

func findCenterElementPos(arr [][]CenteredSEElementType) (int, int, error) {
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr[0]); j++ {
			if arr[i][j] == SpaceCenterElement {
				return i, j, nil
			}
		}
	}
	return 0, 0, errors.New("No Center Element")
}

// MAYBE USE IN FUTURE
type Direction struct {
	dirs [][]int
}

func (d *Direction) getDirections() [][]int {
	return d.dirs
}
