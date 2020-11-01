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


package server_test

import (
	"errors"
	"fmt"
	buildingdb "github.com/app/game/serverlayer/building"
	upgradedb "github.com/app/game/serverlayer/upgrade"
	usergameres "github.com/app/game/serverlayer/usergame"
	mapdb "github.com/app/map"
	"github.com/app/server/middlewares/context"
	"github.com/app/server/middlewares/mongodb"
	"github.com/app/server/middlewares/renderer"
	"github.com/app/server/server"
	"github.com/app/sessions"
	"github.com/app/users"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

var _ = Describe("Server", func() {
	// try to load signing keys for token authority
	// NOTE: DO NOT USE THESE KEYS FOR PRODUCTION! FOR TEST ONLY
	privateSigningKey, err := ioutil.ReadFile("../keys/demo.rsa")
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error loading private signing key: %v", err.Error())))
	}
	publicSigningKey, err := ioutil.ReadFile("../keys/demo.rsa.pub")
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error loading public signing key: %v", err.Error())))
	}

	Describe("Basic sanity test", func() {
		ctx := context.New()

		db := mongodb.New(&mongodb.Options{
			ServerName:   "localhost",
			DatabaseName: "test-go-app",
		})
		_ = db.NewSession()

		// init renderer
		renderer := renderer.New(&renderer.Options{
			IndentJSON: true,
		}, renderer.JSON)

		// set up users resource
		usersResource := users.NewResource(ctx, &users.Options{
			Database: db,
			Renderer: renderer,
		})

		// set up sessions resource
		sessionsResource := sessions.NewResource(ctx, &sessions.Options{
			Database:              db,
			Renderer:              renderer,
			PrivateSigningKey:     privateSigningKey,
			PublicSigningKey:      publicSigningKey,
			UserRepositoryFactory: usersResource.UserRepositoryFactory,
		})

		BuildingFactory := buildingdb.NewBuildingRepositoryFactory()
		UpgradeFactory := upgradedb.NewUpgradeRepositoryFactory()
		MapFactory := mapdb.NewMapRepositoryFactory()

		usergameres.NewResource(ctx, &usergameres.Options{
			Database:                  db,
			Renderer:                  renderer,
			UpgradeRepositoryFactory:  UpgradeFactory,
			BuildingRepositoryFactory: BuildingFactory,
			MapRepositoryFactory:      MapFactory,
		})

		// init server
		s := server.NewServer(&server.Config{
			Context: ctx,
		})

		// set up router
		ac := server.NewAccessController(ctx, renderer)
		router := server.NewRouter(s.Context, ac)

		// add REST resources to router
		router.AddResources(sessionsResource, usersResource)

		// add middlewares
		s.UseContextMiddleware(renderer)
		s.UseMiddleware(sessionsResource.NewAuthenticator())

		s.UseContextMiddleware(renderer)

		s.UseRouter(router)

		It("should serve request", func() {
			// run server and it shouldn't panic
			go s.Run(":8001", server.Options{
				Timeout: 1 * time.Millisecond,
			})
			time.Sleep(100 * time.Millisecond)

			// serve some urls
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest("GET", "/api/sessions", nil)

			s.ServeHTTP(recorder, request)

			request2, _ := http.NewRequest("GET", "/api/game", nil)
			s.ServeHTTP(recorder, request2)
			Expect(recorder.Code).To(Equal(http.StatusForbidden))

			// gracefully stops server
			s.Stop()

		})
	})
})
