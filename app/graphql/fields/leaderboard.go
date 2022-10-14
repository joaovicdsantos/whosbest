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

type LeaderboardField struct {
	DB                 *sql.DB
	leaderboardService *services.LeaderboardService
}

func NewLeaderboardField(db *sql.DB) *LeaderboardField {
	return &LeaderboardField{
		DB: db,
		leaderboardService: &services.LeaderboardService{
			DB: db,
		},
	}
}

func (lf *LeaderboardField) GetOne() *graphql.Field {
	return &graphql.Field{
		Type: types.LeaderboardType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Leaderboard identifier",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id, _ := p.Args["id"].(int)
			leaderboard := lf.leaderboardService.GetOne(id)
			if leaderboard.Id == 0 {
				return nil, fmt.Errorf("leaderboard not found")
			}
			return leaderboard, nil
		},
	}
}

func (lf *LeaderboardField) GetAll() *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(types.LeaderboardType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return lf.leaderboardService.GetAll()
		},
	}
}

func (lf *LeaderboardField) Create() *graphql.Field {
	return &graphql.Field{
		Type: types.LeaderboardType,
		Args: graphql.FieldConfigArgument{
			"title": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "Leaderboard's title",
			},
			"description": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "A description for leaderboard",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			userID := p.Context.Value("user_id").(int)
			var leaderboard models.Leaderboard
			if err := helpers.ParseMapToStruct(p.Args, &leaderboard); err != nil {
				return nil, fmt.Errorf("invalid create params")
			}

			leaderboard.Creator = lf.getCurrentUser(userID)

			validate := validator.New()
			err := validate.Struct(leaderboard)
			if err != nil {
				return nil, err
			}

			createdLeaderboard, err := lf.leaderboardService.Create(leaderboard)
			if err != nil {
				fmt.Println(err)
				return nil, fmt.Errorf("error on leaderboard creation")
			}

			return createdLeaderboard, nil
		},
	}
}

func (lf *LeaderboardField) Update() *graphql.Field {
	return &graphql.Field{
		Type: types.LeaderboardType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Leaderboard identifier",
			},
			"title": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "Leaderboard's title",
			},
			"description": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "A description for leaderboard",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			userID := p.Context.Value("user_id").(int)
			var leaderboard models.Leaderboard
			if err := helpers.ParseMapToStruct(p.Args, &leaderboard); err != nil {
				return nil, fmt.Errorf("invalid update params")
			}

			savedLeaderboard := lf.leaderboardService.GetOne(leaderboard.Id)
			if savedLeaderboard.Id == 0 || savedLeaderboard.Creator.Id != userID {
				return nil, fmt.Errorf("you are not authorized for this or the resource does not exist")
			}

			validate := validator.New()
			err := validate.Struct(leaderboard)
			if err != nil {
				return nil, err
			}

			updatedLeaderboard, err := lf.leaderboardService.Update(leaderboard)
			if err != nil {
				return nil, fmt.Errorf("error on leaderboard update")
			}

			return updatedLeaderboard, nil
		},
	}
}

func (lf *LeaderboardField) Delete() *graphql.Field {
	return &graphql.Field{
		Type: types.LeaderboardType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "Leaderboard identifier",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			userID := p.Context.Value("user_id").(int)
			id, _ := p.Args["id"].(int)

			leaderboard := lf.leaderboardService.GetOne(id)
			if leaderboard.Id == 0 || leaderboard.Creator.Id != userID {
				return nil, fmt.Errorf("you are not authorized for this or the resource does not exist")
			}
			lf.leaderboardService.Delete(leaderboard)
			return leaderboard, nil
		},
	}
}


func (lf *LeaderboardField) getCurrentUser(id int) *models.User {
	userService := &services.UserService{
		DB: lf.DB,
	}
	currentUser := userService.GetOne(id)
	return &currentUser
}
