// The MIT License (MIT)

// Copyright (c) 2015 Hafiz Ismail

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.


package server_test

import (
	"encoding/json"
	"github.com/app/game/applayer/buildinginstance"
	"github.com/app/helpers/version"
	. "github.com/onsi/ginkgo"
	"net/http"
	"os"

	//. "github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"testing"

	"errors"
	"fmt"
	"github.com/app/game/applayer/gamemap"
	"github.com/app/game/applayer/general"
	up "github.com/app/game/applayer/upgrade"
	buildingdb "github.com/app/game/serverlayer/building"
	"github.com/app/game/serverlayer/dbmodels"
	resourcedb "github.com/app/game/serverlayer/gameresources"
	spacedb "github.com/app/game/serverlayer/space"
	upgradedb "github.com/app/game/serverlayer/upgrade"
	usergame "github.com/app/game/serverlayer/usergame"
	mapdb "github.com/app/map"
	mapmodels "github.com/app/map/mapmodels"
	"github.com/app/server/middlewares/mongodb"
	"github.com/app/server/middlewares/renderer"
	"github.com/app/server/server"
	"github.com/app/sessions"
	"github.com/app/users"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"time"
)

var TestDatabaseServerName = "localhost"
var TestDatabaseName = "test_db"
var (
	s *server.Server
)
var Configuration TestConfiguration

type TestConfiguration struct {
	Setup_db    bool `json:"setup_db"`
	Create_user bool `json:"create_user"`
	Create_game bool `json:"create_game"`
}

func TestServer(t *testing.T) {
	//	defineFactories()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}

var _ = BeforeSuite(func() {
	testconffile, e := os.Open("./testresources/testconf.json")
	if e != nil {
		fmt.Printf("Test Config couldnt read error: %v\n", e)
	}
	decoder := json.NewDecoder(testconffile)
	panicIfFails(decoder.Decode(&Configuration))

	privateSigningKey, err := ioutil.ReadFile("./keys/demokeys.rsa")
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error loading private signing key: %v", err.Error())))
	}
	publicSigningKey, err := ioutil.ReadFile("./keys/demokeys.rsa.pub")
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error loading public signing key: %v", err.Error())))
	}

	db := mongodb.New(&mongodb.Options{
		ServerName:   "localhost",
		DatabaseName: "test-go-app",
	})
	_ = db.NewSession()

	// init renderer
	renderer := renderer.New(&renderer.Options{
		IndentJSON: true,
	}, renderer.JSON)

	// set up users resource
	usersResource := users.NewResource(&users.Options{
		Database: db,
		Renderer: renderer,
	})

	// set up sessions resource
	sessionsResource := sessions.NewResource(&sessions.Options{
		Database:              db,
		Renderer:              renderer,
		PrivateSigningKey:     privateSigningKey,
		PublicSigningKey:      publicSigningKey,
		UserRepositoryFactory: usersResource.UserRepositoryFactory,
	})

	BuildingFactory := buildingdb.NewBuildingRepositoryFactory()
	ResourceFactory := resourcedb.NewResourceRepositoryFactory()
	UpgradeFactory := upgradedb.NewUpgradeRepositoryFactory()
	SpaceFactory := spacedb.NewSpaceRepositoryFactory()

	mapResource := mapdb.NewResource(&mapdb.Options{
		Database: db,
		Renderer: renderer,
		MapProvider: &mapdb.MapProvider{
			Retriever: func(mapmodels.LatLng) (v mapdb.MapProviderReturn, err error) {
				dat, err := ioutil.ReadFile("./testresources/lastImage.png")
				if err != nil {
					panic(err)
				}
				return mapdb.MapProviderReturn{
					Image:               dat,
					MapIdentifierAdress: "test",
					MapCompleteAddress:  "test",
				}, err
			}},
	})

	usergameResource := usergame.NewResource(&usergame.Options{
		Database:                  db,
		Renderer:                  renderer,
		UpgradeRepositoryFactory:  UpgradeFactory,
		BuildingRepositoryFactory: BuildingFactory,
		MapRepositoryFactory:      mapResource.MapRepositoryFactory,
		SpaceRepositoryFactory:    SpaceFactory,
		MapProvider:               mapResource.GetMapProvider(),
		ResourceRepositoryFactory: ResourceFactory,
	})

	if Configuration.Setup_db {
		fmt.Println("Database Initiated")
		db.DropDatabase()
		InsertTestValuesIntoResource(usergameResource)
	}

	// init server
	s = server.NewServer(&server.Config{})

	// set up router
	ac := server.NewAccessController(renderer)
	router := server.NewRouter(ac)

	// add REST resources to router
	router.AddResources(sessionsResource, usersResource, usergameResource, mapResource)

	// add middlewares
	//s.UseContextMiddleware(renderer)
	s.UseMiddleware(sessionsResource.NewAuthenticator())

	s.UseRouter(router)

	CertPath := "keys/server.pem"
	KeyPath := "keys/server.key"
	go s.Run(":8008", server.Options{
		Timeout:  1 * time.Millisecond,
		CertPath: CertPath,
		KeyPath:  KeyPath,
	})
	time.Sleep(100 * time.Millisecond)
	fmt.Println("setup complete")
})

