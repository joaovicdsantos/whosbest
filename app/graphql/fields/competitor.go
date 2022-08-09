package fields

import (
	"database/sql"

	"github.com/graphql-go/graphql"
	"github.com/joaovicdsantos/whosbest-api/app/graphql/types"
	"github.com/joaovicdsantos/whosbest-api/app/models"
	"github.com/joaovicdsantos/whosbest-api/app/services"
)

type CompetitorField struct {
	DB *sql.DB
}

func (cf *CompetitorField) GetOne() *graphql.Field {
	return &graphql.Field{
		Type: types.CompetitorType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				DefaultValue: 0,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			competitorService := &services.CompetitorService{
				DB: cf.DB,
			}
			id, _ := p.Args["id"].(int)
			return competitorService.GetOne(id), nil
		},
	}
}

func (cf *CompetitorField) GetAll() *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(types.CompetitorType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			competitorService := &services.CompetitorService{
				DB: cf.DB,
			}
			return competitorService.GetAll()
		},
	}
}

func (cf *CompetitorField) Create() *graphql.Field {
	return &graphql.Field{
		Type: types.CompetitorType,
		Args: graphql.FieldConfigArgument{
			"title": &graphql.ArgumentConfig{
				Type:         graphql.String,
				DefaultValue: "",
			},
			"description": &graphql.ArgumentConfig{
				Type:         graphql.String,
				DefaultValue: "",
			},
			"image_url": &graphql.ArgumentConfig{
				Type:         graphql.String,
				DefaultValue: "",
			},
			"leaderboard": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				DefaultValue: 0,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			competitorService := &services.CompetitorService{
				DB: cf.DB,
			}
			competitor := models.Competitor{
				Title: p.Args["title"].(string),
				Description: p.Args["description"].(string),
				ImageURL: p.Args["image_url"].(string),
				Leaderboard: p.Args["leaderboard"].(int),
			}
			competitorService.Create(competitor)
			return competitor, nil
		},
	}
}

func (cf *CompetitorField) Update() *graphql.Field {
	return &graphql.Field{
		Type: types.CompetitorType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				DefaultValue: 0,
			},
			"title": &graphql.ArgumentConfig{
				Type:         graphql.String,
				DefaultValue: "",
			},
			"description": &graphql.ArgumentConfig{
				Type:         graphql.String,
				DefaultValue: "",
			},
			"image_url": &graphql.ArgumentConfig{
				Type:         graphql.String,
				DefaultValue: "",
			},
			"leaderboard": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				DefaultValue: 0,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			competitorService := &services.CompetitorService{
				DB: cf.DB,
			}
			competitor := models.Competitor{
				Id: p.Args["id"].(int),
				Title: p.Args["title"].(string),
				Description: p.Args["description"].(string),
				ImageURL: p.Args["image_url"].(string),
				Leaderboard: p.Args["leaderboard"].(int),
			}
			competitorService.Update(competitor)
			return competitor, nil
		},
	}
}


func (cf *CompetitorField) Delete() *graphql.Field {
	return &graphql.Field{
		Type: types.CompetitorType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:         graphql.Int,
				DefaultValue: 0,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			competitorService := &services.CompetitorService{
				DB: cf.DB,
			}
			id, _ := p.Args["id"].(int)
			competitor := competitorService.GetOne(id)
			competitorService.Delete(competitor)
			return competitor, nil
		},
	}
}
