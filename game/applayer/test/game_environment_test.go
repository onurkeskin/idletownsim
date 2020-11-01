package test

import (
	building "github.com/app/game/applayer/building"
	buildinginstance "github.com/app/game/applayer/buildinginstance"
	"gopkg.in/mgo.v2/bson"
	//effect "github.com/game/applayer/effect"
	effectdomain "github.com/app/game/applayer/effect/domain"
	game "github.com/app/game/applayer/game"
	gamemap "github.com/app/game/applayer/gamemap"
	gamemapdomain "github.com/app/game/applayer/gamemap/domain"
	general "github.com/app/game/applayer/general"
	up "github.com/app/game/applayer/upgrade"

	"fmt"
	generaldomain "github.com/app/game/applayer/general/domain"

	"testing"
	"time"
)

func TestGameEnv(t *testing.T) {
	building1ID := bson.NewObjectId()
	building2ID := bson.NewObjectId()

	Space1ID := "1"
	Space2ID := "2"
	Space3ID := "3"
	Space4ID := "4"
	baseProduct1Type := "money"
	baseProduct2Type := "happiness"
	bProduct1Val := float64(5)
	bProduct2Val := float64(5)
	testBProduction := generaldomain.IBProducts{general.NewRegularBProduct(baseProduct1Type, bProduct1Val), general.NewRegularBProduct(baseProduct2Type, bProduct2Val)}
	productionTimeSec := int64(3)
	prodTimeSec := float64(1)
	timesWorkedExpected := float64(productionTimeSec) / prodTimeSec

	doubler, _ := general.NewMathValScheme([]string{"*"}, []float64{2})
	add2, _ := general.NewMathValScheme([]string{"+"}, []float64{2})
	eff3 := gamemap.NewCenteredAreaSpaceEffect(
		"TileAreaEffect",
		150,
		[][]gamemap.CenteredSEElementType{
			{gamemap.SpaceEmpty, gamemap.SpaceElement, gamemap.SpaceEmpty},           /*  initializers for row indexed by 0 */
			{gamemap.SpaceElement, gamemap.SpaceCenterElement, gamemap.SpaceElement}, /*  initializers for row indexed by 1 */
			{gamemap.SpaceEmpty, gamemap.SpaceElement, gamemap.SpaceEmpty},           /*  initializers for row indexed by 2 */
		}, []string{baseProduct1Type}, doubler, nil,
	)

	testBuilding := building.NewBuilding(
		building1ID,
		nil,
		nil,
		nil,
		testBProduction,
		int64(prodTimeSec*1000)*time.Millisecond.Nanoseconds(),
		eff3)

	testBuilding2 := building.NewBuilding(
		building2ID,
		nil,
		nil,
		nil,
		testBProduction,
		int64(prodTimeSec*1000)*time.Millisecond.Nanoseconds(),
		eff3)

	testBInstance1 := buildinginstance.FormFromBuilding(testBuilding, "builId1", 1)
	testBInstance1.SetBuiltTimeUnix(time.Now())
	testBInstance1.AddProperty("me", "2")
	testBInstance2 := buildinginstance.FormFromBuilding(testBuilding2, "builId2", 1)
	testBInstance2.SetBuiltTimeUnix(time.Now())

	testBInstance3 := buildinginstance.FormFromBuilding(testBuilding2, "builId3", 1)
	testBInstance3.SetBuiltTimeUnix(time.Now())

	testBInstance4 := buildinginstance.FormFromBuilding(testBuilding2, "builId4", 1)
	testBInstance4.SetBuiltTimeUnix(time.Now())

	s1 := gamemap.NewSpace(
		Space1ID,
		nil,
	)
	s2 := gamemap.NewSpace(
		Space2ID,
		nil,
	)
	s3 := gamemap.NewSpace(
		Space3ID,
		nil,
	)
	s4 := gamemap.NewSpace(
		Space4ID,
		nil,
	)
	s1.AddNeighboorTo(s2, gamemap.East)
	s1.AddNeighboorTo(s3, gamemap.North)
	s1.AddNeighboorTo(s4, gamemap.South)
	s1.AddProperty("me", "2")

	s2.AddNeighboorTo(s1, gamemap.West)
	s3.AddNeighboorTo(s1, gamemap.South)
	s4.AddNeighboorTo(s1, gamemap.North)

	gMap := gamemap.GameMap{[]gamemapdomain.ISpace{}, time.Now(), "v1"}

	gm := game.NewGame(gMap)
	gm.AddSpace(s1)
	gm.AddSpace(s2)
	gm.AddSpace(s3)
	gm.AddSpace(s4)
	gm.PlaceBuildingInstance(Space1ID, testBInstance1)
	gm.PlaceBuildingInstance(Space2ID, testBInstance2)
	gm.PlaceBuildingInstance(Space3ID, testBInstance3)
	gm.PlaceBuildingInstance(Space4ID, testBInstance4)
	//	fmt.Println(game)

	eff := buildinginstance.NewBuildingProductionEffect("qwezcs", 1, []string{baseProduct1Type}, doubler, nil)
	eff2 := gamemap.NewSpaceEffect("qwzxcasd", 200, []string{baseProduct1Type}, add2, nil)

	bchkProps := general.NewObjectProperties()
	bchkProps.AddProperty("me", "2")
	gm.AddBuildingProductModifier(bchkProps, generaldomain.CheckTypeAll, eff)

	//game.RemoveBuildingProductModifier("1", eff)
	chkProps := general.NewObjectProperties()
	chkProps.AddProperty("me", "2")
	gm.AddSpaceModifier(chkProps, generaldomain.CheckTypeAll, eff2)

	//game.RemoveBuildingInstance(Space1ID)
	//game.RemoveSpaceModifier(eff2)
	//testBInstance1.RemoveSpaceEffect()
	//fmt.Println(s1)
	//time.Sleep(time.Duration(productionTimeSec) * time.Second)
	//t.Log(med)

	eff4 := buildinginstance.NewBuildingProductionEffect("eff4", 1, []string{baseProduct1Type}, doubler, nil)
	_ = eff4
	compProperty := general.NewObjectProperties()
	compProperty.AddProperty("building", building1ID.Hex())
	comp := up.ComplishmentDB{
		ID:     "asd",
		Target: []up.ComplishmentTarget{up.ComplishmentTarget{compProperty}},
		Params: []up.ComplishmentParams{
			up.ComplishmentParams{
				InequalitySymbol: "=",
				Value:            1,
			}},
	}
	fn, _ := comp.ComplishmentParse()
	req := up.NewRequirement(
		"req1",
		[]up.ComplishmentTestFunc{fn})
	//fmt.Println(req)

	targetProperty := general.NewObjectProperties()
	targetProperty.AddProperty("building", building1ID.Hex())
	up := up.NewUpgrade("upgradedoubleb1",
		nil,
		req,
		[]up.UpgradeTarget{
			up.UpgradeTarget{
				ObjectProperties: targetProperty}},
		[]effectdomain.IEffect{eff4})
	//_ = up
	gameEnv := game.NewGameEnvironment("gameenv1", "asdasd", "builton1", generaldomain.IBProducts{}, []string{}, gm)
	gameEnv.AddUpgrade(up)
	//gameEnv.RemoveUpgrade(*up)
	ExpectedRes1 := float64((((bProduct1Val)*2*2+float64(2))*2*2*2*2)*float64(timesWorkedExpected) + (bProduct1Val*3*4)*float64(timesWorkedExpected))
	ExpectedRes2 := float64(bProduct2Val * float64(timesWorkedExpected) * 4)

	t1 := time.Now().Nanosecond()
	//game.TickUntil(int64(time.Now().UnixNano()))
	gameEnv.PlayFor(1 * 3 * time.Second.Nanoseconds())
	med := gameEnv.GetResources()
	t2 := time.Now().Nanosecond()
	fmt.Println(fmt.Sprintf("ms:%d", (t2-t1)/int(time.Millisecond)))

	//marshalled, err := json.Marshal(med)
	//fmt.Println(err)
	//fmt.Println(marshalled)
	//fmt.Println(gm)
	if len(med) == 0 {
		t.Error("Wrong production")
	}
	for _, check := range med {
		if check.GetType() == baseProduct1Type {
			if check.GetValue() != ExpectedRes1 {
				t.Error("Current type1 :", baseProduct1Type, " val:", check.GetValue(), " Expected: ", ExpectedRes1)
				t.Error("Wrong production")
			}
		} else if check.GetType() == baseProduct2Type {
			if check.GetValue() != ExpectedRes2 {
				t.Error("Current type1 :", baseProduct2Type, " val:", check.GetValue(), " Expected: ", ExpectedRes2)
				t.Error("Wrong production")
			}
		} else {
			t.Error("Wrong Production type")
			t.Error("Wrong production")
		}
	}
	//t.Log(game)
	t.Log("Game Test Finished")
}
