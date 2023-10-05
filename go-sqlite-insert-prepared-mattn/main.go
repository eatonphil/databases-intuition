package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	// Rewritten to ../go-lib
	"lib"
)

func main() {
	db, err := sql.Open("sqlite3", "data.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("SELECT 1") // Test the connection.
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Generating data")
	data := lib.GenerateData()

	lib.Benchmark(func() {
		lib.PrepareSQL(db)
	}, func() {
		lib.RunSQLInsertPrepared(db, data, nil)
	})
}
