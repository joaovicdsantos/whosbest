package models

type Competitor struct {
	Id            int    `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	ImageURL      string `json:"image_url"`
	Votes         int    `json:"votes"`
	Leaderboard   int    `json:"leaderboard"`
}

func (c *Competitor) Update(competitor Competitor) {
	if competitor.Title != "" && c.Title != competitor.Title {
		c.Title = competitor.Title
	}
	if competitor.Description != "" && c.Description != competitor.Description {
		c.Description = competitor.Description
	}
	if competitor.ImageURL != "" && c.ImageURL != competitor.ImageURL {
		c.ImageURL = competitor.ImageURL
	}

}
