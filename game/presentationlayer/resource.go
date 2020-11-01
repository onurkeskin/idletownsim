package presentationlayer

import (
	"github.com/app/game/presentationlayer/buildingview"
	buildingviewdomain "github.com/app/game/presentationlayer/buildingview/domain"
	"github.com/app/game/presentationlayer/resourceview"
	resourceviewdomain "github.com/app/game/presentationlayer/resourceview/domain"
	"github.com/app/server/domain"
	"net/http"
)

type Options struct {
	Database                      domain.IDatabase
	Renderer                      domain.IRenderer
	ResourceViewRepositoryFactory resourceviewdomain.IResourceViewRepositoryFactory
	BuildingViewRepositoryFactory buildingviewdomain.IBuildingViewRepositoryFactory
}

func NewResource(options *Options) *Resource {

	database := options.Database
	if database == nil {
		panic("users.Options.Database is required")
	}
	renderer := options.Renderer
	if renderer == nil {
		panic("users.Options.Renderer is required")
	}

	resourceViewRepositoryFactory := options.ResourceViewRepositoryFactory
	if resourceViewRepositoryFactory == nil {
		// init default UserRepositoryFactory
		resourceViewRepositoryFactory = resourceview.NewResourceViewRepositoryFactory()
	}

	buildingViewRepositoryFactory := options.BuildingViewRepositoryFactory
	if buildingViewRepositoryFactory == nil {
		// init default UserRepositoryFactory
		buildingViewRepositoryFactory = buildingview.NewBuildingViewRepositoryFactory()
	}

	u := &Resource{options,
		database,
		renderer,
		nil,
		resourceViewRepositoryFactory,
		buildingViewRepositoryFactory,
	}

	viewCreator := ViewCreator{
		brep: u.BuildingViewRepository(nil),
		rrep: u.ResourceViewRepository(nil),
	}
	u.ViewCreator = &viewCreator

	return u
}

type Resource struct {
	options                       *Options
	Database                      domain.IDatabase
	Renderer                      domain.IRenderer
	ViewCreator                   *ViewCreator
	ResourceViewRepositoryFactory resourceviewdomain.IResourceViewRepositoryFactory
	BuildingViewRepositoryFactory buildingviewdomain.IBuildingViewRepositoryFactory
}

func (resource *Resource) Render(w http.ResponseWriter, req *http.Request, status int, v interface{}) {
	resource.Renderer.Render(w, req, status, v)
}

func (resource *Resource) BuildingViewRepository(req *http.Request) buildingviewdomain.IBuildingViewRepository {
	return resource.BuildingViewRepositoryFactory.New(resource.Database)
}
func (resource *Resource) ResourceViewRepository(req *http.Request) resourceviewdomain.IResourceViewRepository {
	return resource.ResourceViewRepositoryFactory.New(resource.Database)
}

func (resource *Resource) GetViewCreator() *ViewCreator {
	return resource.ViewCreator
}
