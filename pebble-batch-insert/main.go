package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cockroachdb/pebble"
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

func prepare() *pebble.DB {
	err := os.RemoveAll("data.pbl")
	if err != nil {
		log.Fatal(err)
	}

	db, err := pebble.Open("data.pbl", &pebble.Options{})
	if err != nil {
		log.Fatal(err)
	}

	start := make([]byte, 4)
	end := make([]byte, 4)
	binary.LittleEndian.PutUint32(start, uint32(0))
	binary.LittleEndian.PutUint32(end, uint32(ROWS-1))
	err = db.DeleteRange(start, end, nil)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func run(db *pebble.DB, rows [][]byte) {
	batchSize := 100_000

	for counter := 0; counter < len(rows); counter += batchSize {
		subset := rows[counter : counter+batchSize]
		assert(len(subset) == batchSize)
		batch := db.NewBatch()

		key := make([]byte, 4)
		for i, row := range subset {
			rowId := counter + i
			binary.LittleEndian.PutUint32(key, uint32(rowId))
			if err := batch.Set(key, row, nil); err != nil {
				log.Fatal(err)
			}
		}

		err := batch.Commit(nil)
		if err != nil {
			log.Fatal(err)
		}

		err = batch.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func generateData(n int) [][]byte {
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

	data := make([][]byte, n)
	for i := 0; i < n; i++ {
		rowBase := i * COLUMN_SIZE * len(COLUMNS)
		data[i] = needed[rowBase : rowBase+COLUMN_SIZE*len(COLUMNS)]
	}

	needed = nil
	return data
}

func main() {
	fmt.Println("Generating data")
	data := generateData(ROWS)

	p := message.NewPrinter(message.MatchLanguage("en"))

	var times []float64
	var throughput []float64
	for runs := 0; runs < 10; runs++ {
		fmt.Println("Preparing run", runs+1)
		db := prepare()

		fmt.Println("Executing run", runs+1)
		t1 := time.Now()
		run(db, data)
		t2 := time.Now()
		diff := t2.Sub(t1).Seconds()
		p.Printf("Completed run in %.2fs\n", diff)
		times = append(times, diff)
		throughput = append(throughput, float64(ROWS)/diff)

		var count uint64
		iter, err := db.NewIter(nil)
		if err != nil {
			log.Fatal(err)
		}
		for iter.First(); iter.Valid(); iter.Next() {
			count++
		}
		if err := iter.Close(); err != nil {
			log.Fatal(err)
		}
		assert(count == ROWS)

		db.Close()
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

	p.Printf("Rows: %d\n", ROWS)
	p.Printf("Timing: %.2f ± %.2fs, Min: %.2fs, Max: %.2fs\n", median, stddev, min, max)
	p.Printf("Throughput: %.2f ± %.2f rows/s, Min: %.2f rows/s, Max: %.2f rows/s\n", t_median, t_stddev, t_min, t_max)
}
