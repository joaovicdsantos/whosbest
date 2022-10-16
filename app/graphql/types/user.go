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
			"leaderboards": &graphql.Field{
				Type:        graphql.NewList(LeaderboardType),
				Description: "User's leaderboards",
			},
		},
	},
)

var UserForLeaderboardType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "UserForLeaderboard",
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
