package models

type Leaderboard struct {
	Id          int           `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Creator     *User         `json:"creator"`
	Competitors *[]Competitor `json:"competitors"`
}

func (l *Leaderboard) Update(leaderboard Leaderboard) {
	if leaderboard.Title != "" && l.Title != leaderboard.Title {
		l.Title = leaderboard.Title
	}
	if leaderboard.Description != "" && l.Description != leaderboard.Description {
		l.Description = leaderboard.Description
	}
}
