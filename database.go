package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	adb, e := sql.Open("postgres", "postgresql://jairo@localhost:26257/suncork?sslmode=disable")
	db = adb
	crash(e)
}
