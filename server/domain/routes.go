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

import "net/http"

// RouteHandlerVersion type
type RouteHandlerVersion string

// RouteHandlers is a map of route version to its handler
type RouteHandlers map[RouteHandlerVersion]http.HandlerFunc

// Route type
// Note that DefaultVersion must exists in RouteHandlers map
// See routes.go for examples
type Route struct {
	Name           string
	Method         string
	Pattern        string
	Queries        []string
	DefaultVersion RouteHandlerVersion
	RouteHandlers  RouteHandlers
	ACLHandler     ACLHandlerFunc
}

// Routes type
type Routes []Route

// Append Returns a new slice of Routes
func (r *Routes) Append(routes ...*Routes) Routes {
	res := Routes{}
	// copy current route
	for _, route := range *r {
		res = append(res, route)
	}
	for _, _routes := range routes {
		for _, route := range *_routes {
			res = append(res, route)
		}
	}
	return res
}
