// +build benchmarks

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime"

	"github.com/ptiger10/pd/benchmarking/profiler/benchmarks"
)

func main() {
	benchmarks.ReadData()
	goBenchmarks := benchmarks.RunGoProfiler()
	pyBenchmarks := benchmarks.RunPythonProfiler()

	// fmt.Println(goBenchmarks, pyBenchmarks)

	table := benchmarks.CompareBenchmarks(
		goBenchmarks, pyBenchmarks,
		benchmarks.SampleSizes, benchmarks.Descriptions)
	_, thisFile, _, _ := runtime.Caller(0)
	basename := "comparison_summary.txt"
	dest := filepath.Join(filepath.Dir(thisFile), basename)
	err := ioutil.WriteFile(dest, []byte(table), 0666)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(">> %v\n", basename)
}
