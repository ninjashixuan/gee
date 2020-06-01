package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "gee.db")
	if err != nil {
		panic("sql Open failed error: " + err.Error())
	}
	defer func() {
		_ = db.Close()
	}()

	_, _ = db.Exec("DROP TABLE IF EXISTS User;")
	_, _ = db.Exec("CREATE TABLE User(Name text, Age integer);")
	result, err := db.Exec("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam")
	if err == nil {
		affected, _ := result.RowsAffected()
		log.Println(affected)
	}

	result, err = db.Exec("INSERT INTO User(`Name`, `Age`) values (?, ?), (?, ?)", "ninja", 12, "hello",14)
	if err == nil {
		affected, _ := result.RowsAffected()
		log.Println(affected)
	}

	row := db.QueryRow("SELECT Name FROM User LIMIT 1")
	var name string
	if err := row.Scan(&name); err == nil {
		log.Println(name)
	}
}
