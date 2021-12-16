package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/joaovicdsantos/whosbest-api/app/helpers"
	"github.com/joaovicdsantos/whosbest-api/app/models"
	"github.com/joaovicdsantos/whosbest-api/app/services"
	"golang.org/x/crypto/bcrypt"
)

var (
	getAllRe   = regexp.MustCompile(`^\/user[\/]*$`)
	getOneRe   = regexp.MustCompile(`^\/user\/(\d+)$`)
	registerRe = regexp.MustCompile(`^\/user[\/]*$`)
)

type UserRoutes struct {
	DB          *sql.DB
	userService *services.UserService
	Payload     interface{}
}

func (u *UserRoutes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u.userService = &services.UserService{
		DB: u.DB,
	}
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost && registerRe.MatchString(r.URL.Path) {
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

	// Authorized routes
	switch {
	case r.Method == http.MethodGet && getAllRe.MatchString(r.URL.Path):
		u.GetAll(w, r)
		break
	case r.Method == http.MethodGet && getOneRe.MatchString(r.URL.Path):
		u.GetOne(w, r)
		break
	default:
		response := helpers.NewResponseError("Method not allowed", http.StatusMethodNotAllowed)
		response.SendResponse(w)
	}
}

func (u *UserRoutes) GetAll(w http.ResponseWriter, r *http.Request) {
	data := u.userService.GetAll()
	response := helpers.NewResponse(data, http.StatusOK)
	response.SendResponse(w)
}

func (u *UserRoutes) GetOne(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(fmt.Sprint(helpers.GetUrlParam(r.URL.Path, getOneRe)))
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
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		response := helpers.NewResponseError("Invalid requisition body", http.StatusBadRequest)
		response.SendResponse(w)
		return
	}
	var user models.User
	if err = json.Unmarshal(body, &user); err != nil {
		response := helpers.NewResponseError("Invalid requisition body", http.StatusBadRequest)
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

	u.userService.Create(user)

	response := helpers.NewResponse(nil, http.StatusCreated)
	response.SendResponse(w)
}

func (u *UserRoutes) Login(w http.ResponseWriter, r *http.Request) {
	u.userService = &services.UserService{
		DB: u.DB,
	}
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		response := helpers.NewResponseError("Invalid requisition body", http.StatusBadRequest)
		response.SendResponse(w)
		return
	}
	var user models.User
	if err = json.Unmarshal(body, &user); err != nil {
		response := helpers.NewResponseError("Invalid requisition body", http.StatusBadRequest)
		response.SendResponse(w)
		return
	}

	foundUser := u.userService.GetOneByUsername(user.Username)
	if foundUser.Id == 0 {
		response := helpers.NewResponseError("Invalid password or user", http.StatusUnauthorized)
		response.SendResponse(w)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		response := helpers.NewResponseError("Invalid password or user", http.StatusUnauthorized)
		response.SendResponse(w)
		return
	}

	token, err := helpers.CreateJwtToken(user)
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
