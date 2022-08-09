package services

import (
	"database/sql"
	"fmt"

	"github.com/joaovicdsantos/whosbest-api/app/models"
)

type CompetitorService struct {
	DB *sql.DB
}

func (s *CompetitorService) GetAll() ([]models.Competitor, error) {
	var competitors []models.Competitor

	sql := "SELECT * FROM Competitors"
	result, err := s.DB.Query(sql)
	if err != nil {
		return []models.Competitor{}, fmt.Errorf("error querying all competitor")
	}

	for result.Next() {
		var competitor models.Competitor

		err = result.Scan(&competitor.Id, &competitor.Title, &competitor.Description, &competitor.ImageURL, &competitor.Votes, &competitor.Leaderboard)
		if err != nil {
			return []models.Competitor{}, fmt.Errorf("error querying all competitor")
		}

		competitors = append(competitors, competitor)
	}

	if competitors == nil {
		competitors = []models.Competitor{}
	}

	return competitors, nil
}

func (s *CompetitorService) GetAllByLeaderboardId(id int) ([]models.Competitor, error) {
	var competitors []models.Competitor

	sql := "SELECT * FROM Competitors WHERE id = $1"

	result, err := s.DB.Query(sql, id)
	if err != nil {
		return []models.Competitor{}, fmt.Errorf("error querying all competitor by leaderboardId")
	}

	for result.Next() {
		var competitor models.Competitor

		err = result.Scan(&competitor.Id, &competitor.Title, &competitor.Description, &competitor.ImageURL, &competitor.Votes, &competitor.Leaderboard)
		if err != nil {
			return []models.Competitor{}, fmt.Errorf("error querying all competitor by leaderboardId")
		}

		competitors = append(competitors, competitor)
	}

	if competitors == nil {
		competitors = []models.Competitor{}
	}

	return competitors, nil
}
func (s *CompetitorService) GetOne(id int) models.Competitor {
	var competitor models.Competitor

	sql := "SELECT * FROM Competitors WHERE Id = $1"
	err := s.DB.QueryRow(sql, id).Scan(&competitor.Id, &competitor.Title, &competitor.Description, &competitor.ImageURL, &competitor.Votes, &competitor.Leaderboard)
	if err != nil {
		return models.Competitor{}
	}

	return competitor
}

func (s *CompetitorService) Create(competitor models.Competitor) error {
	sql := "INSERT INTO Competitors (title, description, imageUrl, votes, leaderboard) VALUES ($1, $2, $3, 0, $4)"
	insert, err := s.DB.Prepare(sql)
	if err != nil {
		return fmt.Errorf("error creating competitor")
	}
	_, err = insert.Exec(competitor.Title, competitor.Description, competitor.ImageURL, competitor.Leaderboard)
	if err != nil {
		return fmt.Errorf("error creating competitor")
	}
	return nil
}

func (s *CompetitorService) Update(competitor models.Competitor) error {
	sql := "UPDATE Competitors SET title = $1, description = $2, imageurl = $3 WHERE Id = $4"
	update, err := s.DB.Prepare(sql)
	if err != nil {
		return fmt.Errorf("error updating competitor")
	}
	_, err = update.Exec(competitor.Title, competitor.Description, competitor.ImageURL, competitor.Id)
	if err != nil {
		return fmt.Errorf("error updating competitor")
	}
	return nil
}

func (s *CompetitorService) Delete(competitor models.Competitor) error {
	sql := "DELETE FROM Competitors WHERE Id = $1"
	del, err := s.DB.Prepare(sql)
	if err != nil {
		return fmt.Errorf("error deleting competitor")
	}
	_, err = del.Exec(competitor.Id)
	if err != nil {
		return fmt.Errorf("error deleting competitor")
	}
	return nil
}

func (s *CompetitorService) Vote(competitor models.Competitor) error {
	sql := "UPDATE Competitors SET Votes = Votes + 1 WHERE Id = $1"
	del, err := s.DB.Prepare(sql)
	if err != nil {
		return fmt.Errorf("error voting for competitor")
	}
	_, err = del.Exec(competitor.Id)
	if err != nil {
		return fmt.Errorf("error voting for competitor")
	}
	return nil
}
