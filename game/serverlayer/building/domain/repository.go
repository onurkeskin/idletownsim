package building

import (
	"github.com/app/helpers/version"
	"github.com/app/server/domain"
)

type IBuildingRepositoryFactory interface {
	New(db domain.IDatabase) IBuildingRepository
}

type IBuildingRepository interface {
	CreateBuilding(buil IBuildingDB) error
	GetBuildings() IBuildingsDB
	GetBuildingById(id string) (IBuildingDB, error)
	GetBuildingsByBuildingId(bid string) (IBuildingsDB, error)
	GetBuildingByBuildingIdVersion(bid string, ver version.Version) (IBuildingDB, error)
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
