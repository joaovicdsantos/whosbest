package handler

import (
	"database/sql"
	"net/http"

	"github.com/joaovicdsantos/whosbest-api/app/graphql"
	"github.com/joaovicdsantos/whosbest-api/app/helpers"
	"github.com/joaovicdsantos/whosbest-api/app/models"
)

type GraphqlRoute struct {
	DB *sql.DB
}


func (gr *GraphqlRoute) HandleGraphqlRequest(w http.ResponseWriter, r *http.Request) {
	var graphqlIn models.GraphQL
	if err := helpers.ParseBodyToStruct(r, &graphqlIn); err != nil {
		response := helpers.NewResponseError(err.Error(), http.StatusBadRequest)
		response.SendResponse(w)
		return
	}

	var graphql = new(graphql.GraphQL)
	graphql.Initialize(gr.DB)

	result := graphql.ExecuteQuery(graphqlIn.Query)
	if result.HasErrors() {
		response := helpers.NewResponseError(result.Errors, http.StatusBadRequest)
		response.SendResponse(w)
		return
	}

	response := helpers.NewResponse(result.Data, http.StatusOK)
	response.SendResponse(w)
}
