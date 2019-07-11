package series

import (
	"math"
	"testing"

	"github.com/ptiger10/pd/options"
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
		name       string
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
		{"float with null", MustNew([]float64{math.NaN(), math.NaN(), 2, 3, 1}), 6, 2, 2, 1, 3, 1, 2, 3, 0.82},
		{"float from string with null", MustNew([]string{"", "", "1", "2", "3"}).ToFloat64(), 6, 2, 2, 1, 3, 1, 2, 3, 0.82},
		{"int from string with null", MustNew([]string{"", "", "1", "2", "3"}).ToInt64(), 6, 2, 2, 1, 3, 1, 2, 3, 0.82},
		{"int", MustNew([]int{2, 1, 3, 4, 5, 6, 7, 8, 9}), 45, 5, 5, 1, 9, 2.5, 5, 7.5, 2.58},
		{"float", MustNew([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}), 45, 5, 5, 1, 9, 2.5, 5, 7.5, 2.58},
		{"float with negative", MustNew([]float64{2, -1, 4, 3}), 8, 2, 2.5, -1, 4, 0.5, 2.5, 3.5, 1.87},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSum := tt.s.Sum()
			if gotSum != tt.wantSum {
				t.Errorf("Sum()returned %v, want %v", gotSum, tt.wantSum)
			}
			gotMean := tt.s.Mean()
			if gotMean != tt.wantMean {
				t.Errorf("Mean()returned %v, want %v", gotMean, tt.wantMean)
			}
			gotMedian := tt.s.Median()
			if gotMedian != tt.wantMedian {
				t.Errorf("Median()returned %v, want %v", gotMedian, tt.wantMedian)
			}
			gotMin := tt.s.Min()
			if gotMin != tt.wantMin {
				t.Errorf("Min()returned %v, want %v", gotMin, tt.wantMin)
			}
			gotMax := tt.s.Max()
			if gotMax != tt.wantMax {
				t.Errorf("Max()returned %v, want %v", gotMax, tt.wantMax)
			}
			gotQ1 := tt.s.Quartile(1)
			if gotQ1 != tt.wantQ1 {
				t.Errorf("Quartile(1)returned %v, want %v", gotQ1, tt.wantQ1)
			}
			gotQ2 := tt.s.Quartile(2)
			if gotQ2 != tt.wantQ2 {
				t.Errorf("Quartile(2)returned %v, want %v", gotQ2, tt.wantQ2)
			}
			gotQ3 := tt.s.Quartile(3)
			if gotQ3 != tt.wantQ3 {
				t.Errorf("Quartile(3)returned %v, want %v", gotQ3, tt.wantQ3)
			}
			gotStd := tt.s.Std()
			if math.Round(gotStd*100)/100 != math.Round(tt.wantStd*100)/100 {
				t.Errorf("Std()returned %v, want %v", gotStd, tt.wantStd)
			}
		})

	}
}

func TestMath_numerics_async(t *testing.T) {
	var tests = []struct {
		name       string
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
		{"float with null", MustNew([]float64{math.NaN(), math.NaN(), 2, 3, 1}), 6, 2, 2, 1, 3, 1, 2, 3, 0.82},
	}
	for _, tt := range tests {
		options.SetAsync(false)
		defer options.RestoreDefaults()
		t.Run(tt.name, func(t *testing.T) {
			gotSum := tt.s.Sum()
			if gotSum != tt.wantSum {
				t.Errorf("Sum()returned %v, want %v", gotSum, tt.wantSum)
			}
			gotMean := tt.s.Mean()
			if gotMean != tt.wantMean {
				t.Errorf("Mean()returned %v, want %v", gotMean, tt.wantMean)
			}
			gotMedian := tt.s.Median()
			if gotMedian != tt.wantMedian {
				t.Errorf("Median()returned %v, want %v", gotMedian, tt.wantMedian)
			}
			gotMin := tt.s.Min()
			if gotMin != tt.wantMin {
				t.Errorf("Min()returned %v, want %v", gotMin, tt.wantMin)
			}
			gotMax := tt.s.Max()
			if gotMax != tt.wantMax {
				t.Errorf("Max()returned %v, want %v", gotMax, tt.wantMax)
			}
			gotQ1 := tt.s.Quartile(1)
			if gotQ1 != tt.wantQ1 {
				t.Errorf("Quartile(1)returned %v, want %v", gotQ1, tt.wantQ1)
			}
			gotQ2 := tt.s.Quartile(2)
			if gotQ2 != tt.wantQ2 {
				t.Errorf("Quartile(2)returned %v, want %v", gotQ2, tt.wantQ2)
			}
			gotQ3 := tt.s.Quartile(3)
			if gotQ3 != tt.wantQ3 {
				t.Errorf("Quartile(3)returned %v, want %v", gotQ3, tt.wantQ3)
			}
			gotStd := tt.s.Std()
			if math.Round(gotStd*100)/100 != math.Round(tt.wantStd*100)/100 {
				t.Errorf("Std()returned %v, want %v", gotStd, tt.wantStd)
			}
		})

	}
}

