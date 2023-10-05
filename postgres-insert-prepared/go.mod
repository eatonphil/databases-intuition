module postgres

go 1.20

replace lib => ../lib

require (
	github.com/jackc/pgx/v5 v5.4.3
	lib v0.0.0-00010101000000-000000000000
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/text v0.13.0 // indirect
)
