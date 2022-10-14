package handler

import (
	"database/sql"
	"net/http"
	"regexp"

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
}

func NewUserRoutes(db *sql.DB) *UserRoutes {
	return &UserRoutes{
		DB: db,
		userService: &services.UserService{
			DB: db,
		},
	}
}

func (u *UserRoutes) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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
	return user.Id != 0
}