var _ = AfterSuite(func() {
	s.Stop()
})

func InsertTestValuesIntoResource(res *usergame.Resource) {
	builrepo := res.BuildingRepository(nil)
	resrepo := res.ResourceRepository(nil)
	upgraderepo := res.UpgradeRepository(nil)

	res1ID := bson.NewObjectId()
	b1Value := dbmodels.ValueDB{res1ID, 100}
	doubler, _ := general.NewMathValScheme([]string{"*"}, []float64{2})
	doublerbyte, _ := json.Marshal(doubler)
	r1 := dbmodels.ResourceDB{
		ID:      res1ID,
		ResType: "regularbproduct",
	}

	b1tileeff := gamemap.NewCenteredAreaSpaceEffect(
		"TileAreaEffect",
		150,
		[][]gamemap.CenteredSEElementType{
			{gamemap.SpaceEmpty, gamemap.SpaceElement, gamemap.SpaceEmpty},           /*  initializers for row indexed by 0 */
			{gamemap.SpaceElement, gamemap.SpaceCenterElement, gamemap.SpaceElement}, /*  initializers for row indexed by 1 */
			{gamemap.SpaceEmpty, gamemap.SpaceElement, gamemap.SpaceEmpty},           /*  initializers for row indexed by 2 */
		}, []string{res1ID.Hex()}, doubler, nil,
	)

	b1dbeff, _ := dbmodels.FormFromEffect(b1tileeff)
	eff4 := buildinginstance.NewBuildingProductionEffect("eff4", 1, []string{res1ID.Hex()}, doubler, nil)
	u1dbeff, _ := dbmodels.FormFromEffect(eff4)
	building1ID := bson.NewObjectId()

	b1 := buildingdb.BuildingDB{
		ID:             building1ID,
		BID:            bson.NewObjectId(),
		BaseValue:      dbmodels.ValuesDB{b1Value},
		UpValScheme:    dbmodels.ValSchemeDB{"mathvalscheme", doublerbyte},
		UpStatScheme:   dbmodels.ValSchemeDB{"mathvalscheme", doublerbyte},
		Yields:         dbmodels.ValuesDB{b1Value},
		ProductionTime: 5 * time.Second.Nanoseconds(),
		SpaceEffect:    *b1dbeff,
		Ver: version.Version{
			Ver:          "1.0.0",
			VerDateAfter: time.Now(),
		},
	}

	compProperty := general.NewObjectProperties()
	compProperty.AddProperty("building", building1ID.Hex())
	comp := up.Complishment{
		ID:     "asd",
		Target: []up.ComplishmentTarget{up.ComplishmentTarget{compProperty}},
		Params: []up.ComplishmentParams{
			up.ComplishmentParams{
				InequalitySymbol: "=",
				Value:            1,
			}},
	}

	req := up.NewRequirement(
		"req1",
		[]up.Complishment{comp})

	targetProperty := general.NewObjectProperties()
	targetProperty.AddProperty("building", building1ID.Hex())
	/*up1 := up.NewUpgrade("upgradedoubleb1",
	nil,
	req,
	[]up.UpgradeTarget{
		up.UpgradeTarget{
			ObjectProperties: targetProperty}},
	[]effectdomain.IEffect{eff4})
	*/
	u1 := dbmodels.UpgradeDB{
		UpgradeID:          bson.NewObjectId(),
		UpgradeRequirement: req,
		Value:              nil,
		UpTargets: []up.UpgradeTarget{
			up.UpgradeTarget{
				ObjectProperties: targetProperty}},
		Effects: []dbmodels.EffectDB{*u1dbeff},
		Ver: version.Version{
			Ver:          "1.0.0",
			VerDateAfter: time.Now(),
		},
		CreationDate: time.Now(),
	}
	//fmt.Println(u1)
	panicIfFails(resrepo.CreateResource(&r1))
	panicIfFails(builrepo.CreateBuilding(&b1))
	panicIfFails(upgraderepo.CreateUpgrade(&u1))

	addresource := func(resource *usergame.Resource, w http.ResponseWriter, req *http.Request, payload *usergame.PostCreateGameHookPayload) error {
		payload.Ge.Resources = append(payload.Ge.Resources, dbmodels.ValueDB{res1ID, 101})
		return nil
	}
	res.ControllerHooks.PostCreateGameHook = append(res.ControllerHooks.PostCreateGameHook, addresource)
}

func panicIfFails(e error) {
	if e != nil {
		panic(e)
	}
}
