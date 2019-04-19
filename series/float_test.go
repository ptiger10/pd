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
		wantMin    float64
		wantMax    float64
	}{
		{
			input:      []float64{-2.5, 0, 0, 3.5, math.NaN()},
			wantSum:    1.0,
			wantCount:  4,
			wantMean:   .25,
			wantMedian: 0,
			wantMin:    -2.5,
			wantMax:    3.5,
		},
		{
			input:      []float64{-1, 3, 4, math.NaN()},
			wantSum:    6,
			wantCount:  3,
			wantMean:   2,
			wantMedian: 3,
			wantMin:    -1,
			wantMax:    4,
		},
		{
			input:      []float64{},
			wantSum:    0,
			wantCount:  0,
			wantMean:   math.NaN(),
			wantMedian: math.NaN(),
			wantMin:    math.NaN(),
			wantMax:    math.NaN(),
		},
		{
			input:      []float64{math.NaN(), math.NaN()},
			wantSum:    0,
			wantCount:  0,
			wantMean:   math.NaN(),
			wantMedian: math.NaN(),
			wantMin:    math.NaN(),
			wantMax:    math.NaN(),
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
			t.Errorf("Count() returned %v for input %v, want %v", gotCount, test.input, test.wantCount)
		}
		gotMean, _ := s.Mean()
		if gotMean != test.wantMean && !math.IsNaN(gotMean) {
			t.Errorf("Mean() returned %v for input %v, want %v", gotMean, test.input, test.wantMean)
		}
		gotMedian, _ := s.Median()
		if gotMedian != test.wantMedian && !math.IsNaN(gotMedian) {
			t.Errorf("Median() returned %v for input %v, want %v", gotMedian, test.input, test.wantMedian)
		}
		gotMin, _ := s.Min()
		if gotMin != test.wantMin && !math.IsNaN(gotMin) {
			t.Errorf("Min() returned %v for input %v, want %v", gotMin, test.input, test.wantMin)
		}
		gotMax, _ := s.Max()
		if gotMax != test.wantMax && !math.IsNaN(gotMax) {
			t.Errorf("Max() returned %v for input %v, want %v", gotMax, test.input, test.wantMax)
		}
	}
}

func TestAdd_Float(t *testing.T) {
	s, _ := New([]float64{math.NaN(), 1, 1, math.NaN(), 1})
	// Adding an int
	s, err := s.AddConst(1)
	if err != nil {
		t.Errorf("Unable to add int to float: %v", err)
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

	// Adding a float
	s, err = s.AddConst(1.0)
	if err != nil {
		t.Errorf("Unable to add float to float: %v", err)
	}

	gotSum, _ = s.Sum()
	wantSum = 9.0
	if gotSum != wantSum {
		t.Errorf("Sum() returned %v, want %v", gotSum, wantSum)
	}
	gotCount = s.Count()
	wantCount = 3
	if gotCount != wantCount {
		t.Errorf("Count() returned %v, want %v", gotCount, wantCount)
	}

	// Adding a string
	unsupported := "dog"
	_, err = s.AddConst(unsupported)
	if err == nil {
		t.Errorf("Returned nil error when adding unsupported type (%T) to Float, want error", unsupported)
	}
}

func TestFloat_Unsupported(t *testing.T) {
	s, _ := New([]float64{1, 2, 3, 4})
	_, err := s.ValueCounts()
	if err == nil {
		t.Errorf("Returned nil error when calling ValueCounts() on float series, want error")
	}
}
