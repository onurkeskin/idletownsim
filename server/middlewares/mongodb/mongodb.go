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


package mongodb

import (
	"github.com/app/server/domain"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"net/http"
	"time"
)

const MongoDbKey domain.ContextKey = "mongokey"

type Options struct {
	ServerName   string
	DatabaseName string
	DialTimeout  time.Duration
}

func New(options *Options) *MongoDB {
	db := &MongoDB{}
	db.options = options
	return db
}

// MongoDatabase implements IDatabase
type MongoDB struct {
	currentDb *mgo.Database
	options   *Options
}

func (db *MongoDB) NewSession() *MongoDBSession {

	mongoOptions := db.options

	// set default DialTimeout value
	if mongoOptions.DialTimeout <= 0 {
		mongoOptions.DialTimeout = 1 * time.Minute
	}

	session, err := mgo.DialWithTimeout(mongoOptions.ServerName, mongoOptions.DialTimeout)
	if err != nil {
		panic(err)
	}
	db.currentDb = session.DB(mongoOptions.DatabaseName)
	return &MongoDBSession{session, mongoOptions}
}

func (db *MongoDB) FindOne(name string, query domain.Query, result interface{}) error {
	return db.currentDb.C(name).Find(query).One(result)
}

func (db *MongoDB) FindAll(name string, query domain.Query, result interface{}, limit int, sort string) error {
	if sort == "" {
		sort = "-_id"
	}
	return db.currentDb.C(name).Find(query).Sort(sort).Limit(limit).All(result)
}

func (db *MongoDB) Count(name string, query domain.Query) (int, error) {
	return db.currentDb.C(name).Find(query).Count()
}

func (db *MongoDB) Insert(name string, obj interface{}) error {
	return db.currentDb.C(name).Insert(obj)
}

func (db *MongoDB) Update(name string, query domain.Query, change domain.Change, result interface{}) error {
	_, err := db.currentDb.C(name).Find(query).Apply(mgo.Change(change), result)
	return err
}

func (db *MongoDB) UpdateAll(name string, query domain.Query, change domain.Query) (int, error) {
	changeInfo, err := db.currentDb.C(name).UpdateAll(query, change)
	if changeInfo == nil {
		return 0, err
	}
	return changeInfo.Updated, err
}

func (db *MongoDB) RemoveOne(name string, query domain.Query) error {
	return db.currentDb.C(name).Remove(query)
}

func (db *MongoDB) RemoveAll(name string, query domain.Query) error {
	_, err := db.currentDb.C(name).RemoveAll(query)
	return err
}

func (db *MongoDB) DropCollection(name string) error {
	return db.currentDb.C(name).DropCollection()
}

func (db *MongoDB) Exists(name string, query domain.Query) bool {
	var result interface{}
	err := db.currentDb.C(name).Find(query).One(result)
	return (err == nil)
}
func (db *MongoDB) DropDatabase() error {
	return db.currentDb.DropDatabase()
}

func (db *MongoDB) EnsureIndex(name string, index mgo.Index) error {
	return db.currentDb.C(name).EnsureIndex(index)
}

// MongoDatabaseSession struct implements IContextMiddleware
type MongoDBSession struct {
	*mgo.Session
	*Options
}

// Handler Returns a middleware HandlerFunc that creates and saves a database session into request context
func (session *MongoDBSession) Handler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	s := session.Clone()
	defer s.Close()
	db := &MongoDB{
		currentDb: s.DB(session.DatabaseName),
	}
	ctx := SetMongoDbCtx(req.Context(), req, db)
	next(w, req.WithContext(ctx))
}

func SetMongoDbCtx(ctx context.Context, r *http.Request, db *MongoDB) context.Context {
	return context.WithValue(ctx, MongoDbKey, db)
}

func GetMongoDbCtx(ctx context.Context) *MongoDB {
	return ctx.Value(MongoDbKey).(*MongoDB)
}
