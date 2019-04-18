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
	}{
		{
			input:      []int64{-2, 0, 3, 3},
			wantSum:    4,
			wantCount:  4,
			wantMean:   1,
			wantMedian: 1.5,
		},
		{
			input:      []int64{-2, 0, 5},
			wantSum:    3,
			wantCount:  3,
			wantMean:   1,
			wantMedian: 0,
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

func TestConstructor_Int(t *testing.T) {
	s, _ := New([]int{1, -2, 3, 0})
	got, _ := s.Sum()
	want := 2.0
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}
func TestConstructor_Int32(t *testing.T) {
	s, _ := New([]int32{1, -2, 3, 0})
	got, _ := s.Sum()
	want := 2.0
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}
func TestConstructor_Int64(t *testing.T) {
	s, _ := New([]int64{1, -2, 3, 0})
	got, _ := s.Sum()
	want := 2.0
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}

func TestConstructor_UInt(t *testing.T) {
	s, err := New([]uint{1, 2, 3, 0})
	if err != nil {
		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
	}
	got, _ := s.Sum()
	want := 6.0
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}
func TestConstructor_UInt32(t *testing.T) {
	s, err := New([]uint32{1, 2, 3, 0})
	if err != nil {
		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
	}
	got, _ := s.Sum()
	want := 6.0
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}
func TestConstructor_UInt64(t *testing.T) {
	s, err := New([]uint64{1, 2, 3, 0})
	if err != nil {
		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
	}
	got, _ := s.Sum()
	want := 6.0
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}

func TestConstructor_InterfaceInt(t *testing.T) {
	s, err := New(
		[]interface{}{float32(1), float64(1.5), 0.5, int32(1), 1, uint64(2), "0.5", "1", nil, complex64(1), "", "n/a", "N/A", "nan", "NaN", math.NaN()},
		SeriesType(Int))
	if err != nil {
		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
		return
	}
	gotSum, _ := s.Sum()
	wantSum := 7.0
	if gotSum != wantSum {
		t.Errorf("Sum() returned %v, want %v", gotSum, wantSum)
	}

	gotCount := s.Count()
	wantCount := 8
	if gotCount != wantCount {
		t.Errorf("Count() returned %v, want %v", gotCount, wantCount)
	}
}
