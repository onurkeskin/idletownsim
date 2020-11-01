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
	"github.com/app/map/mapmodels"
	"github.com/app/server/domain"
)

type IMapRepositoryFactory interface {
	New(db domain.IDatabase) IMapRepository
}

type IMapRepository interface {
	CreateMap(m IMap) error
	UpdateMap(id string, inUser IMap) (IMap, error)
	DeleteMap(id string) error
	GetMapById(id string) (IMap, error)
	//GetMaps() []domain.IMap
	//FilterMaps(field string, query string, lastID string, limit int, sort string) domain.IMap
	FilterMaps(greater mapmodels.LatLng, smaller mapmodels.LatLng, lastID string, limit int, sort string) IMaps
	CountMaps(field string, query string) int
	DeleteAllMaps() error
}
