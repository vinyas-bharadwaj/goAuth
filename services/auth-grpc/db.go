package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)


func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		panic(err)
	}

	query := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT UNIQUE,
		password TEXT
	);`

	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	return db
}