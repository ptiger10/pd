package benchmarks

import "testing"

func BenchmarkSumFloat64_1000000(b *testing.B) {
	benchmarkSumFloat64_100000(b)
}

func BenchmarkMeanFloat64_1000000(b *testing.B) {
	benchmarkMeanFloat64_100000(b)
}

func BenchmarkMedianFloat64_1000000(b *testing.B) {
	benchmarkMedianFloat64_100000(b)
}
