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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
)

var _ = Describe("AccessController Tests", func() {
	Describe("ACLMap Struct", func() {
		Describe("Append()", func() {
			stub := func(req *http.Request, user domain.IUser) (bool, string) {
				return true, ""
			}
			firstMap := domain.ACLMap{
				"first": stub,
			}
			secondMap := domain.ACLMap{
				"second": stub,
			}
			thirdMap := domain.ACLMap{
				"third": stub,
			}
			var result domain.ACLMap
			var result2 domain.ACLMap
			BeforeEach(func() {
				result = firstMap.Append(&secondMap)
				result2 = firstMap.Append(&secondMap, &thirdMap)
			})
			It("should return a new map", func() {
				Expect(result["first"]).ToNot(BeNil())
				Expect(result2["first"]).ToNot(BeNil())
			})
			It("should return a new map", func() {
				Expect(result["second"]).ToNot(BeNil())
				Expect(result2["second"]).ToNot(BeNil())
			})
			It("should return a new map", func() {
				Expect(result["third"]).To(BeNil())
				Expect(result2["second"]).ToNot(BeNil())
			})
		})
	})
})
