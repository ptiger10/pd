# pd
pd (informally known as "GoPandas") is a library for cleaning, aggregating, and transforming data using Series and DataFrames. GoPandas combines a flexible API familiar to Python pandas users with the strengths of Go, including type safety, predictable error handling, and concurrent processing.

Some notable features of GoPandas:
* flexible constructor that supports float, int, string, bool, time.Time, and interface Series
* well-suited to either the Jupyter notebook style of data exploration or conventional programming
* grouping and pivoting
* hierarchical indexing (e.g., multi-level indexes and columns)
* reads from either CSV or any spreadsheet or tabular data structured as [][]interface (e.g., Google Sheets)
* complete test coverage
* uses concurrent processing to achieve faster speeds than Pandas on many fundamental operations, and the performance differentail becomes more pronounced with scale (6x+ superior performance summing two columns in a 500k row spreadsheet - see the most recent [benchmarking table](benchmarking/profiler/comparison_summary.txt)

