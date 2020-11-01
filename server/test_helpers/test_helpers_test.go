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


package test_helpers_test

import (
	"net/http/httptest"

	"github.com/app/server/middlewares/context"
	"github.com/app/server/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test Helpers", func() {

	Describe("NewTestResource", func() {

		It("should return nil routes when options.NilRoutes=true", func() {
			ctx := context.New()
			testResource := test_helpers.NewTestResource(ctx, nil, &test_helpers.TestResourceOptions{
				NilRoutes: true,
			})
			Expect(testResource.Routes()).To(BeNil())
		})
		It("should return routes when options.NilRoutes=false", func() {
			ctx := context.New()
			testResource := test_helpers.NewTestResource(ctx, nil, &test_helpers.TestResourceOptions{
				NilRoutes: false,
			})
			Expect(testResource.Routes()).ToNot(BeNil())
		})
		It("should return context", func() {
			ctx := context.New()
			testResource := test_helpers.NewTestResource(ctx, nil, &test_helpers.TestResourceOptions{
				NilRoutes: false,
			})
			Expect(testResource.Context()).To(Equal(ctx))
		})
	})
	Describe("MapFromJSON", func() {

		It("should map JSON string bytes to map[] if data is a valid JSON", func() {
			data := []byte(`
				{
					"a": "isString",
					"b": 100,
					"c": true
				}
			`)
			body := test_helpers.MapFromJSON(data)
			Expect(body["a"]).To(Equal("isString"))
			Expect(body["b"]).To(Equal(float64(100)))
			Expect(body["c"]).To(Equal(true))

		})

		It("should panic if data is an invalid json ", func() {
			data := []byte("{this is an invalid json}")
			Expect(func() {
				_ = test_helpers.MapFromJSON(data)
			}).Should(Panic())
		})
	})

	Describe("DecodeResponseToType", func() {

		type TestResponseType struct {
			A string `json:"a"`
			B int    `json:"b"`
			C bool   `json:"c"`
		}

		It("should map ResponseRecorder body data to target type if data is a valid JSON", func() {
			data := []byte(`
				{
					"a": "isString",
					"b": 100,
					"c": true
				}
			`)

			var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
			recorder.Body.Write(data)

			var responseType TestResponseType
			test_helpers.DecodeResponseToType(recorder, &responseType)

			Expect(responseType).To(Equal(TestResponseType{
				A: "isString",
				B: 100,
				C: true,
			}))
		})

		It("should not panic if data is an invalid json ", func() {
			data := []byte("{this is an invalid json}")

			var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
			recorder.Body.Write(data)

			Expect(func() {
				var responseType TestResponseType
				test_helpers.DecodeResponseToType(recorder, &responseType)
			}).ShouldNot(Panic())
		})
	})
})
