package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/app/helpers/emailservice"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

//---- User Request API v0 ----

// A ListUsersResponse parameter model.
//
// Used as a response for seeing users.
//
// swagger:response listUsersResponse_v0
type ListUsersResponse_v0 struct {
	// in: body
	Users   Users  `json:"users"`
	LastID  string `json:"last_id, omitempty"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

// A CreateUserRequest model.
//
// This is a CreateUserRequest_v0 request model
//
// swagger:parameters handleCreateUser_v0
type CreateUserRequest_v0 struct {
	// in: body
	User NewUser `json:"user"`
}

// A CreateUserResponse parameter model.
//
// Used as a response for creating users.
//
// swagger:response createUserResponse_v0
type CreateUserResponse_v0 struct {
	// in: body
	User    User   `json:"user,omitempty"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

// A ConfirmUserResponse parameter model.
//
// Used as a response for confirming users.
//
// swagger:response confirmUserResponse_v0
type ConfirmUserResponse_v0 struct {
	// in: body
	Code    string `json:"code,omitempty"`
	User    User   `json:"user,omitempty"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

// A UpdateUsersRequest model.
//
// This is a UpdateUsersRequest_v0 request model
//
// swagger:parameters handleUpdateUsers_v0
type UpdateUsersRequest_v0 struct {
	// in: body
	Action string   `json:"action"`
	IDs    []string `json:"ids"`
}

// A UpdateUsersResponse parameter model.
//
// Used as a response for updating users.
//
// swagger:response updateUsersResponse_v0
type UpdateUsersResponse_v0 struct {
	// in: body
	Action  string   `json:"action,omitempty"`
	IDs     []string `json:"ids,omitempty"`
	Message string   `json:"message,omitempty"`
	Success bool     `json:"success"`
}

// A DeleteAllUsersResponse parameter model.
//
// Used as a response for deleting all users.
//
// swagger:response deleteAllUsersResponse_v0
type DeleteAllUsersResponse_v0 struct {
	// in: body
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

// A GetUserResponse parameter model.
//
// Used as a response for getting user.
//
// swagger:response getUserResponse_v0
type GetUserResponse_v0 struct {
	// in: body
	User    User   `json:"user,omitempty"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

// A UpdateUserRequest model.
//
// This is a UpdateUserRequest_v0 request model
//
// swagger:parameters handleUpdateUser_v0
type UpdateUserRequest_v0 struct {
	// in: body
	User User `json:"user"`
}

// A UpdateUserResponse parameter model.
//
// Used as a response for getting user.
//
// swagger:response updateUserResponse_v0
type UpdateUserResponse_v0 struct {
	// in: body
	User    User   `json:"user,omitempty"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

// A DeleteUserResponse parameter model.
//
// Used as a response for deleting user.
//
// swagger:response deleteUserResponse_v0
type DeleteUserResponse_v0 struct {
	// in: body
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

// A CountUsersResponse parameter model.
//
// Used as a response for counting users.
//
// swagger:response countUsersResponse_v0
type CountUsersResponse_v0 struct {
	// in: body
	Count   int    `json:"count,omitempty"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

// A ErrorResponse parameter model.
//
// Used as a response for errors.
//
// swagger:response errorResponse_v0
type ErrorResponse_v0 struct {
	// in: body
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

func (resource *Resource) DecodeRequestBody(w http.ResponseWriter, req *http.Request, target interface{}) error {
	/*requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))
	*/
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(target)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, fmt.Sprintf("Request body parse error: %v", err.Error()))
		return err
	}
	return nil
}

func (resource *Resource) RenderError(w http.ResponseWriter, req *http.Request, status int, message string) {
	resource.Render(w, req, status, ErrorResponse_v0{
		Message: message,
		Success: false,
	})
}

// HandleListUsers_v0 lists users
// HandleListUsers_v0 swagger:route GET /users users handleListUsers_v0
//
// Gets the user list.
//
// Responses:
//    default: errorResponse_v0
//        200: listUsersResponse_v0
func (resource *Resource) HandleListUsers_v0(w http.ResponseWriter, req *http.Request) {
	repo := resource.UserRepository(req)

	// filter & pagination params
	field := req.FormValue("field")
	query := req.FormValue("q")
	lastID := req.FormValue("last_id")
	perPageStr := req.FormValue("per_page")
	sort := req.FormValue("sort")

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil {
		perPage = 20
	}

	u := repo.FilterUsers(field, query, lastID, perPage, sort)
	users := *u.(*Users)
	if len(users) > 0 {
		lastID = users[len(users)-1].ID.Hex()
	}
	resource.Render(w, req, http.StatusOK, ListUsersResponse_v0{
		Users:   users,
		LastID:  lastID,
		Message: "User list retrieved",
		Success: true,
	})
}

// HandleUpdateList_v0 update a list of users
func (resource *Resource) HandleUpdateUsers_v0(w http.ResponseWriter, req *http.Request) {

	var body UpdateUsersRequest_v0
	err := resource.DecodeRequestBody(w, req, &body)
	if err != nil {
		return
	}

	var message = "User list updated"
	var success bool = true
	var returnStatus = http.StatusOK

	if body.Action == "delete" {
		repo := resource.UserRepository(req)
		err = repo.DeleteUsers(body.IDs)
	} else {
		err = errors.New("Invalid action")
	}
	if err != nil {
		success = false
		message = err.Error()
		returnStatus = http.StatusBadRequest
	}

	resource.Render(w, req, returnStatus, UpdateUsersResponse_v0{
		Action:  body.Action,
		IDs:     body.IDs,
		Message: message,
		Success: success,
	})
}

// HandleDeleteAll_v0 deletes all users
func (resource *Resource) HandleDeleteAllUsers_v0(w http.ResponseWriter, req *http.Request) {
	repo := resource.UserRepository(req)
	_ = repo.DeleteAllUsers()

	resource.Render(w, req, http.StatusOK, DeleteAllUsersResponse_v0{
		Message: "All users deleted",
		Success: true,
	})
}

// HandleCreateUser_v0 creates a new user
// HandleListUsers_v0 swagger:route POST /users users handleCreateUser_v0
//
// Gets the user list.
//
// Responses:
//    default: errorResponse_v0
//        200: createUserResponse_v0
func (resource *Resource) HandleCreateUser_v0(w http.ResponseWriter, req *http.Request) {
	repo := resource.UserRepository(req)
	var body CreateUserRequest_v0
	err := resource.DecodeRequestBody(w, req, &body)
	if err != nil {
		return
	}

	if repo.UserExistsByUsername(body.User.Username) {
		resource.RenderError(w, req, http.StatusBadRequest, "Username already exists")
		return
	}

	if repo.UserExistsByEmail(body.User.Email) {
		resource.RenderError(w, req, http.StatusBadRequest, "User with email address already exists")
		return
	}

	// New user always have no roles assigned until confirmed
	// Set flag to `pending` awaiting user to confirm email
	var newUser = User{
		Username: body.User.Username,
		Email:    body.User.Email,
		Roles:    Roles{},
		Status:   StatusPending,
	}
	//fmt.Println("User Requested acc creation" + newUser.Email + " " + newUser.Username + " " + newUser.HashedPassword)

	// generate new code
	newUser.GenerateConfirmationCode()

	// set password (hashed)
	newUser.SetPassword(body.User.Password)

	// ensure that user obj is valid
	if !newUser.IsValid() {
		resource.RenderError(w, req, http.StatusBadRequest, "Invalid user object")
		return
	}

	err = repo.CreateUser(&newUser)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, "Failed to save user object")
		return
	}

	// run a post-create hook
	// example of a post-create hook: send email / message with confirmation link
	if resource.ControllerHooks.PostCreateUserHook != nil {
		err = resource.ControllerHooks.PostCreateUserHook(resource, w, req, &PostCreateUserHookPayload{
			User: &newUser,
		})
		if err != nil {
			resource.RenderError(w, req, http.StatusBadRequest, err.Error())
			return
		}
	}

	resource.Render(w, req, http.StatusCreated, CreateUserResponse_v0{
		User:    newUser,
		Message: "User created",
		Success: true,
	})
}

