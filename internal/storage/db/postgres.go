package db

import "database/sql"

var Db *sql.DB

func InitDb() {
	Db = connectToDb()
}

func connectToDb() *sql.DB {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/newsletter?sslmode=disable")
	if err != nil {
		panic(err)
	}
	return db
}
