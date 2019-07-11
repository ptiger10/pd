package benchmarks

import (
	"log"
	"math"
	"testing"

	"github.com/ptiger10/pd"
)

func benchmarkSumFloat64_100000(b *testing.B) {
	df, err := pd.ReadCSV("RandomNumbers.csv", pd.ReadOptions{HeaderRows: 1})
	if err != nil {
		log.Fatal(err)
	}
	// run the sum function
	for n := 0; n < b.N; n++ {
		df.Sum()
	}
	got := math.Round(df.Sum().At(0).(float64)*100) / 100
	want := 50408.63
	if got != want {
		b.Errorf("Sum() got %v, want %v", got, want)
	}
}

func benchmarkMeanFloat64_100000(b *testing.B) {
	df, err := pd.ReadCSV("RandomNumbers.csv", pd.ReadOptions{HeaderRows: 1})
	if err != nil {
		log.Fatal(err)
	}
	// run the mean function
	for n := 0; n < b.N; n++ {
		df.Mean()
	}
	got := math.Round(df.Mean().At(0).(float64)*100) / 100
	want := 0.50
	if got != want {
		b.Errorf("Mean() got %v, want %v", got, want)
	}
}

func benchmarkMedianFloat64_100000(b *testing.B) {
	df, err := pd.ReadCSV("RandomNumbers.csv", pd.ReadOptions{HeaderRows: 1})
	if err != nil {
		log.Fatal(err)
	}
	for n := 0; n < b.N; n++ {
		df.Median()
	}
	got := math.Round(df.Median().At(0).(float64)*100) / 100
	want := 0.50
	if got != want {
		b.Errorf("Median() got %v, want %v", got, want)
	}
}
