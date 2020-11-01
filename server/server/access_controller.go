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
	"github.com/app/users"
	"github.com/app/server/domain"
	"net/http"
)

const defaultForbiddenAccessMessage = "Forbidden (403)"
const defaultOKAccessMessage = "OK"

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

// TODO: Currently, AccessController only acts as a gateway for endpoints on router level. Build AC to handler other aspects of ACL
func NewAccessController(renderer domain.IRenderer) *AccessController {
	return &AccessController{domain.ACLMap{}, renderer}
}

// implements IAccessController
type AccessController struct {
	ACLMap   domain.ACLMap
	renderer domain.IRenderer
}

func (ac *AccessController) Add(aclMap *domain.ACLMap) {
	ac.ACLMap = ac.ACLMap.Append(aclMap)
}

func (ac *AccessController) AddHandler(action string, handler domain.ACLHandlerFunc) {
	ac.ACLMap[action] = handler
}

func (ac *AccessController) HasAction(action string) bool {
	fn := ac.ACLMap[action]
	return (fn != nil)
}

func (ac *AccessController) IsHTTPRequestAuthorized(req *http.Request, action string, user domain.IUser) (bool, string) {
	fn := ac.ACLMap[action]
	if fn == nil {
		// by default, if acl action/handler is not defined, request is not authorized
		return false, defaultForbiddenAccessMessage
	}

	result, message := fn(req, user)
	if result && message == "" {
		message = defaultOKAccessMessage
	}
	if !result && message == "" {
		message = defaultForbiddenAccessMessage
	}
	return result, message
}

func (ac *AccessController) NewContextHandler(action string, next http.HandlerFunc) http.HandlerFunc {
	//func (ac *AccessController) NewHandler(action string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		//user := ac.ctx.GetCurrentUserCtx(req)
		user := users.GetUserCtx(req.Context())
		// `user` might be `nil` if has not authenticated.
		// ACL might want to allow anonymous / non-authenticated access (for login, e.g)

		//fmt.Println("At Context handler creation", user)
		result, message := ac.IsHTTPRequestAuthorized(req, action, user)
		if !result {
			ac.renderer.Render(w, req, http.StatusForbidden, ErrorResponse{
				Message: message,
				Success: false,
			})
			return
		}

		next(w, req)
	}
}
