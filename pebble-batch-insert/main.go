package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/cockroachdb/pebble"

	// Rewritten to ../lib
	"lib"
)

func prepare() *pebble.DB {
	err := os.RemoveAll("data.pbl")
	if err != nil {
		log.Fatal(err)
	}

	db, err := pebble.Open("data.pbl", &pebble.Options{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func run(db *pebble.DB, rows [][]any) {
	batchSize := 100_000

	var rowHolder = make([]byte, len(lib.COLUMNS)*lib.COLUMN_SIZE)
	for counter := 0; counter < len(rows); counter += batchSize {
		subset := rows[counter : min(counter+batchSize, len(rows))-1]
		lib.Assert(len(subset) <= batchSize)
		batch := db.NewBatch()

		key := make([]byte, 4)
		for i, row := range subset {
			for j, cell := range row {
				copy(rowHolder[j:j+lib.COLUMN_SIZE], cell.([]byte))
			}

			rowId := counter + i
			binary.LittleEndian.PutUint32(key, uint32(rowId))
			if err := batch.Set(key, rowHolder, nil); err != nil {
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

func main() {
	fmt.Println("Generating data")
	data := lib.GenerateData()

	var db *pebble.DB
	lib.Benchmark(func() {
		if db != nil {
			db.Close()
		}
		db = prepare()
	}, func() {
		run(db, data)
	})
}
