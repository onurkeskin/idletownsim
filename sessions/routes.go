package sessions

import (
	"github.com/app/server/domain"
)

const (
	GetSession    = "GetSession"
	CreateSession = "CreateSession"
	DeleteSession = "DeleteSession"
)

func (resource *Resource) generateRoutes(basePath string) *domain.Routes {
	if basePath == "" {
		basePath = "/api/sessions"
	}
	var baseRoutes = domain.Routes{

		domain.Route{
			Name:           GetSession,
			Method:         "GET",
			Pattern:        "/api/sessions",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.HandleGetSession_v0,
			},
			ACLHandler: resource.HandleGetSessionACL,
		},
		domain.Route{
			Name:           CreateSession,
			Method:         "POST",
			Pattern:        "/api/sessions",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.HandleCreateSession_v0,
			},
			ACLHandler: resource.HandleCreateSessionACL,
		},
		domain.Route{
			Name:           DeleteSession,
			Method:         "DELETE",
			Pattern:        "/api/sessions",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.HandleDeleteSession_v0,
			},
			ACLHandler: resource.HandleDeleteSessionACL,
		},
	}

	routes := domain.Routes{}

	for _, route := range baseRoutes {
		r := domain.Route{
			Name:           route.Name,
			Method:         route.Method,
			Pattern:        basePath,
			DefaultVersion: route.DefaultVersion,
			RouteHandlers:  route.RouteHandlers,
			ACLHandler:     route.ACLHandler,
		}
		routes = routes.Append(&domain.Routes{r})
	}
	resource.routes = &routes
	return resource.routes
}
