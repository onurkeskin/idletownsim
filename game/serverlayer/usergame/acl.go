package game

import (
	"github.com/app/server/domain"
	"net/http"
)

func (resource *Resource) HandleGetGamesACL(req *http.Request, user domain.IUser) (bool, string) {
	if user == nil {
		return false, ""
	}
	return true, ""
}

func (resource *Resource) HandleGetGameACL(req *http.Request, user domain.IUser) (bool, string) {
	if user == nil {
		return false, ""
	}
	return true, ""
}

func (resource *Resource) HandleCreateGameACL(req *http.Request, user domain.IUser) (bool, string) {
	if user == nil {
		return false, ""
	}
	return true, ""
}
func (resource *Resource) HandleGetEligibleBuildingsACL(req *http.Request, user domain.IUser) (bool, string) {
	if user == nil {
		return false, ""
	}
	return true, ""
}

func (resource *Resource) HandleBuyBuildingACL(req *http.Request, user domain.IUser) (bool, string) {
	if user == nil {
		return false, ""
	}
	return true, ""
}

func (resource *Resource) DeleteGameSessionsACL(req *http.Request, user domain.IUser) (bool, string) {
	if user == nil {
		return false, ""
	}
	return true, ""
}

func (resource *Resource) HandleBuyUpgradeACL(req *http.Request, user domain.IUser) (bool, string) {
	if user == nil {
		return false, ""
	}
	return true, ""
}
func (resource *Resource) HandleGetEligibleUpgradesACL(req *http.Request, user domain.IUser) (bool, string) {
	if user == nil {
		return false, ""
	}
	return true, ""
}
