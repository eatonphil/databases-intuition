package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"

	// Rewritte to ../go-lib
	"lib"
)

func prepare(db *sql.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS " + lib.TABLE)
	if err != nil {
		log.Fatal(err)
	}

	ddl := "CREATE TABLE " + lib.TABLE + " (\n  "
	for i, column := range lib.COLUMNS {
		if i > 0 {
			ddl += ",\n  "
		}

		ddl += column + " TEXT"
	}
	ddl += ")" // + " WITH (autovacuum_enabled = false)"
	_, err = db.Exec(ddl)
	if err != nil {
		log.Fatal(err)
	}
}

func run(db *sql.DB, dataFile string) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	mysql.RegisterLocalFile(dataFile)
	query := "LOAD DATA LOCAL INFILE '" + dataFile + "' INTO TABLE " + lib.TABLE
	_, err = tx.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Failed to commit: %s", err)
	}
}

func writeAll(out io.Writer, bytes []byte) {
	written := 0
	for written < len(bytes) {
		n, err := out.Write(bytes[written:])
		if err != nil {
			log.Fatal(err)
		}

		written += n
	}
}

func generateData(outFile string) {
	outRaw, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer outRaw.Sync()
	defer outRaw.Close()

	out := bufio.NewWriter(outRaw)
	defer out.Flush()

	rows := lib.GenerateData()

	for i, row := range rows {
		if i > 0 {
			writeAll(out, []byte("\n"))
		}

		for j, cell := range row {
			if j > 0 {
				writeAll(out, []byte(","))
			}

			writeAll(out, cell.([]byte))
		}
	}
}

func main() {
	db, err := sql.Open("mysql", "mariadbtest:mariadbtest@tcp(localhost)/mariadbtest")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("SELECT 1") // Test the connection.
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Generating data")
	dataFile := "data.csv"
	generateData(dataFile)

	lib.Benchmark(func() {
		prepare(db)
	}, func() {
		run(db, dataFile)
	})
}
