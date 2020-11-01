package userhooks

import (
	"github.com/app/helpers/emailservice"
	"github.com/app/server/domain"
	"github.com/app/users"
	"net/http"
)

func PostCreateUserHookF(resource *users.Resource, w http.ResponseWriter, req *http.Request, payload *users.PostCreateUserHookPayload) error {
	// tochange: Maybe use better thing instead of casting

	var iuser domain.IUser = (payload.User)
	user := iuser.(*users.User)
	toSend := user.Email
	confirmationCode := user.ConfirmationCode

	msg := "Hello" + toSend + " your confirmation code : " + confirmationCode
	title := "Welcome"

	emailhelper.SendMail(toSend, title, msg)

	return nil
}
