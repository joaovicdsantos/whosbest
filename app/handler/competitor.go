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
	competitorGetAllRe = regexp.MustCompile(`^\/competitor[\/]*$`)
	competitorGetOneRe = regexp.MustCompile(`^\/competitor\/(\d+)$`)
	competitorCreateRe = competitorGetAllRe
	competitorUpdateRe = competitorGetOneRe
	competitorDeleteRe = competitorGetOneRe
)

type CompetitorRoutes struct {
	DB                *sql.DB
	competitorService *services.CompetitorService
	Payload           map[string]interface{}
}

func (c *CompetitorRoutes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.competitorService = &services.CompetitorService{
		DB: c.DB,
	}
	w.Header().Set("Content-Type", "application/json")

	unparsedToken, ok := r.Header["Authorization"]
	if !ok {
		response := helpers.NewResponseError("Token is not valid", http.StatusUnauthorized)
		response.SendResponse(w)
		return
	}

	var err error
	c.Payload, err = helpers.ParseJwtToken(fmt.Sprint(unparsedToken[0]))
	if err != nil {
		response := helpers.NewResponseError(err.Error(), http.StatusUnauthorized)
		response.SendResponse(w)
		return
	}

	user := helpers.GetCurrentUser(c.Payload, c.DB)
	if user.Id == 0 {
		response := helpers.NewResponseError("Invalid user", http.StatusUnauthorized)
		response.SendResponse(w)
		return

	}

	// Authorized routes
	switch {
	case r.Method == http.MethodGet && competitorGetAllRe.MatchString(r.URL.Path):
		c.GetAll(w, r)
		break
	case r.Method == http.MethodGet && competitorGetOneRe.MatchString(r.URL.Path):
		c.GetOne(w, r)
		break
	case r.Method == http.MethodPost && competitorCreateRe.MatchString(r.URL.Path):
		c.Create(w, r)
		break
	case r.Method == http.MethodPut && competitorGetOneRe.MatchString(r.URL.Path):
		c.Update(w, r)
		break
	case r.Method == http.MethodDelete && competitorDeleteRe.MatchString(r.URL.Path):
		c.Delete(w, r)
		break
	default:
		response := helpers.NewResponseError("Method not allowed", http.StatusMethodNotAllowed)
		response.SendResponse(w)
	}
}

func (c *CompetitorRoutes) GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := c.competitorService.GetAll()
	if err != nil {
		response := helpers.NewResponseError(err.Error(), http.StatusInternalServerError)
		response.SendResponse(w)
		return
	}
	response := helpers.NewResponse(data, http.StatusOK)
	response.SendResponse(w)
}

func (c *CompetitorRoutes) GetOne(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(fmt.Sprint(helpers.GetUrlParam(r.URL.Path, competitorGetOneRe)))
	data := c.competitorService.GetOne(id)
	if data.Id == 0 {
		response := helpers.NewResponseError("Competitor not found", http.StatusNotFound)
		response.SendResponse(w)
		return
	}
	response := helpers.NewResponse(data, http.StatusOK)
	response.SendResponse(w)
}

func (c *CompetitorRoutes) Create(w http.ResponseWriter, r *http.Request) {

	var competitor models.Competitor
	if err := helpers.ParseBodyToStruct(r, &competitor); err != nil {
		response := helpers.NewResponseError(err.Error(), http.StatusBadRequest)
		response.SendResponse(w)
		return
	}

	leaderboardId := c.getLeadboard(competitor.LeaderboardId).Id
	if leaderboardId == 0 {
		response := helpers.NewResponseError(fmt.Sprintf("Leaderboard id %d invalid", competitor.LeaderboardId), http.StatusNotFound)
		response.SendResponse(w)
		return
	}

	userId := helpers.GetCurrentUser(c.Payload, c.DB).Id
	creatorLeadboardId := c.getLeadboard(leaderboardId).Creator.Id
	if creatorLeadboardId != userId {
		response := helpers.NewResponseError("You are not authorized for this", http.StatusForbidden)
		response.SendResponse(w)
		return
	}

	c.competitorService.Create(competitor)

	response := helpers.NewResponse(nil, http.StatusCreated)
	response.SendResponse(w)
}

func (c *CompetitorRoutes) Update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(fmt.Sprint(helpers.GetUrlParam(r.URL.Path, competitorGetOneRe)))
	data := c.competitorService.GetOne(id)
	if data.Id == 0 {
		response := helpers.NewResponseError("Competitor not found", http.StatusNotFound)
		response.SendResponse(w)
		return
	}

	userId := helpers.GetCurrentUser(c.Payload, c.DB).Id
	creatorLeadboardId := c.getLeadboard(data.LeaderboardId).Creator.Id
	if creatorLeadboardId != userId {
		response := helpers.NewResponseError("You are not authorized for this", http.StatusForbidden)
		response.SendResponse(w)
		return
	}

	var competitor models.Competitor
	if err := helpers.ParseBodyToStruct(r, &competitor); err != nil {
		response := helpers.NewResponseError(err.Error(), http.StatusBadRequest)
		response.SendResponse(w)
		return
	}

	data.Update(competitor)
	err := c.competitorService.Update(data)
	if err != nil {
		response := helpers.NewResponseError(err.Error(), http.StatusInternalServerError)
		response.SendResponse(w)
		return
	}

	response := helpers.NewResponse(data, http.StatusAccepted)
	response.SendResponse(w)
}

func (c *CompetitorRoutes) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(fmt.Sprint(helpers.GetUrlParam(r.URL.Path, competitorGetOneRe)))
	data := c.competitorService.GetOne(id)
	if data.Id == 0 {
		response := helpers.NewResponseError("Competitor not found", http.StatusNotFound)
		response.SendResponse(w)
		return
	}

	userId := helpers.GetCurrentUser(c.Payload, c.DB).Id
	creatorLeadboardId := c.getLeadboard(data.LeaderboardId).Creator.Id
	if creatorLeadboardId != userId {
		response := helpers.NewResponseError("You are not authorized for this", http.StatusForbidden)
		response.SendResponse(w)
		return
	}

	err := c.competitorService.Delete(data)
	if err != nil {
		response := helpers.NewResponseError(err.Error(), http.StatusInternalServerError)
		response.SendResponse(w)
		return
	}

	response := helpers.NewResponse("Deleted", http.StatusNoContent)
	response.SendResponse(w)
}

func (c *CompetitorRoutes) getLeadboard(id int) models.Leaderboard {
	var leaderboardService services.LeaderboardService

	leaderboardService.DB = c.DB

	return leaderboardService.GetOne(id)
}
