package benchmarks

import (
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ptiger10/pd"
	"github.com/ptiger10/pd/dataframe"
)

// Descriptions of the benchmarking tests
var Descriptions = map[string]desc{
	"sum":        desc{1, "Sum one column"},
	"mean":       desc{2, "Simple mean of one column"},
	"median":     desc{3, "Median of one column"},
	"min":        desc{4, "Min of one column"},
	"max":        desc{5, "Max of one column"},
	"std":        desc{6, "Standard deviation of one column"},
	"readCSVSum": desc{7, "Read in CSV then calculate sum"},
	"sum2":       desc{8, "Sum two columns"},
}

// SampleSizes is all the potential sample sizes and the order in which they should appear in the comparison table.
var SampleSizes = []string{"100k", "500k"}

var df100k *dataframe.DataFrame
var df500k *dataframe.DataFrame

func read100k() {
	var err error
	df100k, err = pd.ReadCSV(getPath("100k"), pd.ReadOptions{HeaderRows: 1})
	if err != nil {
		log.Fatal(err)
	}

	got := math.Round(df100k.Sum().At(0).(float64)*100) / 100
	want := 50408.63
	if got != want {
		log.Fatalf("Reading in test data: df.Sum() got %v, want %v", got, want)
	}

	got = math.Round(df100k.Mean().At(0).(float64)*100) / 100
	want = 0.5
	if got != want {
		log.Fatalf("Reading in test data: df.Mean() got %v, want %v", got, want)
	}

	got = math.Round(df100k.Median().At(0).(float64)*100) / 100
	want = 0.50
	if got != want {
		log.Fatalf("Reading in test data: df.Median() got %v, want %v", got, want)
	}

	got = math.Round(df100k.Min().At(0).(float64)*100) / 100
	want = 0.0
	if got != want {
		log.Fatalf("Reading in test data: df.Min() got %v, want %v", got, want)
	}

	got = math.Round(df100k.Max().At(0).(float64)*100) / 100
	want = 1.0
	if got != want {
		log.Fatalf("Reading in test data: df.Max() got %v, want %v", got, want)
	}

	got = math.Round(df100k.Std().At(0).(float64)*100) / 100
	want = 0.29
	if got != want {
		log.Fatalf("Reading in test data: df.Std() got %v, want %v", got, want)
	}

}

func read500k() {
	var err error
	df500k, err = pd.ReadCSV(getPath("500k"), pd.ReadOptions{HeaderRows: 1})
	if err != nil {
		log.Fatal(err)
	}

	got := math.Round(df500k.Sum().At(0).(float64)*100) / 100
	want := 130598.19
	if got != want {
		log.Fatalf("Reading in test data: df.Sum500() got %v, want %v", got, want)
	}
}

// ReadData initializes data for use in comparison tetss
func ReadData() {
	read100k()
	read500k()
}

var files = map[string]string{
	"100k": "../dataRandom100k1Col.csv",
	"500k": "../dataRandom500k2Col.csv",
}

func getPath(s string) string {
	basename, ok := files[s]
	if !ok {
		log.Fatalf("Reading in test data: df.%v not in %v", s, files)
	}
	_, thisFile, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(thisFile), basename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("Reading in test data: df.File does not exist at %s", path)
	}
	return path
}
