package main

import (
	"database/sql"
	"log"
	"os"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("mysql", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Panic(err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}
