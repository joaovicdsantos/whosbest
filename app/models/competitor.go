package models

type Competitor struct {
	Id            int    `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	ImageURL      string `json:"image_url"`
	Votes         int    `json:"votes"`
	LeaderboardId int    `json:"leadboard_id"`
}
