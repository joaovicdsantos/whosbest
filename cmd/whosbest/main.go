package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/joaovicdsantos/whosbest-api/internal/database"
	"github.com/joaovicdsantos/whosbest-api/internal/handler"
	"github.com/joaovicdsantos/whosbest-api/internal/session"
)

type App struct {
	mux *http.ServeMux
}

func (a *App) Initialize(db *sql.DB) {
	a.mux = http.NewServeMux()
	a.SetRoutes(db)
}

func (a *App) SetRoutes(db *sql.DB) {
	// User routes
	userRoutes := handler.NewUserRoutes(db)
	a.mux.HandleFunc("/register", userRoutes.Register)
	a.mux.HandleFunc("/login", userRoutes.Login)

	// Graphql route
	graphqlRoute := handler.GraphqlRoute{DB: db}
	a.mux.HandleFunc("/graphql", graphqlRoute.HandleGraphqlRequest)

	// Websocket
	hub := session.NewHub()
	go hub.Run()
	a.mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		session.ServerWs(hub, db, w, r)
	})

}

func (a *App) Run() {
	log.Fatal(http.ListenAndServe(":3001", a.mux))
}

func main() {
	dbFactory := new(database.DBFactory)
	db := dbFactory.InitDatabase()
	defer db.Close()

	app := new(App)
	app.Initialize(db)
	app.Run()
}
