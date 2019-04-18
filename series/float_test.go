package series

import (
	"math"
	"testing"
)

func TestMath_Float(t *testing.T) {
	var tests = []struct {
		input      []float64
		wantSum    float64
		wantCount  int
		wantMean   float64
		wantMedian float64
	}{
		{
			input:      []float64{-2.5, 0, 0, 3.5, math.NaN()},
			wantSum:    1.0,
			wantCount:  4,
			wantMean:   .25,
			wantMedian: 0,
		},
		{
			input:      []float64{-1, 3, 4, math.NaN()},
			wantSum:    6,
			wantCount:  3,
			wantMean:   2,
			wantMedian: 3,
		},
	}
	for _, test := range tests {
		s, _ := New(test.input)
		gotSum, _ := s.Sum()
		if gotSum != test.wantSum {
			t.Errorf("Sum() returned %v for input %v, want %v", gotSum, test.input, test.wantSum)
		}
		gotCount := s.Count()
		if gotCount != test.wantCount {
			t.Errorf("Count() returned %v for input %v, want %v", gotSum, test.input, test.wantSum)
		}
		gotMean, _ := s.Mean()
		if gotMean != test.wantMean {
			t.Errorf("Mean() returned %v for input %v, want %v", gotSum, test.input, test.wantSum)
		}
		gotMedian, _ := s.Median()
		if gotMedian != test.wantMedian {
			t.Errorf("Median() returned %v for input %v, want %v", gotSum, test.input, test.wantSum)
		}
	}
}

func TestAdd_Float(t *testing.T) {
	s, _ := New([]float64{math.NaN(), 1, 1, math.NaN(), 1})
	s, err := s.AddConst(float64(1))
	if err != nil {
		t.Errorf("Unable to add constant to float: %v", err)
	}
	gotSum, _ := s.Sum()
	wantSum := 6.0
	if gotSum != wantSum {
		t.Errorf("Sum() returned %v, want %v", gotSum, wantSum)
	}

	gotCount := s.Count()
	wantCount := 3
	if gotCount != wantCount {
		t.Errorf("Count() returned %v, want %v", gotCount, wantCount)
	}

	unsupported := "dog"
	_, err = s.AddConst(unsupported)
	if err == nil {
		t.Errorf("Returned nil error when adding unsupported type (%T) to Float, want error", unsupported)
	}
}

func TestConstructor_Float32(t *testing.T) {
	s, _ := New([]float32{1, -2, 3.5, 0})
	got, _ := s.Sum()
	want := 2.5
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}

func TestConstructor_Float64(t *testing.T) {
	s, _ := New([]float64{1, -2, 3.5, 0})
	got, _ := s.Sum()
	want := 2.5
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}

func TestConstructor_InterfaceFloat(t *testing.T) {
	s, err := New(
		[]interface{}{float32(1), float64(1.5), 0.5, int32(1), 1, uint64(2), "0.5", "1", nil, complex64(1), "", "n/a", "N/A", "nan", "NaN", math.NaN()},
		SeriesType(Float))
	if err != nil {
		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
		return
	}

	gotSum, _ := s.Sum()
	wantSum := 8.5
	if gotSum != wantSum {
		t.Errorf("Sum() returned %v, want %v", gotSum, wantSum)
	}

	gotCount := s.Count()
	wantCount := 8
	if gotCount != wantCount {
		t.Errorf("Count() returned %v, want %v", gotCount, wantCount)
	}
}
