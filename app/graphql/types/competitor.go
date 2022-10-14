package types

import "github.com/graphql-go/graphql"

var CompetitorType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Competitor",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.Int,
				Description: "Competitor identifier",
			},
			"title": &graphql.Field{
				Type:        graphql.String,
				Description: "Competitor's name or title",
			},
			"description": &graphql.Field{
				Type:        graphql.String,
				Description: "A description for competitor attributes",
			},
			"image_url": &graphql.Field{
				Type:        graphql.String,
				Description: "A image for represent the competitor",
			},
			"votes": &graphql.Field{
				Type:        graphql.Int,
				Description: "Number of votes",
			},
			"leaderboard": &graphql.Field{
				Type:        graphql.Int,
				Description: "Competitor's leaderboard",
			},
		},
	},
)
