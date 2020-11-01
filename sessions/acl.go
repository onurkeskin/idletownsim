package sessions

import (
	"github.com/app/server/domain"
	"net/http"
)

func (resource *Resource) HandleGetSessionACL(req *http.Request, user domain.IUser) (bool, string) {
	if user == nil {
		return false, ""
	}
	return true, ""
}

func (resource *Resource) HandleCreateSessionACL(req *http.Request, user domain.IUser) (bool, string) {
	// allow anonymous access
	return true, ""
}

func (resource *Resource) HandleDeleteSessionACL(req *http.Request, user domain.IUser) (bool, string) {
	// allow anonymous access
	return true, ""
}
