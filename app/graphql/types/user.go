package types

import "github.com/graphql-go/graphql"

var UserType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.Int,
				Description: "User identifier",
			},
			"username": &graphql.Field{
				Type:        graphql.String,
				Description: "Username to login",
			},
		},
	},
)
