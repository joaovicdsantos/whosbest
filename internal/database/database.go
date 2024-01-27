package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type DBFactory struct {
}

func (d *DBFactory) InitDatabase() *sql.DB {

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	name := os.Getenv("POSTGRES_DB")

	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, name)

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal("Error on connect to database")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error testing the database connection", err)
	}

	return db
}
