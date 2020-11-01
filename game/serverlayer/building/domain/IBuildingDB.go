package building

import (
	buildingdomain "github.com/app/game/applayer/building/domain"
	resourcedomain "github.com/app/game/serverlayer/gameresources/domain"
)

type IBuildingsDB interface{}

type IBuildingDB interface {
	FormIBuilding(provider resourcedomain.IResourceRepository) buildingdomain.IBuilding
}
