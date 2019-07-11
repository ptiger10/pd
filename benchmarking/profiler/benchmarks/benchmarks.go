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
	for n := 0; n < b.N; n++ {
		df.Sum()
	}
	got := math.Round(df.Sum().At(0).(float64)*100) / 100
	want := 50408.63
	if got != want {
		log.Fatalf("Sum() got %v, want %v", got, want)
	}
}

func benchmarkMeanFloat64_100000(b *testing.B) {
	df, _ := pd.ReadCSV("RandomNumbers.csv", pd.ReadOptions{HeaderRows: 1})
	for n := 0; n < b.N; n++ {
		df.Mean()
	}
	got := math.Round(df.Mean().At(0).(float64)*100) / 100
	want := 0.5
	if got != want {
		log.Fatalf("Mean() got %v, want %v", got, want)
	}
}

func benchmarkMedianFloat64_100000(b *testing.B) {
	df, _ := pd.ReadCSV("RandomNumbers.csv", pd.ReadOptions{HeaderRows: 1})
	for n := 0; n < b.N; n++ {
		df.Median()
	}
	got := math.Round(df.Median().At(0).(float64)*100) / 100
	want := 0.50
	if got != want {
		log.Fatalf("Median() got %v, want %v", got, want)
	}
}

func benchmarkMinFloat64_100000(b *testing.B) {
	df, _ := pd.ReadCSV("RandomNumbers.csv", pd.ReadOptions{HeaderRows: 1})
	for n := 0; n < b.N; n++ {
		df.Min()
	}
	got := math.Round(df.Min().At(0).(float64)*100) / 100
	want := 0.0
	if got != want {
		log.Fatalf("Min() got %v, want %v", got, want)
	}
}

func benchmarkMaxFloat64_100000(b *testing.B) {
	df, _ := pd.ReadCSV("RandomNumbers.csv", pd.ReadOptions{HeaderRows: 1})
	for n := 0; n < b.N; n++ {
		df.Max()
	}
	got := math.Round(df.Max().At(0).(float64)*100) / 100
	want := 1.0
	if got != want {
		log.Fatalf("Max() got %v, want %v", got, want)
	}
}

func benchmarkStdFloat64_100000(b *testing.B) {
	df, _ := pd.ReadCSV("RandomNumbers.csv", pd.ReadOptions{HeaderRows: 1})
	for n := 0; n < b.N; n++ {
		df.Max()
	}
	got := math.Round(df.Std().At(0).(float64)*100) / 100
	want := 0.29
	if got != want {
		log.Fatalf("Std() got %v, want %v", got, want)
	}
}

func benchmarkReadSumFloat64_100000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		df, err := pd.ReadCSV("RandomNumbers.csv", pd.ReadOptions{HeaderRows: 1})
		if err != nil {
			log.Fatal(err)
		}
		df.Sum()
	}
	df, _ := pd.ReadCSV("RandomNumbers.csv", pd.ReadOptions{HeaderRows: 1})
	got := math.Round(df.Sum().At(0).(float64)*100) / 100
	want := 50408.63
	if got != want {
		log.Fatalf("ReadCSV() then Sum() got %v, want %v", got, want)
	}
}
