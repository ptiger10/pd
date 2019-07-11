package benchmarks

import (
	"log"
	"testing"

	"github.com/ptiger10/pd"
)

func benchmarkSumFloat64_500000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		df500k.Sum()
	}

}

func benchmarkSumFloat64_100000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		df100k.Sum()
	}
}

func benchmarkMeanFloat64_100000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		df100k.Mean()
	}

}

func benchmarkMedianFloat64_100000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		df100k.Median()
	}

}

func benchmarkMinFloat64_100000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		df100k.Min()
	}

}

func benchmarkMaxFloat64_100000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		df100k.Max()
	}

}

func benchmarkStdFloat64_100000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		df100k.Max()
	}

}

func benchmarkReadSumFloat64_100000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		df, err := pd.ReadCSV(getPath("100k"), pd.ReadOptions{HeaderRows: 1})
		if err != nil {
			log.Fatal(err)
		}
		df.Sum()
	}
}
