package series

import (
	"log"
	"testing"
)

// Floats
func float32Slice(n int) []float32 {
	var l []float32
	for i := 0; i < n; i++ {
		l = append(l, 1)
	}
	return l
}

func float64Slice(n int) []float64 {
	var l []float64
	for i := 0; i < n; i++ {
		l = append(l, 1)
	}
	return l
}

func benchmarkNewFloat32(i int, b *testing.B) {
	v := float32Slice(i)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := New(v)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func benchmarkNewFloat64(i int, b *testing.B) {
	v := float64Slice(i)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := New(v)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// func BenchmarkNewFloat32_1(b *testing.B) { benchmarkNewFloat32(10000, b) }
// func BenchmarkNewFloat64_1(b *testing.B) { benchmarkNewFloat64(10000, b) }
// func BenchmarkNewFloat32_2(b *testing.B) { benchmarkNewFloat32(100000, b) }

// func BenchmarkNewFloat64_2(b *testing.B) { benchmarkNewFloat64(100000, b) }
// func BenchmarkNewFloat32_3(b *testing.B) { benchmarkNewFloat32(1000000, b) }

// func BenchmarkNewFloat64_3(b *testing.B) { benchmarkNewFloat64(1000000, b) }

// func BenchmarkNewFloat32_4(b *testing.B) { benchmarkNewFloat32(10000000, b) }
func BenchmarkNewFloat64_4(b *testing.B) { benchmarkNewFloat64(10000000, b) }
