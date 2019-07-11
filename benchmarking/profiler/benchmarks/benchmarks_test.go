package benchmarks

import "testing"

func BenchmarkSumFloat64_100000(b *testing.B) {
	benchmarkSumFloat64_100000(b)
}

func BenchmarkMeanFloat64_100000(b *testing.B) {
	benchmarkMeanFloat64_100000(b)
}

func BenchmarkMedianFloat64_100000(b *testing.B) {
	benchmarkMedianFloat64_100000(b)
}

func BenchmarkMinFloat64_100000(b *testing.B) {
	benchmarkMinFloat64_100000(b)
}

func BenchmarkMaxFloat64_100000(b *testing.B) {
	benchmarkMaxFloat64_100000(b)
}

func BenchmarkStdFloat64_100000(b *testing.B) {
	benchmarkStdFloat64_100000(b)
}
