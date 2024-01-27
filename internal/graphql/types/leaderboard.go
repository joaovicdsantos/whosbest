package types

import "github.com/graphql-go/graphql"

var LeaderboardType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Leaderboard",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.Int,
				Description: "Leaderboard identifier",
			},
			"title": &graphql.Field{
				Type:        graphql.String,
				Description: "Leaderboard's title",
			},
			"description": &graphql.Field{
				Type:        graphql.String,
				Description: "A description for leaderboard",
			},
			"image_url": &graphql.Field{
				Type:        graphql.String,
				Description: "A image for represent the leaderboard",
			},
			"creator": &graphql.Field{
				Type:        UserForLeaderboardType,
				Description: "Leaderboard's creator",
			},
			"competitors": &graphql.Field{
				Type:        graphql.NewList(CompetitorType),
				Description: "Leaderboard's competitors",
			},
		},
	},
)
