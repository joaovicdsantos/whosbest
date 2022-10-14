package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/joaovicdsantos/whosbest-api/app/graphql"
	"github.com/joaovicdsantos/whosbest-api/app/helpers"
	"github.com/joaovicdsantos/whosbest-api/app/models"
)

type GraphqlRoute struct {
	DB *sql.DB
}

func (gr *GraphqlRoute) HandleGraphqlRequest(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		response := helpers.NewResponseError("Method Not Allowed", http.StatusMethodNotAllowed)
		response.SendResponse(w)
		return
	}

	var graphqlIn models.GraphQL
	if err := helpers.ParseBodyToStruct(r, &graphqlIn); err != nil {
		response := helpers.NewResponseError(err.Error(), http.StatusBadRequest)
		response.SendResponse(w)
		return
	}

	// Token Validation
	unparsedToken, ok := r.Header["Authorization"]
	if !ok {
		response := helpers.NewResponseError("Token is not valid", http.StatusUnauthorized)
		response.SendResponse(w)
		return
	}

	var err error
	payload, err := helpers.ParseJwtToken(fmt.Sprint(unparsedToken[0]))
	if err != nil {
		response := helpers.NewResponseError(err.Error(), http.StatusUnauthorized)
		response.SendResponse(w)
		return
	}

	user := helpers.GetCurrentUser(payload, gr.DB)
	if user.Id == 0 {
		response := helpers.NewResponseError("Invalid user", http.StatusUnauthorized)
		response.SendResponse(w)
		return
	}

	var graphql = new(graphql.GraphQL)
	graphql.Initialize(gr.DB)

	result := graphql.ExecuteQuery(graphqlIn.Query, user.Id)
	if result.HasErrors() {
		response := helpers.NewResponseError(result.Errors, http.StatusBadRequest)
		response.SendResponse(w)
		return
	}

	response := helpers.NewResponse(result.Data, http.StatusOK)
	response.SendResponse(w)
}
