package models

type Competitor struct {
	Id          int    `json:"id"`
	Title       string `json:"title" validate:"max=80"`
	Description string `json:"description" validate:"min=10,max=600"`
	ImageURL    string `json:"image_url" validate:"max=120"`
	Votes       int    `json:"votes" validate:"gte=0"`
	Leaderboard int    `json:"leaderboard"`
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
