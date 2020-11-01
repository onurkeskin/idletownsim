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


package context

import (
	"github.com/app/server/domain"
	"github.com/gorilla/context"
	"net/http"
)

const CurrentUserKey domain.ContextKey = "slumber-mddlwr-context-current-user-key"
const DatabaseKey domain.ContextKey = "slumber-mddlwr-context-database-key"

func New() *Context {
	return &Context{}
}

// Context implements IContext
type Context struct {
}

func (ctx *Context) InjectMiddleware(middleware domain.ContextMiddlewareFunc) domain.MiddlewareFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		middleware(rw, r, next, ctx)
	}
}

func (ctx *Context) Inject(handler domain.ContextHandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		handler(rw, r, ctx)
	}
}

func (ctx *Context) Set(r *http.Request, key interface{}, val interface{}) {
	context.Set(r, key, val)
}

func (ctx *Context) Get(r *http.Request, key interface{}) interface{} {
	return context.Get(r, key)
}

func (ctx *Context) SetCurrentUserCtx(r *http.Request, user domain.IUser) {
	ctx.Set(r, CurrentUserKey, user)
}

func (ctx *Context) GetCurrentUserCtx(r *http.Request) domain.IUser {
	if user := ctx.Get(r, CurrentUserKey); user != nil {
		return user.(domain.IUser)
	}
	return nil
}
