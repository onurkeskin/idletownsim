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


package server

import (
	"net/http"
	"time"

	"github.com/app/server/domain"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"gopkg.in/tylerb/graceful.v1"
)

// Request JSON body limit is set at 5MB (currently not enforced)
const BodyLimitBytes uint32 = 1048576 * 5

// Server type
type Server struct {
	negroni        *negroni.Negroni
	router         *Router
	gracefulServer *graceful.Server
	timeout        time.Duration
}

// Config type
type Config struct {
}

// Options for running the server
type Options struct {
	Timeout         time.Duration
	ShutdownHandler func()
	CertPath        string
	KeyPath         string
}

// NewServer Returns a new Server object
func NewServer(options *Config) *Server {

	// set up server and middlewares
	n := negroni.Classic()

	s := &Server{n, nil, nil, 0}

	return s
}

func (s *Server) UseMiddleware(middleware domain.IMiddleware) *Server {
	// next convert it into negroni style handlerfunc
	s.negroni.Use(negroni.HandlerFunc(middleware.Handler))
	return s
}

func (s *Server) UseRouter(router *Router) *Server {
	// add router and clear mux.context values at the end of request life-times
	s.negroni.UseHandler(context.ClearHandler(router))
	return s
}

func (s *Server) Run(address string, options Options) *Server {
	s.timeout = options.Timeout
	s.gracefulServer = &graceful.Server{
		Timeout:           options.Timeout,
		Server:            &http.Server{Addr: address, Handler: s.negroni},
		ShutdownInitiated: options.ShutdownHandler,
	}
	err := s.gracefulServer.ListenAndServeTLS(options.CertPath, options.KeyPath)
	if err != nil {
		panic(err)
	}
	return s
}

func (s *Server) Stop() {
	s.gracefulServer.Stop(s.timeout)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) *Server {
	s.negroni.ServeHTTP(w, r)
	return s
}
