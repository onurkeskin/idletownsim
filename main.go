package main

import (
	"encoding/json"
	"fmt"
	"github.com/app/appenginehelpers"
	buildingdb "github.com/app/game/serverlayer/building"
	resourcedb "github.com/app/game/serverlayer/gameresources"
	spacedb "github.com/app/game/serverlayer/space"
	upgradedb "github.com/app/game/serverlayer/upgrade"
	usergame "github.com/app/game/serverlayer/usergame"
	mapdb "github.com/app/map"
	"github.com/app/server/middlewares/mongodb"
	"github.com/app/server/middlewares/renderer"
	"github.com/app/server/server"
	"github.com/app/server/swaggerpass"
	"github.com/app/sessions"
	"github.com/app/users"
	"github.com/app/users/hooks"
	"os"
	"time"
)

//serviceppbucket for bucket

var (
	// Version is a compile time constant, injected at build time
	Version string
)

var Configuration GameConfiguration

type GameConfiguration struct {
	Reset_db bool `json:"reset_db"`
	Setup_db bool `json:"setup_db"`
}

func main() {
	testconffile, e := os.Open("./gameconf.json")
	if e != nil {
		fmt.Printf("Game Config couldnt read error: %v\n", e)
	}
	decoder := json.NewDecoder(testconffile)
	panicIfFails(decoder.Decode(&Configuration))

	// try to load signing keys for token authority
	// NOTE: DO NOT USE THESE KEYS FOR PRODUCTION! FOR DEMO ONLY
	prov := appengessentials.GetInstance()
	mapKey := appengessentials.GetKeyFromFile("MapsStaticKey", "keys/googlestaticmapskey")
	prov.SetKeyFor(mapKey[0], mapKey[1])

	privateSigningKey := appengessentials.ReadFileOnBucket("demokeys.rsa")
	publicSigningKey := appengessentials.ReadFileOnBucket("demokeys.rsa.pub")
	// create current project context

	// set up DB session
	db := mongodb.New(&mongodb.Options{
		ServerName:   "dbname/////",
		DatabaseName: "try",
	})
	_ = db.NewSession()

	// set up Renderer (unrolled_render)
	renderer := renderer.New(&renderer.Options{
		IndentJSON: true,
	}, renderer.JSON)

	controllerHooks := users.ControllerHooks{
		PostCreateUserHook:  userhooks.PostCreateUserHookF,
		PostConfirmUserHook: nil,
	}

	// set up users resource
	usersResource := users.NewResource(&users.Options{
		Database:        db,
		Renderer:        renderer,
		ControllerHooks: &controllerHooks,
	})

	// set up sessions resource
	sessionsResource := sessions.NewResource(&sessions.Options{
		PrivateSigningKey:     privateSigningKey,
		PublicSigningKey:      publicSigningKey,
		Database:              db,
		Renderer:              renderer,
		UserRepositoryFactory: usersResource.UserRepositoryFactory,
	})

	BuildingFactory := buildingdb.NewBuildingRepositoryFactory()
	ResourceFactory := resourcedb.NewResourceRepositoryFactory()
	UpgradeFactory := upgradedb.NewUpgradeRepositoryFactory()
	SpaceFactory := spacedb.NewSpaceRepositoryFactory()

	mapResource := mapdb.NewResource(&mapdb.Options{
		Database: db,
		Renderer: renderer,
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

	// init server
	s := server.NewServer(&server.Config{})

	// set up router
	ac := server.NewAccessController(renderer)
	router := server.NewRouter(ac)

	// add REST resources to router
	router.AddResources(sessionsResource, usersResource, usergameResource, mapResource)

	// add Swagger pass
	swag := swaggerpass.Swagger{}
	router.AddRoutes(swag.GenerateRoutes(""))

	// add middlewares
	s.UseMiddleware(sessionsResource.NewAuthenticator())

	// setup router
	s.UseRouter(router)

	if Configuration.Reset_db {
		db.DropDatabase()
		fmt.Println("Database Initiated")
	}

	if Configuration.Setup_db {
		usergameResource.InsertTestValuesIntoResource()
		fmt.Println("Resources Initiated")
	}

	CertPath := "keys/server.pem"
	KeyPath := "keys/server.key"
	// bam!
	s.Run(":443", server.Options{
		Timeout:  10 * time.Second,
		CertPath: CertPath,
		KeyPath:  KeyPath,
	})
}

func panicIfFails(e error) {
	if e != nil {
		panic(e)
	}
}
