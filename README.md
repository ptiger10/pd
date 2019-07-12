# pd
[![Go Report Card](https://goreportcard.com/badge/github.com/ptiger10/pd)](https://goreportcard.com/report/github.com/ptiger10/pd) 
[![GoDoc](https://godoc.org/github.com/ptiger10/pd?status.svg)](https://godoc.org/github.com/ptiger10/pd) 
[![Build Status](https://travis-ci.org/ptiger10/pd.svg?branch=master)](https://travis-ci.org/ptiger10/pd)
[![codecov](https://codecov.io/gh/ptiger10/pd/branch/master/graph/badge.svg)](https://codecov.io/gh/ptiger10/pd)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

pd (informally known as "GoPandas") is a library for cleaning, aggregating, and transforming data using Series and DataFrames. GoPandas combines a flexible API familiar to Python pandas users with the qualities of Go, including type safety, predictable error handling, and fast concurrent processing.

The API is still version 0 and subject to major revisions. Use in production code at your own risk.

Some notable features of GoPandas:
* flexible constructor that supports float, int, string, bool, time.Time, and interface Series
* seamlessly handles null data and type conversions
* well-suited to either the Jupyter notebook style of data exploration or conventional programming
* advanced filtering, grouping, and pivoting
* hierarchical indexing (i.e., multi-level indexes and columns)
* reads from either CSV or any spreadsheet or tabular data structured as [][]interface (e.g., Google Sheets)
* complete test coverage
* minimal dependencies (total package size is <10MB, compared to Pandas at >200MB)
* uses concurrent processing to achieve faster speeds than Pandas on many fundamental operations, and the performance differential becomes more pronounced with scale (6x+ superior performance summing two columns in a 500k row spreadsheet - see the most recent [benchmarking table](benchmarking/profiler/comparison_summary.txt)

## Getting Started
Check out the Jupyter notebook examples in the [guides](https://github.com/ptiger10/pd/tree/master/guides). Github sometimes has trouble rendering .ipynb, backup views are here: [Series](https://nbviewer.jupyter.org/github/ptiger10/pd/blob/master/guides/Series.ipynb?flush_cache=true), [DataFrame](https://nbviewer.jupyter.org/github/ptiger10/pd/blob/master/guides/DataFrame.ipynb?flush_cache=true), [Options](https://nbviewer.jupyter.org/github/ptiger10/pd/blob/master/guides/Options.ipynb?flush_cache=true).

To run the Jupyter notebooks yourself, I recommend lgo (Docker required)
* `cd guides/docker`
* start: `./up.sh`
* stop: `./down.sh`
* rebuild package to newest version: `./up.sh -r`

## Replicating Benchmark Tests
* Requires Python 3.x and pandas
* Download data from [here](https://github.com/ptiger10/pdTestData) and save in benchmarking/profiler
* `go run -tags=benchmarks benchmarking/profiler/main.go`