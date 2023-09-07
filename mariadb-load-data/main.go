package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/montanaflynn/stats"
	"golang.org/x/text/message"
)

func assert(b bool) {
	if !b {
		panic("")
	}
}

const ROWS = 500_000
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
		log.Fatal(err)
	}

	ddl := "CREATE TABLE " + TABLE + " (\n  "
	for i, column := range COLUMNS {
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
	query := "LOAD DATA LOCAL INFILE '" + dataFile + "' INTO TABLE " + TABLE
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

func generateData(n int, outFile string) {
	f, err := os.Open("/dev/random")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	outRaw, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer outRaw.Sync()
	defer outRaw.Close()

	out := bufio.NewWriter(outRaw)
	defer out.Flush()

	totalBytes := COLUMN_SIZE * len(COLUMNS) * n
	needed := make([]byte, 0, totalBytes)

	var buf = make([]byte, 4096)
	for len(needed) == totalBytes {
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

	for i := 0; i < n; i++ {
		rowBase := i * COLUMN_SIZE * len(COLUMNS)
		for j := 0; j < len(COLUMNS); j++ {
			if j > 0 {
				writeAll(out, []byte(","))
			}
			cell := needed[rowBase+j*COLUMN_SIZE : rowBase+(j+1)*COLUMN_SIZE]
			assert(len(cell) == COLUMN_SIZE)
			writeAll(out, cell)
		}

		writeAll(out, []byte("\n"))
	}
	needed = nil
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
	generateData(ROWS, dataFile)

	var times []float64
	var throughput []float64
	for runs := 0; runs < 10; runs++ {
		fmt.Println("Preparing run", runs+1)
		prepare(db)

		fmt.Println("Executing run", runs+1)
		t1 := time.Now()
		run(db, dataFile)
		t2 := time.Now()
		diff := t2.Sub(t1).Seconds()
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

	p := message.NewPrinter(message.MatchLanguage("en"))
	p.Printf("Timing: %.2f ± %.2fs, Min: %.2fs, Max: %.2fs\n", median, stddev, min, max)
	p.Printf("Throughput: %.2f ± %.2f rows/s, Min: %.2f rows/s, Max: %.2f rows/s\n", t_median, t_stddev, t_min, t_max)
}
