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
	var version string
	err = db.QueryRow("SELECT sqlite_version()").Scan(&version) // Test the connection.
	if err != nil {
		log.Fatal(err)
	}
	lib.Assert(version == "3.43.1")

	fmt.Println("Generating data")
	data := lib.GenerateData()

	lib.Benchmark(func() {
		lib.PrepareSQL(db)
	}, func() {
		lib.RunSQLInsertPrepared(db, data, nil)
	})
}
