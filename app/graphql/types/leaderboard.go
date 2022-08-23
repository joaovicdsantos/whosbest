package types

import "github.com/graphql-go/graphql"


var LeaderboardType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Leaderboard",
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
			"creator": &graphql.Field{
				Type: UserType,
			},
			"competitors": &graphql.Field{
				Type: graphql.NewList(CompetitorType),
			},
		},
	},
)
