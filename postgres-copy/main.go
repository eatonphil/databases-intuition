package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"

	// Gets rewritten in go.mod to ../lib/
	"lib"
)

func main() {
	db, err := sql.Open("postgres", "user=pgtest dbname=pgtest password=pgtest sslmode=disable host=127.0.0.1")
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
		lib.RunSQL(db, data, pq.CopyIn(lib.TABLE, lib.COLUMNS...))
	})

	var count uint64
	err = db.QueryRow("SELECT COUNT(1) FROM " + lib.TABLE).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	lib.Assert(count == lib.ROWS)
}
