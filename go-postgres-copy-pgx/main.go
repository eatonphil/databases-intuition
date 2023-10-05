package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	// Gets rewritten in go.mod to ../go-lib
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

	conn, err := db.Conn(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Raw(func(driverConn any) error {
		conn := driverConn.(*stdlib.Conn).Conn()

		lib.Benchmark(func() {
			lib.PrepareSQL(db)
		}, func() {
			conn.CopyFrom(
				context.Background(),
				pgx.Identifier{lib.TABLE},
				lib.COLUMNS,
				pgx.CopyFromRows(data),
			)
		})

		return nil
	})

	var count uint64
	err = db.QueryRow("SELECT COUNT(1) FROM " + lib.TABLE).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	lib.Assert(count == lib.ROWS)
}
