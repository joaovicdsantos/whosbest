package app

import (
	"log"
	"net/http"

	"github.com/joaovicdsantos/whosbest-api/app/handler"
)

type App struct {
	mux *http.ServeMux
}

func (a *App) Initialize() {
	a.mux = http.NewServeMux()
	a.SetRoutes()
}

func (a *App) SetRoutes() {
	// User routes
	a.mux.Handle("/user", &handler.UserRoutes{})
	a.mux.Handle("/user/", &handler.UserRoutes{})
	a.mux.HandleFunc("/login", new(handler.UserRoutes).Login)
}

func (a *App) Run() {
	log.Fatal(http.ListenAndServe(":3001", a.mux))
}