func TestMath_bool(t *testing.T) {
	var tests = []struct {
		s        *Series
		wantSum  float64
		wantMean float64
	}{
		{MustNew([]string{"", "true"}).ToBool(), 1, 1},
		{MustNew([]bool{true, false}), 1, .5},
		{MustNew([]bool{false}), 0, 0},
	}
	for _, tt := range tests {
		gotSum := tt.s.Sum()
		if gotSum != tt.wantSum {
			t.Errorf("Sum()returned %v, want %v", gotSum, tt.wantSum)
		}
		gotMean := tt.s.Mean()
		if gotMean != tt.wantMean {
			t.Errorf("Mean()returned %v, want %v", gotMean, tt.wantMean)
		}
	}
}

func TestMath_unsupported(t *testing.T) {
	var tests = []struct {
		name string
		s    *Series
	}{
		{"string", MustNew([]string{"foo"})},
		{"null", MustNew([]string{})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSum := tt.s.Sum()
			if !math.IsNaN(gotSum) {
				t.Errorf("Sum()returned %v, want NaN", gotSum)
			}
			gotMean := tt.s.Mean()
			if !math.IsNaN(gotMean) {
				t.Errorf("Mean()returned %v, want NaN", gotMean)
			}
			gotMedian := tt.s.Median()
			if !math.IsNaN(gotMedian) {
				t.Errorf("Median()returned %v, want NaN", gotMedian)
			}
			gotMin := tt.s.Min()
			if !math.IsNaN(gotMin) {
				t.Errorf("Min()returned %v, want NaN", gotMin)
			}
			gotMax := tt.s.Max()
			if !math.IsNaN(gotMax) {
				t.Errorf("Max()returned %v, want NaN", gotMax)
			}
			gotQ1 := tt.s.Quartile(1)
			if !math.IsNaN(gotQ1) {
				t.Errorf("Quartile(1)returned %v, want NaN", gotQ1)
			}
			gotQ2 := tt.s.Quartile(2)
			if !math.IsNaN(gotQ2) {
				t.Errorf("Quartile(2)returned %v, want NaN", gotQ2)
			}
			gotQ3 := tt.s.Quartile(3)
			if !math.IsNaN(gotQ3) {
				t.Errorf("Quartile(3)returned %v, want NaN", gotQ3)
			}
			gotStd := tt.s.Std()
			if !math.IsNaN(gotStd) {
				t.Errorf("Std()returned %v, want NaN", gotStd)
			}
		})
	}
}

func TestMath_unsupported_other(t *testing.T) {
	s := MustNew([]float64{})
	got := s.Std()
	if !math.IsNaN(got) {
		t.Errorf("Std()returned %v, want NaN", got)
	}
	got = s.Median()
	if !math.IsNaN(got) {
		t.Errorf("Median()returned %v, want NaN", got)
	}
	s = MustNew([]float64{1})
	got = s.Quartile(10)
	if !math.IsNaN(got) {
		t.Errorf("Quartile()returned %v, want NaN", got)
	}
}
