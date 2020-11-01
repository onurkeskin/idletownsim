package swaggerpass

import (
	"github.com/app/server/domain"
	"net/http"
)

func (swagger *Swagger) HandleSwaggerACL(req *http.Request, user domain.IUser) (bool, string) {
	return true, ""
}
