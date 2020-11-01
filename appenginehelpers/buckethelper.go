package appengessentials

import (
	"cloud.google.com/go/storage"
	"fmt"
	"github.com/pkg/errors"
	appctx "golang.org/x/net/context"
	"io/ioutil"
	"log"
)

var (
	inited        = false
	storageClient *storage.Client
	appcontext    appctx.Context
	// Set this in app.yaml when running in production.
	bucket = "appstoragebucket"
)

func checkAndInit() {
	if !inited {
		toInit()
		inited = true
	}
}

func toInit() {
	appcontext = appctx.Background()
	var err error
	storageClient, err = storage.NewClient(appcontext)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error connecting: %v", err.Error())))
	}
	testBucket()
}

func GetObject(name string) *storage.ObjectHandle {
	checkAndInit()
	return storageClient.Bucket(bucket).Object(name)
}

func GetReaderOnObject(name string) *storage.Reader {
	checkAndInit()
	reader, _ := storageClient.Bucket(bucket).Object(name).NewReader(appcontext)
	return reader
}

func ReadFileOnBucket(name string) []byte {
	checkAndInit()
	handle := storageClient.Bucket(bucket).Object(name)
	if handle == nil {
		panic(errors.New(fmt.Sprintf("Bucket handle nil")))
	}
	reader, err := handle.NewReader(appcontext)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Cant Read file %v :%v", name, err.Error())))
	}
	file, err := ioutil.ReadAll(reader)
	reader.Close()
	if err != nil {
		panic(errors.New(fmt.Sprintf("Cant Read file %v :%v", name, err.Error())))
	}
	return file
}

func testBucket() {
	handle := storageClient.Bucket(bucket)
	if handle == nil {
		panic(errors.New(fmt.Sprintf("Null Bucket")))
	} else {
		log.Println("Connected To Bucket")
	}
}
