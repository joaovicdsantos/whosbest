package fields

import (
	"database/sql"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/graphql-go/graphql"
	"github.com/joaovicdsantos/whosbest-api/app/graphql/types"
	"github.com/joaovicdsantos/whosbest-api/app/helpers"
	"github.com/joaovicdsantos/whosbest-api/app/models"
	"github.com/joaovicdsantos/whosbest-api/app/services"
)

type CompetitorField struct {
	DB                *sql.DB
	competitorService *services.CompetitorService
}

func NewCompetitorField(db *sql.DB) *CompetitorField {
	return &CompetitorField{
		DB: db,
		competitorService: &services.CompetitorService{
			DB: db,
		},
	}
}

func (cf *CompetitorField) GetOne() *graphql.Field {
	return &graphql.Field{
		Type: types.CompetitorType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Competitor identifier",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id, _ := p.Args["id"].(int)
			competitor := cf.competitorService.GetOne(id)
			if competitor.Id == 0 {
				return nil, fmt.Errorf("competitor not found")
			}
			return competitor, nil
		},
	}
}

func (cf *CompetitorField) GetAll() *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(types.CompetitorType),
		Resolve: func(_ graphql.ResolveParams) (interface{}, error) {
			return cf.competitorService.GetAll()
		},
	}
}

func (cf *CompetitorField) Create() *graphql.Field {
	return &graphql.Field{
		Type: types.CompetitorType,
		Args: graphql.FieldConfigArgument{
			"title": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Competitor's name or title",
			},
			"description": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "A description for competitor attributes",
			},
			"image_url": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "A image for represent the competitor",
			},
			"leaderboard": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Competitor's leaderboard",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			userID := p.Context.Value("user_id").(int)
			var competitor models.Competitor
			if err := helpers.ParseMapToStruct(p.Args, &competitor); err != nil {
				return nil, fmt.Errorf("invalid create params")
			}

			leaderboard := cf.getLeaderboard(competitor.Leaderboard)
			if leaderboard.Id == 0 {
				return nil, fmt.Errorf("leaderboard id %d invalid", competitor.Leaderboard)
			}

			if leaderboard.Creator.Id != userID {
				return nil, fmt.Errorf("you are not authorized for this")
			}

			validate := validator.New()
			err := validate.Struct(competitor)
			if err != nil {
				return nil, err
			}

			createdCompetitor, err := cf.competitorService.Create(competitor)
			if err != nil {
				return nil, fmt.Errorf("error on competitor creation")
			}

			return createdCompetitor, nil
		},
	}
}

func (cf *CompetitorField) Update() *graphql.Field {
	return &graphql.Field{
		Type: types.CompetitorType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Competitor identifier",
			},
			"title": &graphql.ArgumentConfig{
				Type:         graphql.String,
				DefaultValue: "",
				Description:  "Competitor's name or title",
			},
			"description": &graphql.ArgumentConfig{
				Type:         graphql.String,
				DefaultValue: "",
				Description:  "A description for competitor attributes",
			},
			"image_url": &graphql.ArgumentConfig{
				Type:         graphql.String,
				DefaultValue: "",
				Description:  "A image for represent the competitor",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			userID := p.Context.Value("user_id").(int)
			var competitor models.Competitor
			if err := helpers.ParseMapToStruct(p.Args, &competitor); err != nil {
				return nil, fmt.Errorf("invalid update params")
			}

			savedCompetitor := cf.competitorService.GetOne(userID)
			if savedCompetitor.Id == 0 {
				return nil, fmt.Errorf("invalid user")
			}

			if !cf.isAuthorized(userID, savedCompetitor) {
				return nil, fmt.Errorf("you are not authorized for this or the resource does not exist")
			}

			validate := validator.New()
			err := validate.Struct(competitor)
			if err != nil {
				return nil, err
			}

			updatedCompetitor, err := cf.competitorService.Update(competitor)
			if err != nil {
				return nil, fmt.Errorf("error on competitor update")
			}
			return updatedCompetitor, nil
		},
	}
}

func (cf *CompetitorField) Delete() *graphql.Field {
	return &graphql.Field{
		Type: types.CompetitorType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "A description for competitor attributes",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			userID := p.Context.Value("user_id").(int)
			id, _ := p.Args["id"].(int)

			competitor := cf.competitorService.GetOne(id)
			if !cf.isAuthorized(userID, competitor) {
				return nil, fmt.Errorf("you are not authorized for this or the resource does not exist")
			}
			cf.competitorService.Delete(competitor)
			return competitor, nil
		},
	}
}

func (cf *CompetitorField) isAuthorized(userID int, competitor models.Competitor) bool {
	return cf.getLeaderboard(competitor.Leaderboard).Creator.Id == userID
}

func (cf *CompetitorField) getLeaderboard(id int) models.Leaderboard {
	leaderboardService := &services.LeaderboardService{
		DB: cf.DB,
	}
	return leaderboardService.GetOne(id)
}
