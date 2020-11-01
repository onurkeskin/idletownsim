package appengessentials

import (
	"sync"
)

type keyProvider struct {
	data map[string]string
}

var instance *keyProvider
var once sync.Once

func GetInstance() *keyProvider {
	once.Do(func() {
		instance = &keyProvider{data: make(map[string]string)}
	})
	return instance
}

var (
	//mutex sync.RWMutex
	data = make(map[string]string)
	//datat = make(map[*http.Request]int64)
)

type KeyProviderFunc func(name string, filePath string) (values []string)

func (prov keyProvider) GetKeyFor(id string) string {
	return prov.data[id]
}

func (prov keyProvider) SetKeyFor(id string, key string) {
	prov.data[id] = key
}
