package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/montanaflynn/stats"
	"golang.org/x/text/message"
)

func assert(b bool) {
	if !b {
		panic("")
	}
}

const ROWS = 10_000_000
const TABLE = "testtable1"

var COLUMNS = []string{
	"a1",
	"b2",
	"c3",
	"d4",
	"e5",
	"f6",
	"g7",
	"h8",
	"g9",
	"h10",
	"i11",
	"j12",
	"k13",
	"l14",
	"m14",
}

const COLUMN_SIZE = 32

func prepare(db *sql.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS " + TABLE)
	if err != nil {
		log.Fatalf("Failed to drop: %s", err)
	}

	ddl := "CREATE TABLE " + TABLE + " (\n  "
	for i, column := range COLUMNS {
		if i > 0 {
			ddl += ",\n  "
		}

		ddl += column + " TEXT"
	}
	ddl += ")"
	_, err = db.Exec(ddl)
	if err != nil {
		log.Fatalf("Failed to create table: %s", err)
	}
}

func run(db *sql.DB, rows [][]any) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	insert := "INSERT INTO " + TABLE + " VALUES (\n  "
	for i := range COLUMNS {
		if i > 0 {
			insert += ",\n  "
		}
		insert += "?"
	}
	insert += "\n)"
	stmt, err := tx.Prepare(insert)
	if err != nil {
		log.Fatalf("Failed to prepare: %s", err)
	}

	for i, row := range rows {
		assert(len(row) == len(COLUMNS))
		_, err = stmt.Exec(row...)
		if err != nil {
			log.Fatalf("Failed to copy row %d: %s", i, err)
		}
	}

	err = stmt.Close()
	if err != nil {
		log.Fatalf("Failed to close: %s", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Failed to commit: %s", err)
	}
}

func generateData(n int) [][]any {
	f, err := os.Open("/dev/random")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	totalBytes := COLUMN_SIZE * len(COLUMNS) * n
	needed := make([]byte, 0, totalBytes)

	var buf = make([]byte, 4096)
	for len(needed) != totalBytes {
		assert(len(needed) <= totalBytes)

		n, err := f.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		for _, c := range buf[:n] {
			if (c >= 'a' && c <= 'z') ||
				(c >= 'A' && c <= 'Z') ||
				(c >= '0' && c <= '9') {
				needed = append(needed, c)
			}

			if len(needed) == totalBytes {
				break
			}
		}
	}

	data := make([][]any, n)
	for i := 0; i < n; i++ {
		rowBase := i * COLUMN_SIZE * len(COLUMNS)
		row := make([]any, len(COLUMNS))
		for j := 0; j < len(COLUMNS); j++ {
			cell := needed[rowBase+j*COLUMN_SIZE : rowBase+(j+1)*COLUMN_SIZE]
			row[j] = cell
			assert(len(cell) == COLUMN_SIZE)
		}
		data[i] = row
	}

	needed = nil
	return data
}

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
	data := generateData(ROWS)

	var times []float64
	var throughput []float64
	for runs := 0; runs < 10; runs++ {
		fmt.Println("Preparing run", runs+1)
		prepare(db)

		fmt.Println("Executing run", runs+1)
		t1 := time.Now()
		run(db, data)
		t2 := time.Now()
		diff := t2.Sub(t1).Seconds()
		times = append(times, diff)
		throughput = append(throughput, float64(ROWS)/diff)
	}

	var count uint64
	err = db.QueryRow("SELECT COUNT(1) FROM " + TABLE).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	assert(count == ROWS)

	median, err := stats.Median(times)
	if err != nil {
		log.Fatal(err)
	}

	min, err := stats.Min(times)
	if err != nil {
		log.Fatal(err)
	}

	max, err := stats.Max(times)
	if err != nil {
		log.Fatal(err)
	}

	stddev, err := stats.StandardDeviation(times)
	if err != nil {
		log.Fatal(err)
	}

	t_median, err := stats.Median(throughput)
	if err != nil {
		log.Fatal(err)
	}

	t_min, err := stats.Min(throughput)
	if err != nil {
		log.Fatal(err)
	}

	t_max, err := stats.Max(throughput)
	if err != nil {
		log.Fatal(err)
	}

	t_stddev, err := stats.StandardDeviation(throughput)
	if err != nil {
		log.Fatal(err)
	}

	p := message.NewPrinter(message.MatchLanguage("en"))
	p.Printf("Timing: %.2f ± %.2fs, Min: %.2fs, Max: %.2fs\n", median, stddev, min, max)
	p.Printf("Throughput: %.2f ± %.2f rows/s, Min: %.2f rows/s, Max: %.2f rows/s\n", t_median, t_stddev, t_min, t_max)
}
