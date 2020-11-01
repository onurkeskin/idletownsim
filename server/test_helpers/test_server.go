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


package test_helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/app/server/domain"
	"github.com/app/server/middlewares/context"
	"github.com/app/server/middlewares/mongodb"
	"github.com/app/server/middlewares/renderer"
	"github.com/app/server/server"
	"github.com/app/sessions"
	sessionsDomain "github.com/app/sessions/domain"
	"net/http"
	"net/http/httptest"
)

type TestServerOptions struct {
	RequestAcceptHeader string
	PrivateSigningKey   []byte
	PublicSigningKey    []byte
	TokenAuthority      sessionsDomain.ITokenAuthority
	Database            domain.IDatabase
	Renderer            domain.IRenderer
	Resources           []domain.IResource
	Middlewares         []interface{}
}

type TestServer struct {
	Options        *TestServerOptions
	Server         *server.Server
	Router         *server.Router
	TokenAuthority sessionsDomain.ITokenAuthority
	Database       domain.IDatabase
	Renderer       domain.IRenderer
}

type AuthOptions struct {
	APIUser domain.IUser
	Token   string
}

func NewTestServer(options *TestServerOptions) *TestServer {

	// set up basic needs for a test server

	ta := options.TokenAuthority
	if ta == nil {
		if options.PrivateSigningKey == nil {
			panic("TestServer.options.PrivateSigningKey is required")
		}
		if options.PublicSigningKey == nil {
			panic("TestServer.options.PublicSigningKey is required")
		}
		ta = sessions.NewTokenAuthority(&sessions.TokenAuthorityOptions{
			PrivateSigningKeyByte: options.PrivateSigningKey,
			PublicSigningKeyByte:  options.PublicSigningKey,
		})
	}

	ctx := context.New()

	// init server
	s := server.NewServer(&server.Config{
		Context: ctx,
	})

	// set up DB session if not specified
	db := options.Database
	if options.Database == nil {
		db = mongodb.New(&mongodb.Options{
			ServerName:   "localhost",
			DatabaseName: "test-go-app",
		})
		_ = db.(*mongodb.MongoDB).NewSession()
	}

	// set up Renderer (unrolled_render)
	r := options.Renderer
	if r == nil {
		r = renderer.New(&renderer.Options{
			IndentJSON: true,
		}, renderer.JSON)
	}

	// set up router
	ac := server.NewAccessController(ctx, r)
	router := server.NewRouter(s.Context, ac)

	// init test server
	ts := TestServer{options, s, router, ta, db, r}

	// add REST resources to router
	for _, resource := range options.Resources {
		ts.AddResources(resource)
	}

	// add middlewares
	for _, middleware := range options.Middlewares {
		ts.AddMiddlewares(middleware)
	}

	return &ts
}

func (ts *TestServer) AddResources(resources ...domain.IResource) {
	for _, resource := range resources {
		ts.Router.AddResources(resource)
	}
}
func (ts *TestServer) AddMiddlewares(middlewares ...interface{}) {
	for _, middleware := range middlewares {
		switch v := middleware.(type) {
		case domain.IMiddleware:
			ts.Server.UseMiddleware(middleware.(domain.IMiddleware))
		case domain.IContextMiddleware:
			ts.Server.UseContextMiddleware(middleware.(domain.IContextMiddleware))
		default:
			fmt.Println("Unknown middleware, skipping", v)
		}
	}
}
func (ts *TestServer) Run() {
	ts.Server.UseRouter(ts.Router)
}
func (ts *TestServer) Request(method string, urlStr string, body interface{}, targetResponse interface{}, authOptions *AuthOptions) *httptest.ResponseRecorder {

	recorder := httptest.NewRecorder()

	var request *http.Request

	// request for version 0.0
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		request, _ = http.NewRequest(method, urlStr, bytes.NewReader(jsonBytes))
	} else {
		request, _ = http.NewRequest(method, urlStr, nil)
	}
	// set API version through accept header
	request.Header.Set("Accept", ts.Options.RequestAcceptHeader)

	if authOptions == nil {
		authOptions = &AuthOptions{nil, ""}
	}
	if authOptions.APIUser != nil {
		// set Authorization header
		token, _ := ts.TokenAuthority.CreateNewSessionToken(&sessions.TokenClaims{
			UserID: authOptions.APIUser.GetID(),
		})
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	} else {
		if authOptions.Token != "" {
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", authOptions.Token))
		}
	}
	// serve request
	ts.Server.ServeHTTP(recorder, request)
	DecodeResponseToType(recorder, &targetResponse)
	return recorder
}
