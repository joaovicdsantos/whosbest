package session

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/joaovicdsantos/whosbest-api/app/helpers"
	"github.com/joaovicdsantos/whosbest-api/app/models"
	"github.com/joaovicdsantos/whosbest-api/app/services"
)

type Methods struct {
	DB *sql.DB
}

func (m *Methods) verifyMethodAndRun(webSocketInput models.WebSocketInput) (models.WebSocketResponse, error) {
	if webSocketInput.Method == "" {
		return models.WebSocketResponse{}, fmt.Errorf("invalid method")
	}

	switch strings.ToLower(webSocketInput.Method) {
	case "vote":
		return m.vote(webSocketInput)
	default:
		return models.WebSocketResponse{}, fmt.Errorf("invalid method")
	}
}

func (m *Methods) vote(webSocketInput models.WebSocketInput) (models.WebSocketResponse, error) {
	var competitorVote models.CompetitorVote
	if err := helpers.ParseMapToStruct(webSocketInput.Value, &competitorVote); err != nil {
		return models.WebSocketResponse{}, fmt.Errorf("invalid competitor")
	}
	
	competitorService := services.CompetitorService{DB: m.DB}
	leaderboardService := services.LeaderboardService{DB: m.DB}
	
	competitor := competitorService.GetOne(competitorVote.Competitor)
	if competitor.Id == 0 {
		return models.WebSocketResponse{}, fmt.Errorf("competitor not found")
	}
	competitorService.Vote(competitor)
	competitor.Votes = competitor.Votes + 1
	
	resultMap, err := helpers.StructToMap(leaderboardService.GetOne(competitor.Leaderboard))
	if err != nil {
		return models.WebSocketResponse{}, err
	}

	var webSocketResponse models.WebSocketResponse
	webSocketResponse.Data = resultMap
	
	return webSocketResponse, nil
}
