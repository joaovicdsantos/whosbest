package models

type User struct {
	Id           int            `json:"id"`
	Username     string         `json:"username"`
	Password     string         `json:"-"`
	Leaderboards *[]Leaderboard `json:"leaderboard"`
}

type UserInput struct {
	User
	Password string `json:"password"`
}
