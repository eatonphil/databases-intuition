package lib

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/montanaflynn/stats"
	"golang.org/x/text/message"
)

func Assert(b bool) {
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
	"m15",
	"m16",
}

const COLUMN_SIZE = 32

var BYTES_PER_ROW = float64(len(COLUMNS) * COLUMN_SIZE)

func PrepareSQL(db *sql.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS " + TABLE)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
}

var private_INSERT_SQL_PREPARED = ""

func RunSQLInsertPrepared(db *sql.DB, rows [][]any, makePlaceholder func(i int) string) {
	if private_INSERT_SQL_PREPARED == "" {
		private_INSERT_SQL_PREPARED = "INSERT INTO " + TABLE + " VALUES (\n  "
		for i := range COLUMNS {
			if i > 0 {
				private_INSERT_SQL_PREPARED += ",\n  "
			}
			placeholder := "?"
			if makePlaceholder != nil {
				placeholder = makePlaceholder(i)
			}
			private_INSERT_SQL_PREPARED += placeholder
		}
		private_INSERT_SQL_PREPARED += "\n)"
	}

	RunSQL(db, rows, private_INSERT_SQL_PREPARED)
}

func RunSQL(db *sql.DB, rows [][]any, query string) {
	Assert(query != "")

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}

	for i, row := range rows {
		Assert(len(row) == len(COLUMNS))
		_, err = stmt.Exec(row...)
		if err != nil {
			log.Fatalf("Failed to add row %d: %s", i, err)
		}
	}

	// This is a PostgreSQL COPY-specific thing.
	if strings.HasPrefix(query, "COPY ") {
		_, err = stmt.Exec()
		if err != nil {
			log.Fatalf("Failed to exec: %s", err)
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

func GenerateData() [][]any {
	f, err := os.Open("/dev/random")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	totalBytes := COLUMN_SIZE * len(COLUMNS) * ROWS
	needed := make([]byte, 0, totalBytes)

	var buf = make([]byte, 4096)
	for len(needed) != totalBytes {
		Assert(len(needed) <= totalBytes)

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

	data := make([][]any, ROWS)
	for i := 0; i < ROWS; i++ {
		rowBase := i * COLUMN_SIZE * len(COLUMNS)
		row := make([]any, len(COLUMNS))
		for j := 0; j < len(COLUMNS); j++ {
			cell := needed[rowBase+j*COLUMN_SIZE : rowBase+(j+1)*COLUMN_SIZE]
			row[j] = cell
			Assert(len(cell) == COLUMN_SIZE)
		}
		data[i] = row
	}

	needed = nil
	return data
}

func Benchmark(prepare func(), run func()) {
	p := message.NewPrinter(message.MatchLanguage("en"))

	var times []float64
	var throughput []float64
	for runs := 0; runs < 10; runs++ {
		fmt.Println("Preparing run", runs+1)
		prepare()

		fmt.Println("Executing run", runs+1)

		t1 := time.Now()
		run()
		t2 := time.Now()

		diff := t2.Sub(t1).Seconds()
		p.Printf("Completed run in %.2fs\n", diff)
		times = append(times, diff)
		throughput = append(throughput, float64(ROWS)/diff)
	}

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

	p.Printf("Rows: %d, Columns: %d, Column Size: %d\n", ROWS, len(COLUMNS), COLUMN_SIZE)
	p.Printf("Timing: %.2f ± %.2fs, Min: %.2fs, Max: %.2fs\n", median, stddev, min, max)
	p.Printf("Throughput: %.2f ± %.2f rows/s, Min: %.2f rows/s, Max: %.2f rows/s\n", t_median, t_stddev, t_min, t_max)
}
