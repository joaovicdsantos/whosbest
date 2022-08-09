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
)

var (
	leaderboardGetAllRe = regexp.MustCompile(`^\/leaderboard[\/]*$`)
	leaderboardGetOneRe = regexp.MustCompile(`^\/leaderboard\/(\d+)$`)
	leaderboardCreateRe = leaderboardGetAllRe
	leaderboardUpdateRe = leaderboardGetOneRe
	leaderboardDeleteRe = leaderboardGetOneRe
)

type LeaderboardRoutes struct {
	DB                 *sql.DB
	leaderboardService *services.LeaderboardService
	Payload            map[string]interface{}
}

func (l *LeaderboardRoutes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.leaderboardService = &services.LeaderboardService{
		DB: l.DB,
	}
	w.Header().Set("Content-Type", "application/json")

	unparsedToken, ok := r.Header["Authorization"]
	if !ok {
		response := helpers.NewResponseError("Token is not valid", http.StatusUnauthorized)
		response.SendResponse(w)
		return
	}

	var err error
	l.Payload, err = helpers.ParseJwtToken(fmt.Sprint(unparsedToken[0]))
	if err != nil {
		response := helpers.NewResponseError(err.Error(), http.StatusUnauthorized)
		response.SendResponse(w)
		return
	}

	user := helpers.GetCurrentUser(l.Payload, l.DB)
	if user.Id == 0 {
		response := helpers.NewResponseError("Invalid user", http.StatusUnauthorized)
		response.SendResponse(w)
		return

	}

	// Authorized routes
	switch {
	case r.Method == http.MethodGet && leaderboardGetAllRe.MatchString(r.URL.Path):
		l.GetAll(w, r)
	case r.Method == http.MethodGet && leaderboardGetOneRe.MatchString(r.URL.Path):
		l.GetOne(w, r)
	case r.Method == http.MethodPost && leaderboardCreateRe.MatchString(r.URL.Path):
		l.Create(w, r)
	case r.Method == http.MethodPut && leaderboardGetOneRe.MatchString(r.URL.Path):
		l.Update(w, r)
	case r.Method == http.MethodDelete && leaderboardDeleteRe.MatchString(r.URL.Path):
		l.Delete(w, r)
	default:
		response := helpers.NewResponseError("Method not allowed", http.StatusMethodNotAllowed)
		response.SendResponse(w)
	}
}

func (l *LeaderboardRoutes) GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := l.leaderboardService.GetAll()
	if err != nil {
		response := helpers.NewResponseError(err.Error(), http.StatusInternalServerError)
		response.SendResponse(w)
		return
	}
	response := helpers.NewResponse(data, http.StatusOK)
	response.SendResponse(w)
}

func (l *LeaderboardRoutes) GetOne(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(fmt.Sprint(helpers.GetUrlParam(r.URL.Path, leaderboardGetOneRe)))
	data := l.leaderboardService.GetOne(id)
	if data.Id == 0 {
		response := helpers.NewResponseError("Leaderboard not found", http.StatusNotFound)
		response.SendResponse(w)
		return
	}
	response := helpers.NewResponse(data, http.StatusOK)
	response.SendResponse(w)
}

func (l *LeaderboardRoutes) Create(w http.ResponseWriter, r *http.Request) {

	var leaderboard models.Leaderboard
	if err := helpers.ParseBodyToStruct(r, &leaderboard); err != nil {
		response := helpers.NewResponseError(err.Error(), http.StatusBadRequest)
		response.SendResponse(w)
		return
	}

	user := helpers.GetCurrentUser(l.Payload, l.DB)

	leaderboard.Creator = &models.User{Id: user.Id}
	l.leaderboardService.Create(leaderboard)

	response := helpers.NewResponse(nil, http.StatusCreated)
	response.SendResponse(w)
}

func (l *LeaderboardRoutes) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(fmt.Sprint(helpers.GetUrlParam(r.URL.Path, leaderboardGetOneRe)))
	data := l.leaderboardService.GetOne(id)
	if data.Id == 0 {
		response := helpers.NewResponseError("Leaderboard not found", http.StatusNotFound)
		response.SendResponse(w)
		return
	}

	userId := helpers.GetCurrentUser(l.Payload, l.DB).Id
	if data.Creator.Id != userId {
		response := helpers.NewResponseError("You are not authorized for this", http.StatusForbidden)
		response.SendResponse(w)
		return
	}

	var leaderboard models.Leaderboard
	if err := helpers.ParseBodyToStruct(r, &leaderboard); err != nil {
		response := helpers.NewResponseError(err.Error(), http.StatusBadRequest)
		response.SendResponse(w)
		return
	}

	data.Update(leaderboard)
	err := l.leaderboardService.Update(data)
	if err != nil {
		response := helpers.NewResponseError(err.Error(), http.StatusInternalServerError)
		response.SendResponse(w)
		return
	}

	response := helpers.NewResponse(data, http.StatusAccepted)
	response.SendResponse(w)
}

func (l *LeaderboardRoutes) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(fmt.Sprint(helpers.GetUrlParam(r.URL.Path, leaderboardGetOneRe)))
	data := l.leaderboardService.GetOne(id)
	if data.Id == 0 {
		response := helpers.NewResponseError("Leaderboard not found", http.StatusNotFound)
		response.SendResponse(w)
		return
	}

	userId := helpers.GetCurrentUser(l.Payload, l.DB).Id
	if data.Creator.Id != userId {
		response := helpers.NewResponseError("You are not authorized for this", http.StatusForbidden)
		response.SendResponse(w)
		return
	}

	err := l.leaderboardService.Delete(data)
	if err != nil {
		response := helpers.NewResponseError(err.Error(), http.StatusInternalServerError)
		response.SendResponse(w)
		return
	}

	response := helpers.NewResponse("Deleted", http.StatusNoContent)
	response.SendResponse(w)
}