// HandleConfirmEmail_v0 confirms user's email address
func (resource *Resource) HandleConfirmUser_v0(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]
	code := req.FormValue("code")

	repo := resource.UserRepository(req)
	_user, err := repo.GetUserById(id)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}

	user := _user.(*User)
	if user.Status != StatusPending {
		resource.RenderError(w, req, http.StatusBadRequest, "User not pending confirmation")
		return
	}

	if !user.IsCodeVerified(code) {
		resource.RenderError(w, req, http.StatusBadRequest, "Invalid code")
		return
	}

	// set user status to `active`
	_updatedUser, err := repo.UpdateUser(id, &User{
		Status: StatusActive,
		Roles:  Roles{RoleUser},
	})
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}
	updatedUser := _updatedUser.(*User)

	// run a post-confirmation hook
	if resource.ControllerHooks.PostConfirmUserHook != nil {
		err = resource.ControllerHooks.PostConfirmUserHook(resource, w, req, &PostConfirmUserHookPayload{
			User: user,
		})
		if err != nil {
			resource.RenderError(w, req, http.StatusBadRequest, err.Error())
			return
		}
	}

	resource.Render(w, req, http.StatusOK, ConfirmUserResponse_v0{
		Code:    code,
		User:    *updatedUser,
		Message: "User confirmed",
		Success: true,
	})
}

