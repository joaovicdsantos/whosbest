package handler

import (
	"fmt"
	"net/http"
	"regexp"
)

var (
	getAllRe   = regexp.MustCompile(`^\/user[\/]*$`)
	getOneRe   = regexp.MustCompile(`^\/user\/(\d+)$`)
	registerRe = regexp.MustCompile(`^\/user[\/]*$`)
)

type UserRoutes struct{}

func (u *UserRoutes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	fmt.Fprintf(w, "GetAll route from user")
}

func (u *UserRoutes) GetOne(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GetOne route from user")
}

func (u *UserRoutes) Register(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Register route from user")
}

func (u *UserRoutes) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Login route from user")
}
