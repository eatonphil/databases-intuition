module sqlite

go 1.21.2

replace lib => ../go-lib

replace lib/gosqlite => ../go-lib/gosqlite

require (
	github.com/eatonphil/gosqlite v0.8.0
	lib v0.0.0-00010101000000-000000000000
	lib/gosqlite v0.0.0-00010101000000-000000000000
)

require (
	github.com/montanaflynn/stats v0.7.1 // indirect
	golang.org/x/text v0.13.0 // indirect
)
