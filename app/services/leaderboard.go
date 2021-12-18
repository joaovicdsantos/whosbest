package services

import (
	"database/sql"
	"fmt"

	"github.com/joaovicdsantos/whosbest-api/app/models"
)

type LeaderboardService struct {
	DB *sql.DB
}

func (s *LeaderboardService) GetAll() ([]models.Leaderboard, error) {
	var leaderboards []models.Leaderboard

	sql := "SELECT * FROM Leaderboard"
	result, err := s.DB.Query(sql)
	if err != nil {
		return []models.Leaderboard{}, fmt.Errorf("Error querying all leaderboard")
	}

	for result.Next() {
		var leaderboard models.Leaderboard
		var creatorId int

		err = result.Scan(&leaderboard.Id, &leaderboard.Title, &leaderboard.Description, &creatorId)
		if err != nil {
			return []models.Leaderboard{}, fmt.Errorf("Error querying all leaderboard")
		}

		s.loadUser(&leaderboard, creatorId)

		leaderboards = append(leaderboards, leaderboard)
	}

	return leaderboards, nil
}

func (s *LeaderboardService) GetOne(id int) models.Leaderboard {
	var leaderboard models.Leaderboard
	var creatorId int

	sql := "SELECT * FROM Leaderboard WHERE Id = $1"
	err := s.DB.QueryRow(sql, id).Scan(&leaderboard.Id, &leaderboard.Title, &leaderboard.Description, &creatorId)
	if err != nil {
		return models.Leaderboard{}
	}

	s.loadUser(&leaderboard, creatorId)

	return leaderboard
}

func (s *LeaderboardService) Create(leaderboard models.Leaderboard) error {
	sql := "INSERT INTO Leaderboard (title, description, creatorId) VALUES ($1, $2, $3)"
	insert, err := s.DB.Prepare(sql)
	if err != nil {
		return fmt.Errorf("Error creating leaderboard")
	}
	_, err = insert.Exec(leaderboard.Title, leaderboard.Description, leaderboard.Creator.Id)
	if err != nil {
		return fmt.Errorf("Error creating leaderboard")
	}
	return nil
}

func (s *LeaderboardService) Update(leaderboard models.Leaderboard) error {
	sql := "UPDATE Leaderboard SET title = $1, description = $2 WHERE Id = $3"
	update, err := s.DB.Prepare(sql)
	if err != nil {
		return fmt.Errorf("Error updating leaderboard")
	}
	_, err = update.Exec(leaderboard.Title, leaderboard.Description, leaderboard.Id)
	if err != nil {
		return fmt.Errorf("Error updating leaderboard")
	}
	return nil
}

func (s *LeaderboardService) Delete(leaderboard models.Leaderboard) error {
	sql := "DELETE FROM Leaderboard WHERE Id = $1"
	del, err := s.DB.Prepare(sql)
	if err != nil {
		return fmt.Errorf("Error deleting leaderboard")
	}
	_, err = del.Exec(leaderboard.Id)
	if err != nil {
		return fmt.Errorf("Error deleting leaderboard")
	}
	return nil
}

func (s *LeaderboardService) loadUser(leaderboard *models.Leaderboard, id int) {
	var userService UserService
	userService.DB = s.DB
	user := userService.GetOne(id)
	leaderboard.Creator = &user
}
