# Databases Intuition

These are a series of programs to get a vague idea of what's
possible. The point **is not** to do benchmark wars.

There is **very limited value** in comparing these results across
database.

If there are obvious ways a program can do better though, please open
an issue!

# Machine

I am running these tests on a dedicated bare metal instance, [OVH
Rise-1](https://eco.us.ovhcloud.com/#filterType=range_element&filterValue=rise).

* RAM: 64 GB DDR4 ECC 2,133 MHz
* Disk: 2x450 GB SSD NVMe in Soft RAID
* Processor: Intel Xeon E3-1230v6 - 4c/8t - 3.5 GHz/3.9 GHz
* `uname --kernel-release`: 6.3.8-100.fc37.x86_64

## Databases

* SQLite: 3.40.1
* PostgreSQL: 15.1
* MariaDB: 10.5.21
* Pebble: v0.0.0-20230907151110-6f6852d52f18

## Methodology

For each program, the main part of the program is run 10 times and we
collect median, standard deviation, min and max of the 10 runs.

For example, if a program generates 10M rows, the generation is not
part of the time measurement. Only insertion of generated data is.

# Inserts

For all variations:

* No indexes (keeps things simple for the database)
* No disabling fsync or other anti-durability tricks

And we insert 10M rows of two different sizes:

* 16 columns, 32 bytes each
* 3 columns, 8 bytes each

Caveats to think about:

* Throughput under sustained load (i.e. not just 10M rows once, but
  10M new rows loaded 100 times) may be different/worse
* PostgreSQL doesn't have a page cache so it's likely incurring many
  more `write` syscalls for its B-Tree than MariaDB or SQLite

## MariaDB LOAD DATA LOCAL

[Source](./mariadb-load-data)

Uses a single `LOAD DATA LOCAL` query.

100M Rows, 16 columns, each column 32 bytes:

```
Timing: 77.96 ± 1.43s, Min: 75.69s, Max: 81.30s
Throughput: 128,266.45 ± 2,324.53 rows/s, Min: 123,004.14 rows/s, Max: 132,125.36 rows/s
```

100M Rows, 3 columns, each column 8 bytes:

```
Timing: 23.33 ± 0.27s, Min: 22.71s, Max: 23.87s
Throughput: 428,680.78 ± 5,069.67 rows/s, Min: 418,996.47 rows/s, Max: 440,284.99 rows/s
```

## PostgreSQL COPY FROM

[Source](./postgres-copy)

Uses a single `COPY` query.

100M Rows, 16 columns, each column 32 bytes:

```
Timing: 104.53 ± 2.40s, Min: 102.57s, Max: 110.08s
Throughput: 95,665.37 ± 2,129.25 rows/s, Min: 90,847.08 rows/s, Max: 97,490.96 rows/s
```

100M Rows, 3 columns, each column 8 bytes:

```
Timing: 8.16 ± 0.43s, Min: 7.44s, Max: 8.80s
Throughput: 1,225,986.47 ± 66,631.53 rows/s, Min: 1,136,581.82 rows/s, Max: 1,343,441.37 rows/s
```

## SQLite Parameterized INSERT

[Source](./sqlite-insert-prepared)

Parameterizes an `INSERT` query and calls the prepared statement for
each row.

100M Rows, 16 columns, each column 32 bytes:

```
Timing: 52.67 ± 1.70s, Min: 49.91s, Max: 55.46s
Throughput: 189,862.60 ± 6,175.26 rows/s, Min: 180,316.37 rows/s, Max: 200,346.56 rows/s
```

100M Rows, 3 columns, each column 8 bytes:

```
Timing: 16.03 ± 0.27s, Min: 15.69s, Max: 16.74s
Throughput: 623,745.44 ± 10,296.45 rows/s, Min: 597,525.44 rows/s, Max: 637,489.51 rows/s
```

## Pebble Batch Insert

[Source](./pebble-batch-insert)

Splits up the inserts into batches of 100,000 (Pebble has a max
batch size of 4GB).

100M Rows, 16 columns, each column 32 bytes:

```
Timing: 82.97 ± 2.99s, Min: 78.17s, Max: 87.81s
Throughput: 120,524.57 ± 4,365.51 rows/s, Min: 113,883.97 rows/s, Max: 127,918.42 rows/s
```

100M Rows, 3 columns, each column 8 bytes:

```
Timing: 15.32 ± 0.70s, Min: 13.82s, Max: 15.93s
Throughput: 652,938.75 ± 31,429.74 rows/s, Min: 627,658.58 rows/s, Max: 723,686.13 rows/s
```

# Similar work

## Towards Inserting One Billion Rows in SQLite Under A Minute

https://avi.im/blag/2021/fast-sqlite-inserts/

This study gets about 100M rows of 3 columns (all <= 8 bytes wide)
into SQLite in under 30s using Rust.

It makes a number of concessions to ACID-compliance that you wouldn't
actually want to use. So it's value is somewhat limited.

That said, it is solid work since even with all the concessions made
here (and the change to 3 columns of 8 bytes each), the best I could
get the SQLite insert to do was 100M in ~180s rather than <30s Avi
gets.
