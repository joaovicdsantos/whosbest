package types

import "github.com/graphql-go/graphql"

var CompetitorType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Competitor",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"image_url": &graphql.Field{
				Type: graphql.String,
			},
			"votes": &graphql.Field{
				Type: graphql.Int,
			},
			"leaderboard_id": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
