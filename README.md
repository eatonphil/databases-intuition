# Databases Intuition

These are a series of programs to get a ballpark idea of time it takes
to do various common database operations.

The point **is not** to do benchmark wars. There is **very limited
value** in comparing these results across database.

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

# SELECTs

For all variations:

* No indexes, all values in all fields are random

The query run (or approximated in the Pebble case) is:

```sql
SELECT COUNT(1) FROM testtable1 WHERE a1 <> b2
```

Forcing a full-table scan.

## MariaDB

[Source](./mariadb-select-count)

10M Rows, 16 columns, each column 32 bytes:

```
Timing: 5.31 ± 4.00s, Min: 5.28s, Max: 15.73s
Throughput: 1,883,869.01 ± 514,123.28 rows/s, Min: 635,838.73 rows/s, Max: 1,894,079.77 rows/s
```

10M Rows, 3 columns, each column 8 bytes:

```
Timing: 3.73 ± 0.49s, Min: 3.71s, Max: 5.05s
Throughput: 2,677,891.78 ± 273,913.40 rows/s, Min: 1,980,034.10 rows/s, Max: 2,695,158.11 rows/s
```

## PostgreSQL

[Source](./postgres-select-count)

### 10M Rows, 16 columns, each column 32 bytes

lib/pq:

```
Timing: 0.74 ± 0.29s, Min: 0.69s, Max: 1.69s
Throughput: 13,579,422.60 ± 2,537,298.51 rows/s, Min: 5,905,851.18 rows/s, Max: 14,413,751.50 rows/s
```

jackc/pgx:

```
Timing: 0.76 ± 0.36s, Min: 0.72s, Max: 1.95s
Throughput: 13,219,135.30 ± 2,598,054.57 rows/s, Min: 5,127,178.31 rows/s, Max: 13,860,787.47 rows/s
```

### 10M Rows, 3 columns, each column 8 bytes

lib/pq

```
Timing: 0.33 ± 0.04s, Min: 0.32s, Max: 0.46s
Throughput: 30,613,018.35 ± 2,717,800.48 rows/s, Min: 21,728,298.74 rows/s, Max: 31,416,940.39 rows/s
```

jackc/pgx:

```
Timing: 0.32 ± 0.04s, Min: 0.32s, Max: 0.44s
Throughput: 31,078,449.79 ± 2,483,125.81 rows/s, Min: 22,796,606.48 rows/s, Max: 31,118,025.44 rows/s
```

## SQLite

[Source](./sqlite-select-count)

10M Rows, 16 columns, each column 32 bytes:

```
Timing: 2.14 ± 0.02s, Min: 2.10s, Max: 2.16s
Throughput: 4,664,568.46 ± 39,568.72 rows/s, Min: 4,639,582.64 rows/s, Max: 4,767,690.39 rows/s
```

10M Rows, 3 columns, each column 8 bytes:

```
Timing: 0.65 ± 0.01s, Min: 0.65s, Max: 0.66s
Throughput: 15,469,367.49 ± 123,336.16 rows/s, Min: 15,065,057.65 rows/s, Max: 15,492,194.30 rows/s
```

# INSERTs

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

10M Rows, 16 columns, each column 32 bytes:

```
Timing: 77.96 ± 1.43s, Min: 75.69s, Max: 81.30s
Throughput: 128,266.45 ± 2,324.53 rows/s, Min: 123,004.14 rows/s, Max: 132,125.36 rows/s
```

10M Rows, 3 columns, each column 8 bytes:

```
Timing: 23.33 ± 0.27s, Min: 22.71s, Max: 23.87s
Throughput: 428,680.78 ± 5,069.67 rows/s, Min: 418,996.47 rows/s, Max: 440,284.99 rows/s
```

## PostgreSQL COPY FROM

[Source](./postgres-copy)

Uses a single `COPY` query.

### 10M Rows, 16 columns, each column 32 bytes

lib/pq:

```
Timing: 104.53 ± 2.40s, Min: 102.57s, Max: 110.08s
Throughput: 95,665.37 ± 2,129.25 rows/s, Min: 90,847.08 rows/s, Max: 97,490.96 rows/s
```

jackc/pgx:

```
Timing: 46.54 ± 1.60s, Min: 44.09s, Max: 49.51s
Throughput: 214,869.42 ± 7,265.10 rows/s, Min: 201,991.37 rows/s, Max: 226,801.07 rows/s
```

### 10M Rows, 3 columns, each column 8 bytes:

lib/pq:

```
Timing: 8.16 ± 0.43s, Min: 7.44s, Max: 8.80s
Throughput: 1,225,986.47 ± 66,631.53 rows/s, Min: 1,136,581.82 rows/s, Max: 1,343,441.37 rows/s
```

jackc/pgx:

```
Timing: 5.20 ± 0.44s, Min: 4.71s, Max: 5.96s
Throughput: 1,923,722.79 ± 156,820.46 rows/s, Min: 1,676,894.32 rows/s, Max: 2,124,966.60 rows/s
```

## SQLite Parameterized INSERT

[Source](./sqlite-insert-prepared)

Parameterizes an `INSERT` query and calls the prepared statement for
each row.

10M Rows, 16 columns, each column 32 bytes:

```
Timing: 52.67 ± 1.70s, Min: 49.91s, Max: 55.46s
Throughput: 189,862.60 ± 6,175.26 rows/s, Min: 180,316.37 rows/s, Max: 200,346.56 rows/s
```

10M Rows, 3 columns, each column 8 bytes:

```
Timing: 16.03 ± 0.27s, Min: 15.69s, Max: 16.74s
Throughput: 623,745.44 ± 10,296.45 rows/s, Min: 597,525.44 rows/s, Max: 637,489.51 rows/s
```

Note: You can get better performance by using
[bvinc/go-sqlite-lite](https://github.com/bvinc/go-sqlite-lite) or [my
fork of it](https://github.com/eatonphil/gosqlite). For one, they
don't go through `database/sql`. But there's probably more to it than
that.

These are not popular libraries though so using them here felt instead
of mattn/go-sqlite feels a little misleading.

## Pebble Batch Insert

[Source](./pebble-batch-insert)

Splits up the inserts into batches of 100,000 (Pebble has a max
batch size of 4GB).

10M Rows, 16 columns, each column 32 bytes:

```
Timing: 82.97 ± 2.99s, Min: 78.17s, Max: 87.81s
Throughput: 120,524.57 ± 4,365.51 rows/s, Min: 113,883.97 rows/s, Max: 127,918.42 rows/s
```

10M Rows, 3 columns, each column 8 bytes:

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
gets:

```
Timing: 159.10 ± 1.78s, Min: 156.84s, Max: 163.68s
Throughput: 628,553.38 ± 6,936.19 rows/s, Min: 610,943.23 rows/s, Max: 637,597.30 rows/s
```
