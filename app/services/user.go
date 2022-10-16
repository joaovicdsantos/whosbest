package services

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/joaovicdsantos/whosbest-api/app/models"
)

type UserService struct {
	DB *sql.DB
}

func (s *UserService) GetAll() ([]models.User, error) {
	var users []models.User

	sql := "SELECT id, username, password FROM Users"
	result, err := s.DB.Query(sql)
	if err != nil {
		return []models.User{}, fmt.Errorf("error querying all users")
	}

	for result.Next() {
		var user models.User

		err = result.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			return []models.User{}, fmt.Errorf("error querying all users")
		}

		s.loadLeaderboards(&user)

		users = append(users, user)
	}

	if users == nil {
		users = []models.User{}
	}

	return users, nil
}

func (s *UserService) GetOne(id int) models.User {
	var user models.User

	sql := "SELECT id, username, password FROM Users WHERE Id = $1"
	err := s.DB.QueryRow(sql, id).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return models.User{}
	}

	s.loadLeaderboards(&user)

	return user
}

func (s *UserService) GetOneForLeaderboard(id int) models.User {
	var user models.User

	sql := "SELECT id, username, password FROM Users WHERE Id = $1"
	err := s.DB.QueryRow(sql, id).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return models.User{}
	}

	return user
}

func (s *UserService) GetOneByUsername(username string) models.User {
	var user models.User

	sql := "SELECT id, username, password FROM Users WHERE Username = $1"
	err := s.DB.QueryRow(sql, username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return models.User{}
	}

	return user
}

func (s *UserService) Create(user models.User) error {
	sql := "INSERT INTO Users (username, password) VALUES ($1, $2)"
	insert, err := s.DB.Prepare(sql)
	if err != nil {
		return fmt.Errorf("error creating user")
	}
	_, err = insert.Exec(strings.TrimSpace(user.Username), user.Password)
	if err != nil {
		return fmt.Errorf("error creating user")
	}
	return nil
}

func (s *UserService) loadLeaderboards(user *models.User) {
	var leaderboarService LeaderboardService
	leaderboarService.DB = s.DB
	leaderboards, _ := leaderboarService.GetAllByCreatorId(user.Id)
	user.Leaderboards = &leaderboards
}
