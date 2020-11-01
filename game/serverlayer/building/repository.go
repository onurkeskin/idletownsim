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


package building

import (
	"errors"
	"fmt"
	buildingdbdomain "github.com/app/game/serverlayer/building/domain"
	"github.com/app/helpers/version"
	"github.com/app/server/domain"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"log"
	//"time"
)

const BuildingCollection string = "buildings"

type BuildingRepository struct {
	DB domain.IDatabase
}

// Createbuilding Insert new building document into the database
func (repo *BuildingRepository) CreateBuilding(_building buildingdbdomain.IBuildingDB) error {
	building := _building.(*BuildingDB)
	if building.ID.Hex() == "" {
		building.ID = bson.NewObjectId()
	}
	return repo.DB.Insert(BuildingCollection, building)
}

// GetUsers Get list of users
func (repo *BuildingRepository) GetBuildings() buildingdbdomain.IBuildingsDB {
	buildings := BuildingsDB{}
	err := repo.DB.FindAll(BuildingCollection, nil, &buildings, 50, "")
	if err != nil {
		return nil
	}
	return &buildings
}

// GetUser Get user specified by the id
func (repo *BuildingRepository) GetBuildingById(id string) (buildingdbdomain.IBuildingDB, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	var building BuildingDB
	err := repo.DB.FindOne(BuildingCollection, domain.Query{"_id": bson.ObjectIdHex(id)}, &building)
	return &building, err
}

func (repo *BuildingRepository) GetBuildingsByBuildingId(bid string) (buildingdbdomain.IBuildingsDB, error) {
	if !bson.IsObjectIdHex(bid) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", bid))
	}

	err := repo.DB.EnsureIndex(BuildingCollection, mgo.Index{
		Key: []string{
			"$text:bid",
		},
		Background: true,
		Sparse:     false,
	})
	q := domain.Query{"bid": bson.ObjectIdHex(bid)}
	var buildings BuildingsDB
	err = repo.DB.FindAll(BuildingCollection, q, &buildings, 999999, "")

	return &buildings, err

}

func (repo *BuildingRepository) GetBuildingByBuildingIdVersion(bid string, ver version.Version) (buildingdbdomain.IBuildingDB, error) {
	if !bson.IsObjectIdHex(bid) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", bid))
	}

	err := repo.DB.EnsureIndex(BuildingCollection, mgo.Index{
		Key: []string{
			"$text:bid",
		},
		Background: true,
		Sparse:     false,
	})
	q := domain.Query{"bid": bson.ObjectIdHex(bid)}
	q["version"] = ver
	var building BuildingDB
	err = repo.DB.FindAll(BuildingCollection, q, &building, 999999, "")

	return &building, err
}
