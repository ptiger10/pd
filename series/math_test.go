package series

import (
	"math"
	"testing"
)

func TestSeriesMath(t *testing.T) {
	s, _ := New([]int{1, 2, 3})
	if sum := s.Sum(); sum != 6 {
		t.Errorf("s.Sum() returned %v, want %v", sum, 6)
	}
	if mean := s.Mean(); mean != 2 {
		t.Errorf("s.Mean() returned %v, want %v", mean, 2)
	}

}

func TestMath_numerics(t *testing.T) {
	var tests = []struct {
		s          *Series
		wantSum    float64
		wantMean   float64
		wantMedian float64
		wantMin    float64
		wantMax    float64
		wantQ1     float64
		wantQ2     float64
		wantQ3     float64
		wantStd    float64
	}{
		{mustNew([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}), 45, 5, 5, 1, 9, 2.5, 5, 7.5, 2.58},
		{mustNew([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}), 45, 5, 5, 1, 9, 2.5, 5, 7.5, 2.58},
		{mustNew([]float64{-1, 2, 3, 4}), 8, 2, 2.5, -1, 4, 0.5, 2.5, 3.5, 1.87},
	}
	for _, test := range tests {
		gotSum := test.s.Sum()
		if gotSum != test.wantSum {
			t.Errorf("Sum() returned %v for input\n%v, want %v", gotSum, test.s, test.wantSum)
		}
		gotMean := test.s.Mean()
		if gotMean != test.wantMean {
			t.Errorf("Mean() returned %v for input\n%v, want %v", gotMean, test.s, test.wantMean)
		}
		gotMedian := test.s.Median()
		if gotMedian != test.wantMedian {
			t.Errorf("Median() returned %v for input\n%v, want %v", gotMedian, test.s, test.wantMedian)
		}
		gotMin := test.s.Min()
		if gotMin != test.wantMin {
			t.Errorf("Min() returned %v for input\n%v, want %v", gotMin, test.s, test.wantMin)
		}
		gotMax := test.s.Max()
		if gotMax != test.wantMax {
			t.Errorf("Max() returned %v for input\n%v, want %v", gotMax, test.s, test.wantMax)
		}
		gotQ1 := test.s.Quartile(1)
		if gotQ1 != test.wantQ1 {
			t.Errorf("Quartile(1) returned %v for input\n%v, want %v", gotQ1, test.s, test.wantQ1)
		}
		gotQ2 := test.s.Quartile(2)
		if gotQ2 != test.wantQ2 {
			t.Errorf("Quartile(2) returned %v for input\n%v, want %v", gotQ2, test.s, test.wantQ2)
		}
		gotQ3 := test.s.Quartile(3)
		if gotQ3 != test.wantQ3 {
			t.Errorf("Quartile(3) returned %v for input\n%v, want %v", gotQ3, test.s, test.wantQ3)
		}
		gotStd := test.s.Std()
		if math.Round(gotStd*100)/100 != math.Round(test.wantStd*100)/100 {
			t.Errorf("Std() returned %v for input\n%v, want %v", gotStd, test.s, test.wantStd)
		}
	}
}

func TestMath_bool(t *testing.T) {
	var tests = []struct {
		s        *Series
		wantSum  float64
		wantMean float64
	}{
		{mustNew([]bool{true, false}), 1, .5},
		{mustNew([]bool{false}), 0, 0},
	}
	for _, test := range tests {
		gotSum := test.s.Sum()
		if gotSum != test.wantSum {
			t.Errorf("Sum() returned %v for input\n%v, want %v", gotSum, test.s, test.wantSum)
		}
		gotMean := test.s.Mean()
		if gotMean != test.wantMean {
			t.Errorf("Mean() returned %v for input\n%v, want %v", gotMean, test.s, test.wantMean)
		}
	}
}
