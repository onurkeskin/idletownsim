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


package domain_test

import (
	"github.com/app/server/domain"
	"github.com/app/server/middlewares/context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Middleware Tests", func() {
	Describe("ContextHandlerFunc Type", func() {
		Describe("ServeHTTP()", func() {
			It("should be working", func() {
				recorder := httptest.NewRecorder()
				request, _ := http.NewRequest("GET", "/api/test", nil)
				ctx := context.New()

				var handler domain.ContextHandlerFunc
				handler = func(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
					Expect(req.URL.Path).To(Equal("/api/test"))
					Expect(req.Method).To(Equal("GET"))
				}
				handler.ServeHTTP(recorder, request, ctx)

			})

		})

	})

	Describe("MiddlewareFunc Type", func() {
		Describe("ServeHTTP()", func() {
			It("should be working", func() {
				recorder := httptest.NewRecorder()
				request, _ := http.NewRequest("GET", "/api/test", nil)

				next := func(w http.ResponseWriter, req *http.Request) {
					Expect(req.URL.Path).To(Equal("/api/test"))
					Expect(req.Method).To(Equal("GET"))
				}

				var handler domain.MiddlewareFunc
				handler = func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
					Expect(req.URL.Path).To(Equal("/api/test"))
					Expect(req.Method).To(Equal("GET"))
					next(w, req)
				}
				handler.ServeHTTP(recorder, request, next)

			})

		})

	})

	Describe("ContextMiddlewareFunc Type", func() {
		Describe("ServeHTTP()", func() {
			It("should be working", func() {
				recorder := httptest.NewRecorder()
				request, _ := http.NewRequest("GET", "/api/test", nil)
				ctx := context.New()

				next := func(w http.ResponseWriter, req *http.Request) {
					val := ctx.Get(req, "TESTKEY").(string)

					Expect(req.URL.Path).To(Equal("/api/test"))
					Expect(req.Method).To(Equal("GET"))
					Expect(val).To(Equal("TESTVALUE"))
				}

				var handler domain.ContextMiddlewareFunc
				handler = func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc, ctx domain.IContext) {
					Expect(req.URL.Path).To(Equal("/api/test"))
					Expect(req.Method).To(Equal("GET"))

					ctx.Set(req, "TESTKEY", "TESTVALUE")

					next(w, req)
				}
				handler.ServeHTTP(recorder, request, next, ctx)

			})

		})

	})
})
