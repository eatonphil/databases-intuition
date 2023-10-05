module sqlite

go 1.21.2

replace lib => ../go-lib

// Only to force SQLite 3.43.1
replace github.com/mattn/go-sqlite3 => github.com/eatonphil/go-sqlite3 v1.14.14-0.20231005193852-b20e196f602b

require (
	github.com/mattn/go-sqlite3 v0.0.0-00010101000000-000000000000
	lib v0.0.0-00010101000000-000000000000
)

require (
	github.com/montanaflynn/stats v0.7.1 // indirect
	golang.org/x/text v0.13.0 // indirect
)
