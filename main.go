package main

import (
	"log"

	"github.com/joaovicdsantos/whosbest-api/app"
	"github.com/joaovicdsantos/whosbest-api/config"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbFactory := new(config.DBFactory)
	db := dbFactory.InitDatabase()
	defer db.Close()

	app := new(app.App)
	app.Initialize(db)
	app.Run()
}
