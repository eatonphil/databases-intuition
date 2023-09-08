module sqlite

go 1.21.1

replace lib => ../lib

require (
	github.com/mattn/go-sqlite3 v1.14.17
	lib v0.0.0-00010101000000-000000000000
)

require (
	github.com/montanaflynn/stats v0.7.1 // indirect
	golang.org/x/text v0.13.0 // indirect
)
