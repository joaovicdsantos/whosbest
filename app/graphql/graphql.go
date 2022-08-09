package graphql

import (
	"database/sql"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/joaovicdsantos/whosbest-api/app/graphql/fields"
)

type GraphQL struct {
	Schema graphql.Schema
}

func (g *GraphQL) Initialize(db *sql.DB) {
	var err error
	g.Schema, err = g.defineSchema(db)
	if err != nil {
		panic(err)
	}
}

func (g *GraphQL) defineSchema(db *sql.DB) (graphql.Schema, error) {
	competitorField := fields.CompetitorField{
		DB: db,
	}
	return graphql.NewSchema(
		graphql.SchemaConfig{
			Query: graphql.NewObject(
				graphql.ObjectConfig{
					Name: "Core_Query",
					Fields: graphql.Fields{
						"competitor": competitorField.GetOne(),
						"competitors": competitorField.GetAll(),
					},
				}),
			Mutation: graphql.NewObject(
				graphql.ObjectConfig{
					Name: "Core_Mutation",
					Fields: graphql.Fields{
						"createCompetitor": competitorField.Create(),
						"updateCompetitor": competitorField.Update(),
						"deleteCompetitor": competitorField.Delete(),
					},
				}),
		},
	)
}

func (g *GraphQL) ExecuteQuery(query string) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        g.Schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}
