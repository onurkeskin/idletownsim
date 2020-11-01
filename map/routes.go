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


package maps

import (
	"github.com/app/server/domain"
	"strings"
)

const (
	CreateMap    = "CreateMap"
	GetMap       = "GetMap"
	GetCloseMaps = "GetCloseMaps"
)

const defaultBasePath = "/api/maps"

func (resource *Resource) generateRoutes(basePath string) *domain.Routes {
	if basePath == "" {
		basePath = defaultBasePath
	}
	var baseRoutes = domain.Routes{
		domain.Route{
			Name:           GetMap,
			Method:         "GET",
			Pattern:        "/api/maps/{id}",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.HandleGetMap_v0,
			},
			ACLHandler: resource.HandleGetMapACL,
		},
		domain.Route{
			Name:           GetCloseMaps,
			Method:         "GET",
			Pattern:        "/api/maps",
			Queries:        []string{"lat", "lng"},
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.HandleGetCloseMaps_v0,
			},
			ACLHandler: resource.HandleGetCloseMapsACL,
		},
		domain.Route{
			Name:           CreateMap,
			Method:         "POST",
			Pattern:        "/api/maps",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.HandleCreateMapLatLng_v0,
			},
			ACLHandler: resource.HandleCreateMapLatLngACL,
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
