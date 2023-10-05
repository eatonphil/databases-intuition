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

	lib.Benchmark(func() {
		libgosqlite.Prepare(db)
	}, func() {
		libgosqlite.Insert(db, data)
	})
}
