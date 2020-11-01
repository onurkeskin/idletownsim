package maps

import (
	"fmt"
	"github.com/app/server/domain"
	//"github.com/gorilla/mux"
	"net/http"
)

func (resource *Resource) HandleGetMapACL(req *http.Request, user domain.IUser) (bool, string) {
	if user == nil {
		return false, ""
	}
	return true, ""
}

func (resource *Resource) HandleGetCloseMapsACL(req *http.Request, user domain.IUser) (bool, string) {
	fmt.Println(user)
	if user == nil {
		return false, ""
	}
	return true, ""
}

func (resource *Resource) HandleCreateMapLatLngACL(req *http.Request, user domain.IUser) (bool, string) {
	if user == nil {
		return false, ""
	}
	return true, ""
}
