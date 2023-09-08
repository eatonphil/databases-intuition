package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	// Rewritten to ../lib
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
	lib.PrepareSQL(db)
	lib.RunSQLInsertPrepared(db, data, nil)

	lib.Benchmark(func() {
	}, func() {
		var n int64
		err := db.QueryRow("SELECT COUNT(1) FROM " + lib.TABLE + " WHERE a1 <> b2").Scan(&n)
		if err != nil {
			log.Fatal(err)
		}

		lib.Assert(n == lib.ROWS)
	})
}
