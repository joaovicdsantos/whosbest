package fields

import (
	"database/sql"

	"github.com/graphql-go/graphql"
	"github.com/joaovicdsantos/whosbest-api/app/graphql/types"
	"github.com/joaovicdsantos/whosbest-api/app/services"
)

func CompetitorField(db *sql.DB) *graphql.Field {
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
				DB: db,
			}
			if p.Args["id"] != 0 {
				id, _ := p.Args["id"].(int)
				return competitorService.GetOne(id), nil
			}
			competitors, _ := competitorService.GetAll()
			return competitors, nil
		},
	}
}
