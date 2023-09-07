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

## Methodology

For each program, the main part of the program is run 10 times and we
collect median, standard deviation, min and max of the 10 runs.

For example, if a program generates 10M rows, the generation is not
part of the time measurement. Only insertion of generated data is.

# Inserts

For all variations we're:

* Loading 10M rows with 16 32-byte columns per row
* No indexes (keeps things simple for the database)
* No disabling fsync or other anti-durability tricks

Caveats to think about:

* Throughput under sustained load (i.e. not just 10M rows once, but
  10M new rows loaded 100 times) may be different/worse

## MariaDB LOAD DATA LOCAL

Uses a single `LOAD DATA LOCAL` query.

```
Timing: 77.96 ± 1.43s, Min: 75.69s, Max: 81.30s
Throughput: 128,266.45 ± 2,324.53 rows/s, Min: 123,004.14 rows/s, Max: 132,125.36 rows/s
```

## PostgreSQL COPY FROM

```
Timing: 104.53 ± 2.40s, Min: 102.57s, Max: 110.08s
Throughput: 95,665.37 ± 2,129.25 rows/s, Min: 90,847.08 rows/s, Max: 97,490.96 rows/s
```

## SQLite Parameterized INSERT

```
Timing: 52.67 ± 1.70s, Min: 49.91s, Max: 55.46s
Throughput: 189,862.60 ± 6,175.26 rows/s, Min: 180,316.37 rows/s, Max: 200,346.56 rows/s
```
