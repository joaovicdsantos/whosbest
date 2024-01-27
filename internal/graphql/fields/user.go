package fields

import (
	"database/sql"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/joaovicdsantos/whosbest-api/internal/graphql/types"
	"github.com/joaovicdsantos/whosbest-api/internal/services"
)

type UserField struct {
	DB          *sql.DB
	userService *services.UserService
}

func NewUserField(db *sql.DB) *UserField {
	return &UserField{
		DB: db,
		userService: &services.UserService{
			DB: db,
		},
	}
}

func (uf *UserField) GetOne() *graphql.Field {
	return &graphql.Field{
		Type: types.UserType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.Int),
				Description: "User identifier",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id, _ := p.Args["id"].(int)
			user := uf.userService.GetOne(id)
			if user.Id == 0 {
				return nil, fmt.Errorf("user not found")
			}
			return user, nil
		},
	}
}

func (uf *UserField) GetAll() *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(types.UserType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return uf.userService.GetAll()
		},
	}
}
