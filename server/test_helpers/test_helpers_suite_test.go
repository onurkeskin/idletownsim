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


package test_helpers_test

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"testing"
)

var TestDatabaseServerName = "localhost"
var TestDatabaseName = "test_db"

var privateSigningKey []byte
var publicSigningKey []byte

func TestTestHelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Helpers Suite")
}

var _ = BeforeSuite(func() {
	// each test node uses its own database
	TestDatabaseName = fmt.Sprintf("test_db_node%v", config.GinkgoConfig.ParallelNode)

	var err error
	privateSigningKey, err = ioutil.ReadFile("../keys/demo.rsa")
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error loading private signing key: %v", err.Error())))
	}
	publicSigningKey, err = ioutil.ReadFile("../keys/demo.rsa.pub")
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error loading public signing key: %v", err.Error())))
	}
})
