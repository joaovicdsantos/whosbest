package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/joaovicdsantos/whosbest-api/app/helpers"
	"github.com/joaovicdsantos/whosbest-api/app/models"
	"github.com/joaovicdsantos/whosbest-api/app/services"
	"golang.org/x/crypto/bcrypt"
)

var (
	userGetAllRe   = regexp.MustCompile(`^\/user[\/]*$`)
	userGetOneRe   = regexp.MustCompile(`^\/user\/(\d+)$`)
	userRegisterRe = regexp.MustCompile(`^\/user[\/]*$`)
)

type UserRoutes struct {
	DB          *sql.DB
	userService *services.UserService
	Payload     map[string]interface{}
}

func (u *UserRoutes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u.userService = &services.UserService{
		DB: u.DB,
	}
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost && userRegisterRe.MatchString(r.URL.Path) {
		u.Register(w, r)
		return
	}

	unparsedToken, ok := r.Header["Authorization"]
	if !ok {
		response := helpers.NewResponseError("Token is not valid", http.StatusUnauthorized)
		response.SendResponse(w)
		return
	}

	var err error
	u.Payload, err = helpers.ParseJwtToken(fmt.Sprint(unparsedToken[0]))
	if err != nil {
		response := helpers.NewResponseError(err.Error(), http.StatusUnauthorized)
		response.SendResponse(w)
		return
	}

	user := helpers.GetCurrentUser(u.Payload, u.DB)
	if user.Id == 0 {
		response := helpers.NewResponseError("Invalid user", http.StatusUnauthorized)
		response.SendResponse(w)
		return

	}
	// Authorized routes
	switch {
	case r.Method == http.MethodGet && userGetAllRe.MatchString(r.URL.Path):
		u.GetAll(w, r)
		break
	case r.Method == http.MethodGet && userGetOneRe.MatchString(r.URL.Path):
		u.GetOne(w, r)
		break
	default:
		response := helpers.NewResponseError("Method not allowed", http.StatusMethodNotAllowed)
		response.SendResponse(w)
	}
}

func (u *UserRoutes) GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := u.userService.GetAll()
	if err != nil {
		response := helpers.NewResponseError(err.Error(), http.StatusInternalServerError)
		response.SendResponse(w)
		return
	}
	response := helpers.NewResponse(data, http.StatusOK)
	response.SendResponse(w)
}

func (u *UserRoutes) GetOne(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(fmt.Sprint(helpers.GetUrlParam(r.URL.Path, userGetOneRe)))
	data := u.userService.GetOne(id)
	if data.Id == 0 {
		response := helpers.NewResponseError("User not found", http.StatusNotFound)
		response.SendResponse(w)
		return
	}
	response := helpers.NewResponse(data, http.StatusOK)
	response.SendResponse(w)
}

func (u *UserRoutes) Register(w http.ResponseWriter, r *http.Request) {
	var user models.UserInput

	if err := helpers.ParseBodyToStruct(r, &user); err != nil {
		response := helpers.NewResponseError(err.Error(), http.StatusBadRequest)
		response.SendResponse(w)
		return
	}

	if u.verifyUserByUsername(user.Username) {
		response := helpers.NewResponseError("User already exists", http.StatusConflict)
		response.SendResponse(w)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		response := helpers.NewResponseError("Internal error in the server consult the admin", http.StatusInternalServerError)
		response.SendResponse(w)
		return
	}

	user.Password = string(hashedPassword)

	err = u.userService.Create(models.User{
		Id:       user.Id,
		Username: user.Username,
		Password: user.Password,
	})
	if err != nil {
		response := helpers.NewResponseError(err.Error(), http.StatusInternalServerError)
		response.SendResponse(w)
		return
	}

	response := helpers.NewResponse(nil, http.StatusCreated)
	response.SendResponse(w)
}

func (u *UserRoutes) Login(w http.ResponseWriter, r *http.Request) {
	u.userService = &services.UserService{
		DB: u.DB,
	}
	w.Header().Set("Content-Type", "application/json")

	var user models.UserInput

	if err := helpers.ParseBodyToStruct(r, &user); err != nil {
		response := helpers.NewResponseError(err.Error(), http.StatusBadRequest)
		response.SendResponse(w)
		return
	}

	foundUser := u.userService.GetOneByUsername(user.Username)
	if foundUser.Id == 0 {
		response := helpers.NewResponseError("Invalid password or user", http.StatusUnauthorized)
		response.SendResponse(w)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		response := helpers.NewResponseError("Invalid password or user", http.StatusUnauthorized)
		response.SendResponse(w)
		return
	}

	token, err := helpers.CreateJwtToken(models.User{
		Id:       user.Id,
		Username: user.Username,
		Password: user.Password,
	})
	if err != nil {
		response := helpers.NewResponseError("Unable to generate access token", http.StatusInternalServerError)
		response.SendResponse(w)
		return
	}

	response := helpers.NewResponse(token, http.StatusOK)
	response.SendResponse(w)
}

func (u *UserRoutes) verifyUserByUsername(username string) bool {
	user := u.userService.GetOneByUsername(username)
	if user.Id != 0 {
		return true
	} else {
		return false
	}
}
