package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ptiger10/pd/benchmarking/profiler/benchmarks"
)

func main() {
	goBenchmarks := benchmarks.RunGoProfiler()
	pyBenchmarks := benchmarks.RunPythonProfiler()

	table := benchmarks.CompareBenchmarks(goBenchmarks, pyBenchmarks, benchmarks.Descriptions)
	dest := "comparison_summary.txt"
	err := ioutil.WriteFile(dest, []byte(table), 0666)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(">> %v\n", dest)

}
