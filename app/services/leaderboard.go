package services

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/joaovicdsantos/whosbest-api/app/models"
)

type LeaderboardService struct {
	DB *sql.DB
}

func (s *LeaderboardService) GetAll() ([]models.Leaderboard, error) {
	var leaderboards []models.Leaderboard

	sql := "SELECT id, title, description, creator FROM Leaderboards"
	result, err := s.DB.Query(sql)
	if err != nil {
		return []models.Leaderboard{}, fmt.Errorf("error querying all leaderboard")
	}

	for result.Next() {
		var leaderboard models.Leaderboard
		var creatorId int

		err = result.Scan(&leaderboard.Id, &leaderboard.Title, &leaderboard.Description, &creatorId)
		if err != nil {
			return []models.Leaderboard{}, fmt.Errorf("error querying all leaderboard")
		}

		s.loadUser(&leaderboard, creatorId)
		s.loadCompetitors(&leaderboard)

		leaderboards = append(leaderboards, leaderboard)
	}

	if leaderboards == nil {
		leaderboards = []models.Leaderboard{}
	}

	return leaderboards, nil
}

func (s *LeaderboardService) GetOne(id int) models.Leaderboard {
	var leaderboard models.Leaderboard
	var creatorId int

	sql := "SELECT id, title, description, creator FROM Leaderboards WHERE Id = $1"
	err := s.DB.QueryRow(sql, id).Scan(&leaderboard.Id, &leaderboard.Title, &leaderboard.Description, &creatorId)
	if err != nil {
		return models.Leaderboard{}
	}

	s.loadUser(&leaderboard, creatorId)
	s.loadCompetitors(&leaderboard)

	return leaderboard
}

func (s *LeaderboardService) Create(leaderboard models.Leaderboard) (models.Leaderboard, error) {
	sql := "INSERT INTO Leaderboards (title, description, creator) VALUES ($1, $2, $3) RETURNING id"
	insert, err := s.DB.Prepare(sql)
	if err != nil {
		return models.Leaderboard{}, fmt.Errorf("error creating leaderboard")
	}

	err = insert.QueryRow(
		strings.TrimSpace(leaderboard.Title),
		strings.TrimSpace(leaderboard.Description),
		leaderboard.Creator.Id,
	).Scan(&leaderboard.Id)
	if err != nil {
		fmt.Println(err)
		return models.Leaderboard{}, fmt.Errorf("error creating leaderboard")
	}

	return leaderboard, nil
}

func (s *LeaderboardService) Update(leaderboard models.Leaderboard) (models.Leaderboard, error) {
	sql := "UPDATE Leaderboards SET title = $1, description = $2 WHERE Id = $3"
	current := s.GetOne(leaderboard.Id)

	if len(strings.Trim(leaderboard.Title, "")) > 0 {
		current.Title = leaderboard.Title
	}

	if len(strings.Trim(leaderboard.Description, "")) > 0 {
		current.Description = leaderboard.Description
	}

	update, err := s.DB.Prepare(sql)
	if err != nil {
		return models.Leaderboard{}, fmt.Errorf("error updating leaderboard")
	}

	_, err = update.Exec(
		strings.TrimSpace(current.Title),
		strings.TrimSpace(current.Description),
		leaderboard.Id,
	)
	if err != nil {
		return models.Leaderboard{}, fmt.Errorf("error updating leaderboard")
	}

	return current, nil
}

func (s *LeaderboardService) Delete(leaderboard models.Leaderboard) error {
	sql := "DELETE FROM Leaderboards WHERE Id = $1"
	del, err := s.DB.Prepare(sql)
	if err != nil {
		return fmt.Errorf("error deleting leaderboard")
	}
	_, err = del.Exec(leaderboard.Id)
	if err != nil {
		return fmt.Errorf("error deleting leaderboard")
	}
	return nil
}

func (s *LeaderboardService) loadUser(leaderboard *models.Leaderboard, id int) {
	var userService UserService
	userService.DB = s.DB
	user := userService.GetOne(id)
	leaderboard.Creator = &user
}

func (s *LeaderboardService) loadCompetitors(leaderboard *models.Leaderboard) {
	var competitorService CompetitorService
	competitorService.DB = s.DB
	competitors, _ := competitorService.GetAllByLeaderboardId(leaderboard.Id)
	leaderboard.Competitors = &competitors
}
