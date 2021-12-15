package services

import (
	"database/sql"
	"fmt"

	"github.com/joaovicdsantos/whosbest-api/app/models"
)

type UserService struct {
	DB *sql.DB
}

func (u *UserService) GetAll() []models.User {
	sql := "SELECT * FROM Users"
	result, err := u.DB.Query(sql)
	if err != nil {
		panic("SQL ERROR")
	}

	var users []models.User
	for result.Next() {
		var user models.User

		err = result.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			panic("SQL ERROR")
		}

		users = append(users, user)
	}

	return users
}

func (u *UserService) GetOne(id int) models.User {
	sql := "SELECT * FROM Users WHERE Id = $1"
	var user models.User
	err := u.DB.QueryRow(sql, id).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return models.User{}
	}

	return user
}

func (u *UserService) GetOneByUsername(username string) models.User {
	sql := "SELECT * FROM Users WHERE Username = $1"

	var user models.User
	err := u.DB.QueryRow(sql, username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return models.User{}
	}

	return user
}

func (u *UserService) Create(user models.User) {
	sql := "INSERT INTO Users (username, password) VALUES ($1, $2)"
	insert, err := u.DB.Prepare(sql)
	if err != nil {
		panic("SQL ERROR")
	}
	result, err := insert.Exec(user.Username, user.Password)
	if err != nil {
		panic("SQL ERROR")
	}
	fmt.Println(result)
}
