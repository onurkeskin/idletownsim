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
	"fmt"
	"github.com/app/server/domain"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AcceptHeaders Tests", func() {
	var _ = Describe("NewAcceptHeadersFromString", func() {

		type testMap struct {
			TestValue       string
			ExpectedLen     int
			ExpectedResults domain.AcceptHeaders
		}
		type testMaps []testMap

		var tests = testMaps{
			testMap{
				TestValue: "",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{}, 1},
				},
			},
			testMap{
				TestValue: ";",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{}, 1},
				},
			},
			testMap{
				TestValue: ";q=",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{
						Parameters: domain.MediaTypeParams{"q": ""},
					}, 1},
				},
			},
			testMap{
				TestValue: "application/json",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{
						String:  "application/json",
						Type:    "application",
						Tree:    "",
						SubType: "json",
						Suffix:  "",
					}, 1},
				},
			},
			testMap{
				TestValue: "application/json;q=",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{
						String:     "application/json",
						Type:       "application",
						Tree:       "",
						SubType:    "json",
						Suffix:     "",
						Parameters: domain.MediaTypeParams{"q": ""},
					}, 1},
				},
			},
			testMap{
				TestValue: "application/json;q",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{
						String:     "application/json",
						Type:       "application",
						Tree:       "",
						SubType:    "json",
						Suffix:     "",
						Parameters: domain.MediaTypeParams{"q": ""},
					}, 1},
				},
			},
			testMap{
				TestValue: "application/json;  q=0.9 ",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{
						String:     "application/json",
						Type:       "application",
						Tree:       "",
						SubType:    "json",
						Suffix:     "",
						Parameters: domain.MediaTypeParams{"q": "0.9"},
					}, 0.9},
				},
			},
			testMap{
				TestValue: "application/vnd.sgk.v1+json ",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{
						String:  "application/vnd.sgk.v1+json",
						Type:    "application",
						Tree:    "vnd",
						SubType: "sgk.v1",
						Suffix:  "json",
					}, 1},
				},
			},
			testMap{
				TestValue: "application/vnd.sgk.v1+json;q=0.8",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{
						String:  "application/vnd.sgk.v1+json",
						Type:    "application",
						Tree:    "vnd",
						SubType: "sgk.v1",
						Suffix:  "json",
						Parameters: domain.MediaTypeParams{
							"q": "0.8",
						},
					}, 0.8},
				},
			},
			testMap{
				TestValue: "application/vnd.sgk+json; q=0.8 ;version=1.0",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{
						String:  "application/vnd.sgk+json",
						Type:    "application",
						Tree:    "vnd",
						SubType: "sgk",
						Suffix:  "json",
						Parameters: domain.MediaTypeParams{
							"q":       "0.8",
							"version": "1.0",
						},
					}, 0.8},
				},
			},
			testMap{
				TestValue: "application/vnd.sgk.rest-api-server.v1+json; q=0.8 ;version=1.0,*/*",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{
						String:  "application/vnd.sgk.rest-api-server.v1+json",
						Type:    "application",
						Tree:    "vnd",
						SubType: "sgk.rest-api-server.v1",
						Suffix:  "json",
						Parameters: domain.MediaTypeParams{
							"q":       "0.8",
							"version": "1.0",
						},
					}, 0.8},
					domain.AcceptHeader{domain.MediaType{
						String:  "*/*",
						Type:    "*",
						Tree:    "",
						SubType: "*",
						Suffix:  "",
					}, 1},
				},
			},
			testMap{
				TestValue: "application/vnd.sgk+json; q=0.8 ;version=1.0,application/json , */*;q=noninteger",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{
						String:  "application/vnd.sgk+json",
						Type:    "application",
						Tree:    "vnd",
						SubType: "sgk",
						Suffix:  "json",
						Parameters: domain.MediaTypeParams{
							"q":       "0.8",
							"version": "1.0",
						},
					}, 0.8},
					domain.AcceptHeader{domain.MediaType{
						String:  "application/json",
						Type:    "application",
						Tree:    "",
						SubType: "json",
						Suffix:  "",
					}, 1},
					domain.AcceptHeader{domain.MediaType{
						String:  "*/*",
						Type:    "*",
						Tree:    "",
						SubType: "*",
						Suffix:  "",
						Parameters: domain.MediaTypeParams{
							"q": "noninteger",
						},
					}, 1},
				},
			},
		}

		for _, test := range tests {
			testValue := test.TestValue
			expectedResults := test.ExpectedResults
			Context(fmt.Sprintf("when `Accept=%v`", testValue), func() {
				result := domain.NewAcceptHeadersFromString(testValue)
				It("parses OK", func() {
					Expect(len(result)).To(Equal(len(expectedResults)))
					for i, _ := range expectedResults {
						Expect(result[i].MediaType).To(Equal(expectedResults[i].MediaType))
						Expect(result[i].QualityFactor).To(Equal(expectedResults[i].QualityFactor))
					}
				})
			})
		}

	})
})
