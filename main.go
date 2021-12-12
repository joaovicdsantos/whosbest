package main

import "github.com/joaovicdsantos/whosbest-api/app"

func main() {
	app := new(app.App)
	app.Initialize()
	app.Run()
}
