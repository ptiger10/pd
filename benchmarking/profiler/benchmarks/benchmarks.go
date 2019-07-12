// +build benchmarks

package benchmarks

import (
	"log"
	"testing"

	"github.com/ptiger10/pd"
	"github.com/ptiger10/pd/options"
)

func benchmarkSumFloat64_5m(b *testing.B) {
	for n := 0; n < b.N; n++ {
		df5m.Sum()
	}
}

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

func benchmarkSyncMeanFloat64_100000(b *testing.B) {
	options.SetAsync(false)
	for n := 0; n < b.N; n++ {
		df100k.Mean()
	}
	options.RestoreDefaults()
}

func benchmarkMeanFloat64_500000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		df500k.Mean()
	}
}

func benchmarkSyncMeanFloat64_500000(b *testing.B) {
	options.SetAsync(false)
	for n := 0; n < b.N; n++ {
		df500k.Mean()
	}
	options.RestoreDefaults()
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

func benchmarkSyncStdFloat64_100000(b *testing.B) {
	options.SetAsync(false)
	for n := 0; n < b.N; n++ {
		df100k.Std()
	}
	options.RestoreDefaults()
}

func benchmarkStdFloat64_500000(b *testing.B) {
	for n := 0; n < b.N; n++ {
		df500k.Max()
	}
}

func benchmarkSyncStdFloat64_500000(b *testing.B) {
	options.SetAsync(false)
	for n := 0; n < b.N; n++ {
		df500k.Std()
	}
	options.RestoreDefaults()
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
