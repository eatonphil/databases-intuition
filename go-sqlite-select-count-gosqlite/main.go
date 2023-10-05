package main

import (
	"fmt"
	"log"

	"github.com/eatonphil/gosqlite"

	// Rewritten to ../go-lib
	"lib"
	libgosqlite "lib/gosqlite"
)

func main() {
	db, err := gosqlite.Open("data.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Exec("SELECT 1") // Test the connection.
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Generating data")
	data := lib.GenerateData()
	libgosqlite.Prepare(db)
	libgosqlite.Insert(db, data)

	lib.Benchmark(func() {
	}, func() {
		var n int64
		stmt, err := db.Prepare("SELECT COUNT(1) FROM " + lib.TABLE + " WHERE a1 <> b2")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		for {
			hasRow, err := stmt.Step()
			if err != nil {
				log.Fatal(err)
			}

			if !hasRow {
				break
			}

			err = stmt.Scan(&n)
			if err != nil {
				log.Fatal(err)
			}

			break
		}

		lib.Assert(n == lib.ROWS)
	})
}
