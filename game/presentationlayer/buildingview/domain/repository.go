package buildingview

import (
	"github.com/app/server/domain"
)

type IBuildingViewRepositoryFactory interface {
	New(db domain.IDatabase) IBuildingViewRepository
}

type IBuildingViewRepository interface {
	CreateBuildingView(buil IBuildingView) error
	GetBuildingsView() IBuildingsView
	GetBuildingViewById(id string) (IBuildingView, error)
	/*
		CreateUser(user domain.IUser) error
		GetUsers() domain.IUsers
		FilterUsers(field string, query string, lastID string, limit int, sort string) domain.IUsers
		CountUsers(field string, query string) int
		DeleteUsers(ids []string) error
		DeleteAllUsers() error
		GetUserById(id string) (domain.IUser, error)
		GetUserByUsername(username string) (domain.IUser, error)
		UserExistsByUsername(username string) bool
		UserExistsByEmail(email string) bool
		UpdateUser(id string, inUser domain.IUser) (domain.IUser, error)
		DeleteUser(id string) error
	*/
}
