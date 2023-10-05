package libgosqlite

import (
	"log"

	"github.com/eatonphil/gosqlite"

	// Rewritten to ../go-lib
	"lib"
)


func Prepare(db *gosqlite.Conn) {
	err := db.Exec("DROP TABLE IF EXISTS " + lib.TABLE)
	if err != nil {
		log.Fatalf("Failed to drop: %s", err)
	}

	ddl := "CREATE TABLE " + lib.TABLE + " (\n  "
	for i, column := range lib.COLUMNS {
		if i > 0 {
			ddl += ",\n  "
		}

		ddl += column + " TEXT"
	}
	ddl += ")"
	err = db.Exec(ddl)
	if err != nil {
		log.Fatalf("Failed to create table: %s", err)
	}
}

func Insert(db *gosqlite.Conn, rows [][]any) {
	err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	insert := "INSERT INTO " + lib.TABLE + " VALUES (\n  "
	for i := range lib.COLUMNS {
		if i > 0 {
			insert += ",\n  "
		}
		insert += "?"
	}
	insert += "\n)"
	stmt, err := db.Prepare(insert)
	if err != nil {
		log.Fatalf("Failed to prepare: %s", err)
	}

	for i, row := range rows {
		lib.Assert(len(row) == len(lib.COLUMNS))
		err = stmt.Exec(row...)
		if err != nil {
			log.Fatalf("Failed to copy row %d: %s", i, err)
		}
	}

	err = stmt.Close()
	if err != nil {
		log.Fatalf("Failed to close: %s", err)
	}

	err = db.Commit()
	if err != nil {
		log.Fatalf("Failed to commit: %s", err)
	}
}
