package users

import (
	"golang.org/x/net/context"

	"github.com/app/server/domain"
	"net/http"
)

const UserInfoKey domain.ContextKey = "UserInfoKey"

func GetUserCtx(ctx context.Context) domain.IUser {
	if _user := ctx.Value(UserInfoKey); _user != nil {
		return _user.(domain.IUser)
	}
	return nil
}
func NewContextWithUser(ctx context.Context, req *http.Request, user domain.IUser) context.Context {
	return context.WithValue(ctx, UserInfoKey, user)
}
