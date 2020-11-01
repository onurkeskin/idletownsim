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


package domain

import (
	"gopkg.in/mgo.v2"
)

type Query map[string]interface{}
type Change mgo.Change
type Index mgo.Index

// Database interface
type IDatabase interface {
	Insert(name string, obj interface{}) error
	Update(name string, query Query, change Change, result interface{}) error
	UpdateAll(name string, query Query, change Query) (int, error)
	FindOne(name string, query Query, result interface{}) error
	FindAll(name string, query Query, result interface{}, limit int, sort string) error
	Count(name string, query Query) (int, error)
	RemoveOne(name string, query Query) error
	RemoveAll(name string, query Query) error
	Exists(name string, query Query) bool
	DropCollection(name string) error
	DropDatabase() error
	EnsureIndex(name string, index mgo.Index) error
}
