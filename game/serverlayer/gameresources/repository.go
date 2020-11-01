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


package resource

import (
	"errors"
	"fmt"
	dbmodels "github.com/app/game/serverlayer/dbmodels"

	"github.com/app/server/domain"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"log"
	//"time"
)

const resourceCollection string = "resources"

type ResourceRepository struct {
	DB domain.IDatabase
}

// Createbuilding Insert new building document into the database
func (repo *ResourceRepository) CreateResource(resource *dbmodels.ResourceDB) error {
	if resource.ID.Hex() == "" {
		resource.ID = bson.NewObjectId()
	}
	return repo.DB.Insert(resourceCollection, resource)
}

// GetUsers Get list of users
func (repo *ResourceRepository) GetResources() dbmodels.ResourcesDB {
	resources := dbmodels.ResourcesDB{}
	err := repo.DB.FindAll(resourceCollection, nil, &resources, 50, "")
	if err != nil {
		return nil
	}
	return resources
}

// GetUser Get user specified by the id
func (repo *ResourceRepository) GetResourceById(id string) (*dbmodels.ResourceDB, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	var resource dbmodels.ResourceDB
	err := repo.DB.FindOne(resourceCollection, domain.Query{"_id": bson.ObjectIdHex(id)}, &resource)

	return &resource, err
}
