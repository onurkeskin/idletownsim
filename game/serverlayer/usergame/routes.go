package game

import (
	"github.com/app/server/domain"
	"strings"
)

const (
	GetGame              = "GetGames"
	CreateGameByPos      = "CreateGame"
	CreateGameByMapID    = "CreateGameByMapID"
	GetEligibleBuildings = "GetEligibleBuildings"
	BuyBuilding          = "BuyBuilding"
	DeleteGameSession    = "DeleteGameSession"
	BuyUpgrade           = "BuyUpgrade"
	GetEligibleUpgrades  = "GetEligibleUpgrades"
	GetGames             = "GetGames"
)

const defaultBasePath = "/api/games"

func (resource *Resource) generateRoutes(basePath string) *domain.Routes {
	if basePath == "" {
		basePath = defaultBasePath
	}
	var baseRoutes = domain.Routes{

		domain.Route{
			Name:           GetGames,
			Method:         "GET",
			Pattern:        "/api/games",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.HandleGetGames_v0,
			},
			ACLHandler: resource.HandleGetGamesACL,
		},
		domain.Route{
			Name:           GetGame,
			Method:         "GET",
			Pattern:        "/api/game/{id}",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.HandleGetGame_v0,
			},
			ACLHandler: resource.HandleGetGameACL,
		},
		domain.Route{
			Name:           CreateGameByPos,
			Method:         "POST",
			Pattern:        "/api/game/bypos",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.HandleCreateGameByPos_v0,
			},
			ACLHandler: resource.HandleCreateGameACL,
		},
		domain.Route{
			Name:           CreateGameByMapID,
			Method:         "POST",
			Pattern:        "/api/game/bymid",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.HandleCreateGameByMapID_v0,
			},
			ACLHandler: resource.HandleCreateGameACL,
		},
		domain.Route{
			Name:           GetEligibleBuildings,
			Method:         "GET",
			Pattern:        "/api/game/{id}/buildings",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.HandleGetEligibleBuildings_v0,
			},
			ACLHandler: resource.HandleGetEligibleBuildingsACL,
		},
		domain.Route{
			Name:           BuyBuilding,
			Method:         "POST",
			Pattern:        "/api/game/{id}/space/{sid}/building/buy/{bid}",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.HandleBuyBuilding_v0,
			},
			ACLHandler: resource.HandleBuyBuildingACL,
		},
		domain.Route{
			Name:           GetEligibleUpgrades,
			Method:         "GET",
			Pattern:        "/api/game/{id}/upgrades",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.HandleGetEligibleUpgrades_v0,
			},
			ACLHandler: resource.HandleGetEligibleUpgradesACL,
		},
		domain.Route{
			Name:           BuyUpgrade,
			Method:         "POST",
			Pattern:        "/api/game/{id}/upgrade/buy/{uid}",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.HandleBuyUpgrade_v0,
			},
			ACLHandler: resource.HandleBuyUpgradeACL,
		},
		domain.Route{
			Name:           DeleteGameSession,
			Method:         "DELETE",
			Pattern:        "/api/games",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.DeleteGameSessions_v0,
			},
			ACLHandler: resource.DeleteGameSessionsACL,
		},
	}

	routes := domain.Routes{}

	for _, route := range baseRoutes {
		r := domain.Route{
			Name:           route.Name,
			Method:         route.Method,
			Pattern:        strings.Replace(route.Pattern, defaultBasePath, basePath, -1),
			DefaultVersion: route.DefaultVersion,
			RouteHandlers:  route.RouteHandlers,
			ACLHandler:     route.ACLHandler,
		}
		routes = routes.Append(&domain.Routes{r})
	}
	resource.routes = &routes
	return resource.routes
}
