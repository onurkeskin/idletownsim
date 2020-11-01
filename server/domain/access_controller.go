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


package domain

import (
	"net/http"
)

type ACLHandlerFunc func(*http.Request, IUser) (bool, string)

type ACLMap map[string]ACLHandlerFunc

func (m *ACLMap) Append(maps ...*ACLMap) ACLMap {
	res := ACLMap{}
	// copy current map
	for k, v := range *m {
		res[k] = v
	}
	for _, _maps := range maps {
		for k, v := range *_maps {
			res[k] = v
		}
	}
	return res
}

type IAccessController interface {
	Add(*ACLMap)
	AddHandler(name string, handler ACLHandlerFunc)
	HasAction(string) bool
	IsHTTPRequestAuthorized(req *http.Request, action string, user IUser) (bool, string)
	NewContextHandler(string, http.HandlerFunc) http.HandlerFunc
	//	Render(w http.ResponseWriter, req *http.Request, status int, v interface{})
}
