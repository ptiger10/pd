package benchmarks

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path"
	"runtime"
	"testing"
	"time"
)

type desc struct {
	order int
	str   string
}

// Descriptions of the benchmarking tests
var Descriptions = map[string]desc{
	"sum": desc{1, "Sum one column"},
	// "mean":       desc{2, "Simple mean of one column"},
	// "median":     desc{3, "Median of one column"},
	// "min":        desc{4, "Min of one column"},
	// "max":        desc{5, "Max of one column"},
	// "std":        desc{6, "Standard deviation of one column"},
	// "readCSVSum": desc{7, "Read in CSV then calculate sum"},
	"sum2": desc{8, "Sum two columns"},
}

// SampleSizes is all the potential sample sizes and the order in which they should appear in the comparison table.
var SampleSizes = []string{"100k", "500k"}

// RunGoProfiler specifies all the benchmarks to profile and return in the benchmark table.
func RunGoProfiler() Results {
	fmt.Println("Profiling Go")
	Results := Results{
		"100k": {
			"sum": ProfileGo(benchmarkSumFloat64_100000),
			// 	"mean":       ProfileGo(benchmarkMeanFloat64_100000),
			// 	"median":     ProfileGo(benchmarkMedianFloat64_100000),
			// 	"min":        ProfileGo(benchmarkMinFloat64_100000),
			// 	"max":        ProfileGo(benchmarkMaxFloat64_100000),
			// 	"std":        ProfileGo(benchmarkStdFloat64_100000),
			// 	"readCSVSum": ProfileGo(benchmarkReadSumFloat64_100000),
		},
		"500k": {
			"sum2": ProfileGo(benchmarkSumFloat64_500000),
		},
	}
	return Results
}

// Results contains benchmarking results
// {"num of samples": {"test1": "10ms"...}}
type Results map[string]map[string]Result

// A Result of benchmarking data in the form [string, float64]
type Result []interface{}

// ProfileGo runs the normal Go benchmarking command but formats the result as a rounded string
// and raw ns float
func ProfileGo(f func(b *testing.B)) Result {
	benchmark := testing.Benchmark(f).NsPerOp()
	var speed string
	switch {
	case benchmark < int64(time.Microsecond):
		speed = fmt.Sprintf("%vns", benchmark)
	case benchmark < int64(time.Millisecond):
		speed = fmt.Sprintf("%.1fμs", float64(benchmark)/float64(time.Microsecond))
	case benchmark < int64(time.Second):
		speed = fmt.Sprintf("%.1fms", float64(benchmark)/float64(time.Millisecond))
	default:
		speed = fmt.Sprintf("%.2fs", float64(benchmark)/float64(time.Second))
	}
	return Result{speed, float64(benchmark)}
}

// RunPythonProfiler executes main.py in this directory, which is expected to return JSON
// in the form of Results. This command is expected to be initiated from the directory above.
func RunPythonProfiler() Results {
	fmt.Println("Profiling Python")
	_, thisFile, _, _ := runtime.Caller(0)
	script := "main.py"
	scriptPath := path.Join(path.Dir(thisFile), script)
	cmd := exec.Command("python", scriptPath)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	var r Results
	err = json.Unmarshal(out, &r)
	if err != nil {
		log.Fatal(err)
	}
	return r
}
