package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

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
}

func (u *UserRoutes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u.userService = &services.UserService{
		DB: u.DB,
	}
	w.Header().Set("Content-Type", "application/json")
	switch {
	case r.Method == http.MethodGet && getAllRe.MatchString(r.URL.Path):
		u.GetAll(w, r)
		break
	case r.Method == http.MethodGet && getOneRe.MatchString(r.URL.Path):
		u.GetOne(w, r)
		break
	case r.Method == http.MethodPost && registerRe.MatchString(r.URL.Path):
		u.Register(w, r)
		break
	}
}

func (u *UserRoutes) GetAll(w http.ResponseWriter, r *http.Request) {
	response, err := json.Marshal(map[string][]models.User{
		"users": u.userService.GetAll(),
	})
	if err != nil {
		panic("Error on format response")
	}
	w.Write(response)
}

func (u *UserRoutes) GetOne(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GetOne route from user")
}

func (u *UserRoutes) Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var user models.User
	if err = json.Unmarshal(body, &user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if u.verifyUserByUsername(user.Username) {
		w.WriteHeader(http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)
	u.userService.Create(user)
}

func (u *UserRoutes) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Login route from user")
}

func (u *UserRoutes) verifyUserByUsername(username string) bool {
	user := u.userService.GetOneByUsername(username)
	if user.Id != 0 {
		return true
	} else {
		return false
	}
}
