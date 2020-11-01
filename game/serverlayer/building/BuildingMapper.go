package building

import (
	building "github.com/app/game/applayer/building"
	//gamemapdomain "github.com/app/game/applayer/gamemap/domain"
	//generaldomain "github.com/app/game/applayer/general/domain"
	serverlayerbuildingdomain "github.com/app/game/serverlayer/building/domain"
	serverlayerresourcedomain "github.com/app/game/serverlayer/gameresources/domain"
	"github.com/app/server/domain"
	"net/http"
)

type Options struct {
	Database                  domain.IDatabase
	BuildingRepositoryFactory serverlayerbuildingdomain.IBuildingRepositoryFactory
	ResourceRepositoryFactory serverlayerresourcedomain.IResourceRepositoryFactory
}

type BuildingMapper struct {
	Database                  domain.IDatabase
	BuildingRepositoryFactory serverlayerbuildingdomain.IBuildingRepositoryFactory
	ResourceRepositoryFactory serverlayerresourcedomain.IResourceRepositoryFactory

	mapped map[string]serverlayerbuildingdomain.IBuildingDB
}

func NewBuildingMapper(options *Options) *BuildingMapper {
	database := options.Database
	if database == nil {
		panic("game.Options.Database is required")
	}

	BuildingRepositoryFactory := options.BuildingRepositoryFactory
	if BuildingRepositoryFactory == nil {
		panic("building.options.BuildingRepositoryFactory is required")
	}
	ResourceRepositoryFactory := options.ResourceRepositoryFactory
	if ResourceRepositoryFactory == nil {
		panic("resource.options.BuildingRepositoryFactory is required")
	}

	manager := &BuildingMapper{
		Database:                  database,
		BuildingRepositoryFactory: BuildingRepositoryFactory,
		ResourceRepositoryFactory: ResourceRepositoryFactory,
		mapped: make(map[string]serverlayerbuildingdomain.IBuildingDB),
	}

	return manager
}
func (r *BuildingMapper) GetBuildingDBMapFor(id string) *BuildingDB {
	v, ok := r.mapped[id]
	if ok {
		return v.(*BuildingDB)
	}

	builRepo := r.BuildingRepository(nil)
	//resRepo := r.ResourceRepository(nil)
	builDB, err := builRepo.GetBuildingById(id)
	if err != nil {
		return nil
	}
	//FORM UPSCHEME
	//	buil := builDB.FormIBuilding(resRepo)
	r.mapped[id] = builDB
	return builDB.(*BuildingDB)
}

func (r *BuildingMapper) GetBuildingMapFor(id string) *building.Building {
	builRepo := r.BuildingRepository(nil)
	resRepo := r.ResourceRepository(nil)
	v, ok := r.mapped[id]
	if ok {
		return v.FormIBuilding(resRepo).(*building.Building)
	}

	builDB, err := builRepo.GetBuildingById(id)
	if err != nil {
		return nil
	}
	//FORM UPSCHEME
	buil := builDB.FormIBuilding(resRepo)
	r.mapped[id] = builDB
	return buil.(*building.Building)
}

func (r *BuildingMapper) GetBuildingsDBMapFor(bid string) BuildingsDB {
	builRepo := r.BuildingRepository(nil)

	_builsDB, err := builRepo.GetBuildingsByBuildingId(bid)
	if err != nil {
		return nil
	}
	builsDB := *_builsDB.(*BuildingsDB)

	for _, v := range builsDB {
		r.mapped[v.ID.Hex()] = &v
	}
	//FORM UPSCHEME

	return builsDB
}

func (r *BuildingMapper) GetMappedBuildings() BuildingsDB {
	builRepo := r.BuildingRepository(nil)
	_builsDB := builRepo.GetBuildings()
	return *_builsDB.(*BuildingsDB)
}

func (r *BuildingMapper) MapBuildingDB(b BuildingDB) *building.Building {
	resRepo := r.ResourceRepository(nil)
	return b.FormIBuilding(resRepo).(*building.Building)
}

func (mapper *BuildingMapper) BuildingRepository(req *http.Request) serverlayerbuildingdomain.IBuildingRepository {
	return mapper.BuildingRepositoryFactory.New(mapper.Database)
}

func (mapper *BuildingMapper) ResourceRepository(req *http.Request) serverlayerresourcedomain.IResourceRepository {
	return mapper.ResourceRepositoryFactory.New(mapper.Database)
}
