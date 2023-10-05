package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx"

	// Gets rewritten to ../lib
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
		lib.RunSQLInsertPrepared(db, data, func(i int) string {
			return fmt.Sprintf("$%d", i+1)
		})
	})
}
