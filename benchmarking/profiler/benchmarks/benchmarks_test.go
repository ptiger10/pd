package benchmarks

import (
	"testing"
)

func BenchmarkMath(b *testing.B) {
	benchmarks := []struct {
		name string
		fn   func(b *testing.B)
	}{
		// {"100k sum 1 column", benchmarkSumFloat64_100000},
		// {"100k read then sum 1 column", benchmarkReadSumFloat64_100000},
		// {"100k mean 1 column", benchmarkMeanFloat64_100000},
		// {"100k sync mean 1 column", benchmarkSyncMeanFloat64_100000},
		// {"100k median 1 column", benchmarkMedianFloat64_100000},
		// {"100k min 1 column", benchmarkMinFloat64_100000},
		// {"100k max 1 column", benchmarkMaxFloat64_100000},
		// {"100k std 1 column", benchmarkStdFloat64_100000},
		// {"100k sync std 1 column", benchmarkSyncStdFloat64_100000},
		// {"500k std 2 columns", benchmarkStdFloat64_500000},
		// {"500k sync std 2 columns", benchmarkSyncStdFloat64_500000},
		// {"500k sum 2 columns", benchmarkSumFloat64_500000},
		// {"500k mean 2 columns", benchmarkMeanFloat64_500000},
		// {"500k sync mean 2 columns", benchmarkSyncMeanFloat64_500000},
		{"5m sum 1 column", benchmarkSumFloat64_5m},
	}
	ReadData()
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			bm.fn(b)
		})
	}
}
