package series

import (
	"math"
	"testing"
)

func TestMath_Int(t *testing.T) {
	var tests = []struct {
		input      []int64
		wantSum    float64
		wantCount  int
		wantMean   float64
		wantMedian float64
		wantMin    float64
		wantMax    float64
	}{
		{
			input:      []int64{-2, 0, 3, 3},
			wantSum:    4,
			wantCount:  4,
			wantMean:   1,
			wantMedian: 1.5,
			wantMin:    -2,
			wantMax:    3,
		},
		{
			input:      []int64{-2, 0, 5},
			wantSum:    3,
			wantCount:  3,
			wantMean:   1,
			wantMedian: 0,
			wantMin:    -2,
			wantMax:    5,
		},
		{
			input:      []int64{},
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
