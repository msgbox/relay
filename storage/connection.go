package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/bmizerany/pq"
	"os"
)

var DB_NAME = os.Getenv("DATABASE_NAME")
var DB_USER = os.Getenv("DATABASE_USER")
var DB_PASS = os.Getenv("DATABASE_PASS")
var DB_HOST = os.Getenv("DATABASE_HOST")
var DB_PORT = os.Getenv("DATABASE_PORT")
var DB_SSL = os.Getenv("DATABASE_SSL")

func Connect() *sql.DB {
	db, err := sql.Open("postgres", fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=%s", DB_NAME, DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_SSL))
	if err != nil {
		panic(err)
	}
	return db
}