func (resource *Resource) HandleResendConfirmation_v0(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]

	repo := resource.UserRepository(req)
	_user, err := repo.GetUserById(id)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}

	user := _user.(*User)
	if user.Status != StatusPending {
		resource.RenderError(w, req, http.StatusBadRequest, "User has already confirmed their email")
		return
	}

	toSend := user.Email
	confirmationCode := user.ConfirmationCode

	msg := "Hello" + toSend + " your confirmation code : " + confirmationCode
	title := "Welcome"

	emailhelper.SendMail(toSend, title, msg)

	resource.Render(w, req, http.StatusOK, ConfirmUserResponse_v0{
		Message: "Mail Sent",
		Success: true,
	})

}

// HandleGetUser_v0 gets user object
func (resource *Resource) HandleGetUser_v0(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]

	repo := resource.UserRepository(req)
	_user, err := repo.GetUserById(id)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, "User not found")
		return
	}
	user := _user.(*User)

	resource.Render(w, req, http.StatusOK, GetUserResponse_v0{
		User:    *user,
		Message: "User retrieved",
		Success: true,
	})
}

// HandleUpdateUser_v0 updates user object
func (resource *Resource) HandleUpdateUser_v0(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]

	var body UpdateUserRequest_v0
	err := resource.DecodeRequestBody(w, req, &body)
	if err != nil {
		return
	}

	repo := resource.UserRepository(req)
	_user, err := repo.UpdateUser(id, &body.User)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}
	user := _user.(*User)

	resource.Render(w, req, http.StatusOK, UpdateUserResponse_v0{
		User:    *user,
		Message: "User updated",
		Success: true,
	})
}

// HandleDelete_v0 deletes user object
func (resource *Resource) HandleDeleteUser_v0(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]

	repo := resource.UserRepository(req)

	err := repo.DeleteUser(id)
	if err != nil {
		resource.RenderError(w, req, http.StatusBadRequest, err.Error())
		return
	}

	resource.Render(w, req, http.StatusOK, DeleteUserResponse_v0{
		Message: "User deleted",
		Success: true,
	})
}

func (resource *Resource) HandleCountUsers_v0(w http.ResponseWriter, req *http.Request) {

	// filter & pagination params
	field := req.FormValue("field")
	query := req.FormValue("q")

	repo := resource.UserRepository(req)
	count := repo.CountUsers(field, query)

	resource.Render(w, req, http.StatusOK, CountUsersResponse_v0{
		Count:   count,
		Message: "Users count retrieved",
		Success: true,
	})
}
