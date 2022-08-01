package app

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/joaovicdsantos/whosbest-api/app/graphql"
	"github.com/joaovicdsantos/whosbest-api/app/handler"
	"github.com/joaovicdsantos/whosbest-api/app/helpers"
	"github.com/joaovicdsantos/whosbest-api/app/models"
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
	userRoutes := new(handler.UserRoutes)
	userRoutes.DB = db
	a.mux.Handle("/user", userRoutes)
	a.mux.Handle("/user/", userRoutes)
	a.mux.HandleFunc("/login", userRoutes.Login)

	// Leaderboard routes
	leaderboardRoutes := new(handler.LeaderboardRoutes)
	leaderboardRoutes.DB = db
	a.mux.Handle("/leaderboard", leaderboardRoutes)
	a.mux.Handle("/leaderboard/", leaderboardRoutes)

	// Competitor routes
	competitorRoutes := new(handler.CompetitorRoutes)
	competitorRoutes.DB = db
	a.mux.Handle("/competitor", competitorRoutes)
	a.mux.Handle("/competitor/", competitorRoutes)

	var graphql = new(graphql.GraphQL)
	graphql.Initialize(db)
	a.mux.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		var graphqlIn models.GraphQL
		if err := helpers.ParseBodyToStruct(r, &graphqlIn); err != nil {
			response := helpers.NewResponseError(err.Error(), http.StatusBadRequest)
			response.SendResponse(w)
			return
		}
		response := helpers.NewResponse(graphql.ExecuteQuery(graphqlIn.Query), http.StatusOK)
		response.SendResponse(w)
	})
}

func (a *App) Run() {
	log.Fatal(http.ListenAndServe(":3001", a.mux))
}
