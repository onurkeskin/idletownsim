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
	"golang.org/x/net/context"
	"net/http"
)

type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

func (h ContextHandlerFunc) ServeHTTP(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
	h(ctx, rw, r)
}

type MiddlewareFunc func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (m MiddlewareFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	m(rw, r, next)
}

type ContextMiddlewareFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (m ContextMiddlewareFunc) ServeHTTP(ctx context.Context, rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	m(ctx, rw, r, next)
}

type IMiddleware interface {
	Handler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type ContextMiddleware interface {
	Handler(ctx context.Context, rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}
