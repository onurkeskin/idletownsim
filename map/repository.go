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


package maps

import (
	"errors"
	"fmt"
	mapdomain "github.com/app/map/domain"
	"github.com/app/map/mapmodels"
	"github.com/app/server/domain"
	"gopkg.in/mgo.v2"
	"log"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"log"
	"time"
)

// User collection name
const MapsCollection string = "maps"

type MapRepository struct {
	DB domain.IDatabase
}

// CreateUser Insert new user document into the database
func (repo *MapRepository) CreateMap(_map mapdomain.IMap) error {
	cmap := _map.(*Map)
	cmap.ID = bson.NewObjectId()
	cmap.CreatedDate = time.Now()
	cmap.LastModifiedDate = time.Now()
	return repo.DB.Insert(MapsCollection, cmap)
}

func (repo *MapRepository) CountMaps(field string, query string) int {
	q := domain.Query{}
	if query != "" {
		if field != "" {
			q[field] = domain.Query{
				"$regex":   fmt.Sprintf("^%v.*", query),
				"$options": "i",
			}
		} else {
			// if not field is specified, we do a text search on pre-defined text index
			q["$text"] = domain.Query{
				"$search": query,
			}
		}
	}

	count, err := repo.DB.Count(MapsCollection, q)
	if err != nil {
		return 0
	}
	return count
}

func (repo *MapRepository) GetMapById(id string) (mapdomain.IMap, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	var _map Map
	err := repo.DB.FindOne(MapsCollection, domain.Query{"_id": bson.ObjectIdHex(id)}, &_map)
	return &_map, err
}

func (repo *MapRepository) UpdateMap(id string, _inMap mapdomain.IMap) (mapdomain.IMap, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}
	inMap := _inMap.(*Map)

	update := domain.Query{
		"lastModifiedDate": time.Now(),
	}
	if inMap.MapRaw != nil && len(inMap.MapRaw) != 0 {
		update["rawmap"] = inMap.MapRaw
	}
	if inMap.MapFundamentalCoordinates != (mapmodels.LatLng{}) {
		update["fundamentalcoordinates"] = inMap.MapFundamentalCoordinates
	}
	if inMap.MapCompleteAddress != "" {
		update["mapadress"] = inMap.MapCompleteAddress
	}
	if inMap.MapIdentifierAdress != "" {
		update["mapadress"] = inMap.MapIdentifierAdress
	}
	if inMap.ParsedRelations != nil && len(inMap.ParsedRelations) > 0 {
		update["parsedrelations"] = inMap.ParsedRelations
	}
	//TODO DO SMT ABOUT STATUS
	/*
		if inMap.Status != nil && len(inMap.ParsedRelations) > 0 {
			update["parsedrelations"] = inMap.ParsedRelations
		}
	*/

	query := domain.Query{"_id": bson.ObjectIdHex(id)}
	change := domain.Change{
		Update:    domain.Query{"$set": update},
		ReturnNew: true,
	}
	var changedMap Map
	err := repo.DB.Update(MapsCollection, query, change, &changedMap)
	return &changedMap, err
}

func (repo *MapRepository) DeleteMap(id string) error {
	if !bson.IsObjectIdHex(id) {
		return errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}
	err := repo.DB.RemoveOne(MapsCollection, domain.Query{"_id": bson.ObjectIdHex(id)})
	return err
}

func (repo *MapRepository) DeleteAllMaps() error {
	//TODO
	return nil
}

func (repo *MapRepository) FilterMaps(greater mapmodels.LatLng, smaller mapmodels.LatLng, lastID string, limit int, sort string) mapdomain.IMaps {
	toRet := Maps{}

	// ensure that collection has the right text index
	// refactor building collection index
	err := repo.DB.EnsureIndex(MapsCollection, mgo.Index{
		Key: []string{
			"$text:mapversion",
			"$text:lastModifiedDate",
			"$text:status",
		},
		Background: true,
		Sparse:     true,
	})
	if err != nil {
		log.Println("FilterMaps: EnsureIndex", err.Error())
	}
	// parse sort string
	allowedSortMap := map[string]bool{
		"_id":  true,
		"-_id": true,
	}
	// ensure that sort string is allowed
	// we are basically concerned about sorting on un-indexed keys
	if !allowedSortMap[sort] {
		sort = "-_id" // set it to default sort
	}

	q := domain.Query{}
	if lastID != "" && bson.IsObjectIdHex(lastID) {
		if sort == "_id" {
			q["_id"] = domain.Query{
				"$gt": bson.ObjectIdHex(lastID),
			}
		} else {
			q["_id"] = domain.Query{
				"$lt": bson.ObjectIdHex(lastID),
			}
		}
	}

	if greater != (mapmodels.LatLng{}) && smaller != (mapmodels.LatLng{}) {
		q["fundamentalcoordinates.lat"] = domain.Query{
			"$gte": greater.Lat,
			"$lte": smaller.Lat,
		}
		q["fundamentalcoordinates.lng"] = domain.Query{
			"$gte": greater.Lng,
			"$lte": smaller.Lng,
		}
	} else {

		if greater != (mapmodels.LatLng{}) {
			q["fundamentalcoordinates.lat"] = domain.Query{
				"$gte": greater.Lat,
			}
			q["fundamentalcoordinates.lng"] = domain.Query{
				"$gte": greater.Lng,
			}
		}

		if smaller != (mapmodels.LatLng{}) {
			q["fundamentalcoordinates.lat"] = domain.Query{
				"$lte": smaller.Lat,
			}
			q["fundamentalcoordinates.lng"] = domain.Query{
				"$lte": smaller.Lng,
			}
		}
	}

	//fmt.Println(q)
	err = repo.DB.FindAll(MapsCollection, q, &toRet, limit, sort)
	if err != nil {
		return &Maps{}
	}
	return &toRet
}
