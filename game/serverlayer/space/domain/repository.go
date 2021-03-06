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


package space

import (
	dbmodels "github.com/app/game/serverlayer/dbmodels"
	"github.com/app/server/domain"
)

type ISpaceRepositoryFactory interface {
	New(db domain.IDatabase) ISpaceRepository
}

type ISpaceRepository interface {
	CreateSpace(_Space *dbmodels.SpaceDB) error
	CreateSpaces(Spaces dbmodels.SpacesDB) error
	GetSpacesForGMap(id string) (dbmodels.SpacesDB, error)
	GetSpaceById(id string) (*dbmodels.SpaceDB, error)
	UpdateSpace(id string, inSpace *dbmodels.SpaceDB) (*dbmodels.SpaceDB, error)
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
