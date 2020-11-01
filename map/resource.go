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
	mapdomain "github.com/app/map/domain"
	"github.com/app/server/domain"
	"net/http"
)

type PostCreateMapHookPayload struct {
	//User domain.IUser
}

type ControllerHooks struct {
	PostCreateMapHook func(resource *Resource, w http.ResponseWriter, req *http.Request, payload *PostCreateMapHookPayload) error
}

type Options struct {
	BasePath             string
	Database             domain.IDatabase
	Renderer             domain.IRenderer
	MapRepositoryFactory mapdomain.IMapRepositoryFactory
	MapProvider          mapdomain.IMapProvider
	ControllerHooks      *ControllerHooks
}

func NewResource(options *Options) *Resource {

	database := options.Database
	if database == nil {
		panic("maps.Options.Database is required")
	}
	renderer := options.Renderer
	if renderer == nil {
		panic("maps.Options.Renderer is required")
	}

	mapRepositoryFactory := options.MapRepositoryFactory
	if mapRepositoryFactory == nil {
		// init default UserRepositoryFactory
		mapRepositoryFactory = NewMapRepositoryFactory()
	}

	mapProvider := options.MapProvider
	if mapProvider == nil {
		mapProvider = &MapProvider{Retriever: MapImageByGoogleMaps}
	}

	controllerHooks := options.ControllerHooks
	if controllerHooks == nil {
		controllerHooks = &ControllerHooks{}
	}

	u := &Resource{options, nil,
		database,
		renderer,
		mapRepositoryFactory,
		mapProvider,
		controllerHooks,
	}
	u.generateRoutes(options.BasePath)
	return u
}

// UsersResource implements IResource
type Resource struct {
	options              *Options
	routes               *domain.Routes
	Database             domain.IDatabase
	Renderer             domain.IRenderer
	MapRepositoryFactory mapdomain.IMapRepositoryFactory
	MapProvider          mapdomain.IMapProvider
	ControllerHooks      *ControllerHooks
}

func (resource *Resource) Routes() *domain.Routes {
	return resource.routes
}

func (resource *Resource) Render(w http.ResponseWriter, req *http.Request, status int, v interface{}) {
	resource.Renderer.Render(w, req, status, v)
}

func (resource *Resource) MapRepository(req *http.Request) mapdomain.IMapRepository {
	return resource.MapRepositoryFactory.New(resource.Database)
}

func (resource *Resource) GetMapProvider() mapdomain.IMapProvider {
	return resource.MapProvider
}
