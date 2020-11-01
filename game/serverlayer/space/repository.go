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
	"errors"
	"fmt"
	spacedbdomain "github.com/app/game/serverlayer/space/domain"
	//"github.com/app/helpers/version"
	"gopkg.in/mgo.v2"

	dbmodels "github.com/app/game/serverlayer/dbmodels"
	//Spacedbdomain "github.com/app/game/serverlayer/Space/domain"

	"github.com/app/server/domain"
	"gopkg.in/mgo.v2/bson"
	//"log"
	//"time"
)

const SpaceCollection string = "Spaces"

type SpaceRepositoryFactory struct {
}

func NewSpaceRepositoryFactory() spacedbdomain.ISpaceRepositoryFactory {
	return &SpaceRepositoryFactory{}
}
func (factory *SpaceRepositoryFactory) New(db domain.IDatabase) spacedbdomain.ISpaceRepository {
	return &SpaceRepository{
		DB: db,
	}
}

type SpaceRepository struct {
	DB domain.IDatabase
}

// CreateSpace Insert new Space document into the database
func (repo *SpaceRepository) CreateSpace(Space *dbmodels.SpaceDB) error {
	Space.ID = bson.NewObjectId()
	return repo.DB.Insert(SpaceCollection, Space)
}

func (repo *SpaceRepository) CreateSpaces(Spaces dbmodels.SpacesDB) error {
	for _, v := range Spaces {
		v.ID = bson.NewObjectId()
		err := repo.DB.Insert(SpaceCollection, v)
		//TO DO
		if err != nil {
			return err
		}
	}
	return nil
}

// GetUsers Get list of users
func (repo *SpaceRepository) GetSpacesForGMap(id string) (dbmodels.SpacesDB, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	err := repo.DB.EnsureIndex(SpaceCollection, mgo.Index{
		Key: []string{
			"$text:mapid",
		},
		Background: true,
		Sparse:     false,
	})
	q := domain.Query{"mapid": bson.ObjectIdHex(id)}
	var Spaces dbmodels.SpacesDB
	err = repo.DB.FindAll(SpaceCollection, q, &Spaces, 999999, "")

	//fmt.Println(err)
	return Spaces, err
}

// GetUser Get user specified by the id
func (repo *SpaceRepository) GetSpaceById(id string) (*dbmodels.SpaceDB, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	var Space dbmodels.SpaceDB
	err := repo.DB.FindOne(SpaceCollection, domain.Query{"_id": bson.ObjectIdHex(id)}, &Space)
	return &Space, err
}

func (repo *SpaceRepository) UpdateSpace(id string, inSpace *dbmodels.SpaceDB) (*dbmodels.SpaceDB, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	//inSpace := _inSpace.(*dbmodels.SpaceDB)

	// serialize to a sub-set of allowed User fields to update
	update := domain.Query{}
	if inSpace.BelongingMapID.Hex() != "" {
		update["mapid"] = inSpace.BelongingMapID
	}

	/*
		if inSpace.Element != nil {
			//update["element"] = map[string]interface{}{"type": (inSpace.Element).BuildableType, "buildable": (inSpace.Element).BuildableJSON}
			//update["element.type"] = (inSpace.Element).BuildableType
			//update["element.buildable"] = (inSpace.Element).BuildableJSON
			if inSpace.Element.BuildableType != "" && inSpace.Element.BuildableVersion != (version.Version{}) && inSpace.Element.Buildable != nil {
				update["occupier"] = inSpace.Element
			} else {

				if inSpace.Element.BuildableType != "" {
					update["occupier.type"] = inSpace.Element.BuildableType
				}
				if inSpace.Element.BuildableVersion != (version.Version{}) {
					update["occupier.version"] = inSpace.Element.BuildableVersion
				}
				if inSpace.Element.Buildable != nil {
					update["occupier.buildable"] = inSpace.Element.Buildable
				}
			}
		}
	*/

	if inSpace.Element != nil {
		update["occupier"] = inSpace.Element
	}

	if len(inSpace.Around) != 0 {
		ok := false
		for _, any := range inSpace.Around {
			if len(any) > 0 {
				ok = true
				break
			}
		}
		if ok {
			update["around"] = inSpace.Around
		}
	}

	query := domain.Query{"_id": bson.ObjectIdHex(id)}
	change := domain.Change{
		Update:    domain.Query{"$set": update},
		ReturnNew: true,
	}

	var changedSpace dbmodels.SpaceDB
	err := repo.DB.Update(SpaceCollection, query, change, &changedSpace)
	//fmt.Println(err)
	//fmt.Println(changedSpace)
	return &changedSpace, err
}
