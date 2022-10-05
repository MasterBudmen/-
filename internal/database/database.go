package database

import (
	"database/sql"
	"fmt"
	"log"

	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDb() *sql.DB {
	//HEROKU
	//connStr := os.Getenv("DATABASE_URL")
	//LOCAL
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("host"), os.Getenv("port"), os.Getenv("user"), os.Getenv("password"), os.Getenv("dbname"))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
