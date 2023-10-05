package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"

	// Rewritten to ../lib
	"lib"
)

func main() {
	db, err := sql.Open("pgx", "user=pgtest dbname=pgtest password=pgtest sslmode=disable host=127.0.0.1")
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
	lib.RunSQLInsertPrepared(db, data, func(i int) string {
		return fmt.Sprintf("$%d", i+1)
	})

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
