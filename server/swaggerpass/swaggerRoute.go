package swaggerpass

import (
	"github.com/app/server/domain"
	"strings"
)

type Swagger struct {
}

const (
	SwaggerPort = "swagger"
)

const defaultBasePath = "/api/swagger.json"

func (swagger Swagger) GenerateRoutes(basePath string) *domain.Routes {

	if basePath == "" {
		basePath = defaultBasePath
	}
	var baseRoutes = domain.Routes{
		domain.Route{
			Name:           SwaggerPort,
			Method:         "GET",
			Pattern:        "/api/swagger.json",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": swagger.HandleSwagger_v0,
			},
			ACLHandler: swagger.HandleSwaggerACL,
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

	return &routes
}
