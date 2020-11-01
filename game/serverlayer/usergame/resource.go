package game

import (
	"encoding/json"
	"github.com/app/game/applayer/buildinginstance"
	"github.com/app/game/applayer/gamemap"
	"github.com/app/game/applayer/general"
	up "github.com/app/game/applayer/upgrade"
	"github.com/app/game/presentationlayer"
	buildingdb "github.com/app/game/serverlayer/building"
	serverlayerbuilding "github.com/app/game/serverlayer/building"
	serverlayerbuildingdomain "github.com/app/game/serverlayer/building/domain"
	"github.com/app/game/serverlayer/dbmodels"
	serverlayerresourcedomain "github.com/app/game/serverlayer/gameresources/domain"
	serverlayerspacedomain "github.com/app/game/serverlayer/space/domain"
	serverlayerupgradedomain "github.com/app/game/serverlayer/upgrade/domain"
	serverlayergamedomain "github.com/app/game/serverlayer/usergame/domain"
	"github.com/app/helpers/version"
	serverlayermapdomain "github.com/app/map/domain"
	"github.com/app/server/domain"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

type PostCreateGameHookPayload struct {
	Ge *dbmodels.GameEnvironmentDB
}
type ControllerHooks struct {
	PostCreateGameHook []func(resource *Resource, w http.ResponseWriter, req *http.Request, payload *PostCreateGameHookPayload) error
}

type Options struct {
	BasePath                  string
	Database                  domain.IDatabase
	Renderer                  domain.IRenderer
	GameRepositoryFactory     serverlayergamedomain.IGameRepositoryFactory
	UpgradeRepositoryFactory  serverlayerupgradedomain.IUpgradeRepositoryFactory
	BuildingRepositoryFactory serverlayerbuildingdomain.IBuildingRepositoryFactory
	MapRepositoryFactory      serverlayermapdomain.IMapRepositoryFactory
	ResourceRepositoryFactory serverlayerresourcedomain.IResourceRepositoryFactory
	SpaceRepositoryFactory    serverlayerspacedomain.ISpaceRepositoryFactory
	MapProvider               serverlayermapdomain.IMapProvider
	GameManager               serverlayergamedomain.IGameManager
	ViewResource              *presentationlayer.Resource
	ControllerHooks           *ControllerHooks
}

func NewResource(options *Options) *Resource {

	database := options.Database
	if database == nil {
		panic("game.Options.Database is required")
	}
	renderer := options.Renderer
	if renderer == nil {
		panic("game.Options.Renderer is required")
	}
	MapRepositoryFactory := options.MapRepositoryFactory
	if MapRepositoryFactory == nil {
		panic("maps.options.MapRepositoryFactory is required")
	}
	UpgradeRepositoryFactory := options.UpgradeRepositoryFactory
	if UpgradeRepositoryFactory == nil {
		panic("upgrades.options.UpgradeRepositoryFactory is required")
	}
	BuildingRepositoryFactory := options.BuildingRepositoryFactory
	if BuildingRepositoryFactory == nil {
		panic("building.options.BuildingRepositoryFactory is required")
	}
	ResourceRepositoryFactory := options.ResourceRepositoryFactory
	if ResourceRepositoryFactory == nil {
		panic("resource.options.BuildingRepositoryFactory is required")
	}
	SpaceRepositoryFactory := options.SpaceRepositoryFactory
	if SpaceRepositoryFactory == nil {
		panic("resource.options.SpaceRepositoryFactory is required")
	}

	mapProvider := options.MapProvider
	if mapProvider == nil {
		panic("maps.options.MapProvider is required")
	}

	gameRepositoryFactory := options.GameRepositoryFactory
	if gameRepositoryFactory == nil {
		// init default RevokedTokenRepositoryFactory
		gameRepositoryFactory = NewGameRepositoryFactory()
	}

	buildingMapper := serverlayerbuilding.NewBuildingMapper(&serverlayerbuilding.Options{database, BuildingRepositoryFactory, ResourceRepositoryFactory})

	gameManager := options.GameManager
	if gameManager == nil {
		gameManager = InitManager(10 * time.Minute)
	}

	controllerHooks := options.ControllerHooks
	if controllerHooks == nil {
		controllerHooks = &ControllerHooks{}
	}

	viewResource := options.ViewResource
	if options.ViewResource == nil {
		viewResource = presentationlayer.NewResource(
			&presentationlayer.Options{
				Database: database,
				Renderer: renderer,
			})
	}

	resource := &Resource{options, nil,
		database,
		renderer,
		gameRepositoryFactory,
		UpgradeRepositoryFactory,
		BuildingRepositoryFactory,
		ResourceRepositoryFactory,
		MapRepositoryFactory,
		SpaceRepositoryFactory,
		*buildingMapper,
		mapProvider,
		gameManager,
		viewResource,
		controllerHooks,
	}

	resource.generateRoutes(options.BasePath)
	return resource
}

// gameResource implements IResource
type Resource struct {
	options                   *Options
	routes                    *domain.Routes
	Database                  domain.IDatabase
	Renderer                  domain.IRenderer
	GameRepositoryFactory     serverlayergamedomain.IGameRepositoryFactory
	UpgradeRepositoryFactory  serverlayerupgradedomain.IUpgradeRepositoryFactory
	BuildingRepositoryFactory serverlayerbuildingdomain.IBuildingRepositoryFactory
	ResourceRepositoryFactory serverlayerresourcedomain.IResourceRepositoryFactory
	MapRepositoryFactory      serverlayermapdomain.IMapRepositoryFactory
	SpaceRepositoryFactory    serverlayerspacedomain.ISpaceRepositoryFactory
	BuildingMap               serverlayerbuilding.BuildingMapper
	MapProvider               serverlayermapdomain.IMapProvider
	GameManager               serverlayergamedomain.IGameManager
	ViewResource              *presentationlayer.Resource
	ControllerHooks           *ControllerHooks
}

func (resource *Resource) Routes() *domain.Routes {
	return resource.routes
}

func (resource *Resource) Render(w http.ResponseWriter, req *http.Request, status int, v interface{}) {
	resource.Renderer.Render(w, req, status, v)
}

func (resource *Resource) GameRepository(req *http.Request) serverlayergamedomain.IGameRepository {
	return resource.GameRepositoryFactory.New(resource.Database)
}

func (resource *Resource) UpgradeRepository(req *http.Request) serverlayerupgradedomain.IUpgradeRepository {
	return resource.UpgradeRepositoryFactory.New(resource.Database)
}

func (resource *Resource) MapRepository(req *http.Request) serverlayermapdomain.IMapRepository {
	return resource.MapRepositoryFactory.New(resource.Database)
}

func (resource *Resource) BuildingRepository(req *http.Request) serverlayerbuildingdomain.IBuildingRepository {
	return resource.BuildingRepositoryFactory.New(resource.Database)
}

func (resource *Resource) ResourceRepository(req *http.Request) serverlayerresourcedomain.IResourceRepository {
	return resource.ResourceRepositoryFactory.New(resource.Database)
}

func (resource *Resource) SpaceRepository(req *http.Request) serverlayerspacedomain.ISpaceRepository {
	return resource.SpaceRepositoryFactory.New(resource.Database)
}

func (resource *Resource) BuildingMapper(req *http.Request) serverlayerbuilding.BuildingMapper {
	return resource.BuildingMap
}

func (resource *Resource) ViewRes(req *http.Request) *presentationlayer.Resource {
	return resource.ViewResource
}

/*
func (resource *Resource) UserRepository(req *http.Request) usersDomain.IUserRepository {
	return resource.UserRepositoryFactory.New(resource.Database)
}
*/

func (res *Resource) InsertTestValuesIntoResource() {
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
		ID:  building1ID,
		BID: bson.NewObjectId(),
		//BuildingGlobalIdentifier: "factorylevel1",
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

	addresource := func(resource *Resource, w http.ResponseWriter, req *http.Request, payload *PostCreateGameHookPayload) error {
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
