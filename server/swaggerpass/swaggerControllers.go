package swaggerpass

import (
	"fmt"
	"github.com/app/appenginehelpers"
	"log"
	"net/http"
	"strconv"
)

// swagger pass
func (swag *Swagger) HandleSwagger_v0(w http.ResponseWriter, req *http.Request) {
	file := appengessentials.ReadFileOnBucket("swagger.json")

	f := fmt.Sprintf("%v", len(file))

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, api_key, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT")
	w.Header().Set("Connection", "close")
	w.Header().Set("Content-Length", f)

	w.Header().Del("Transfer-Encoding")

	_, err := w.Write(file)
	if err != nil {
		log.Fatal(err)
	}

	return
}

func IntToString(input_num int64) string {
	// to convert a float number to a string
	return strconv.FormatInt(input_num, 64)
}
