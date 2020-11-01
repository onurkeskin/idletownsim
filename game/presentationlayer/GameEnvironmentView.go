package presentationlayer

import (
	buildingviewdomain "github.com/app/game/presentationlayer/buildingview/domain"
	effectviewdomain "github.com/app/game/presentationlayer/effectview/domain"
	resourceviewdomain "github.com/app/game/presentationlayer/resourceview/domain"

	"time"
)

type GamesView []GameView

type GameView struct {
	GameID            string                           `json:"id,omitempty"`
	GameCreationTime  time.Time                        `json:"gamect"`
	LastIterationTime time.Time                        `json:"gamelitime"`
	Resources         resourceviewdomain.IResourceView `json:"resources"`
	Spaces            SpacesView                       `json:"spaces"`
}

type SpacesView []SpaceView

type SpaceView struct {
	SpaceID          string                           `json:"id,omitempty" bson:"_id,omitempty"`
	InMapID          string                           `json:"inmapid" bson:"inmapid"`
	ExpectedResource resourceviewdomain.IResourceView `json:"expectedincome"`
	AppliedEffects   effectviewdomain.IEffectView     `json:"effectson"`
	Occupier         buildingviewdomain.IBuildingView `json:"spaceoccupier"`
}
