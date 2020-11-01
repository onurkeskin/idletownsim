package appengessentials

import (
	"io/ioutil"
	"strings"
)

func GetKeyFromFile(name string, filePath string) (values []string) {
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil
	}
	arr := []string{name, strings.TrimSpace(string(dat))}
	return arr
}
