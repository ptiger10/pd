package series

import (
	"math"
	"testing"
)

func TestConstructor_Bool(t *testing.T) {
	s, err := New([]bool{true, true, false, true})
	if err != nil {
		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
	}
	got, _ := s.Sum()
	want := 3.0
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}

func TestConstructor_InterfaceBool(t *testing.T) {
	s, err := New(
		[]interface{}{float32(1), float64(1.5), 0.5, int32(1), 1, uint64(2), "0.5", "1", nil, complex64(1), "", "n/a", "N/A", "nan", "NaN", math.NaN()},
		SeriesType(Bool))
	if err != nil {
		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
		return
	}
	gotSum, _ := s.Sum()
	wantSum := 8.0
	if gotSum != wantSum {
		t.Errorf("Sum() returned %v, want %v", gotSum, wantSum)
	}

	gotCount := s.Count()
	wantCount := 8
	if gotCount != wantCount {
		t.Errorf("Count() returned %v, want %v", gotCount, wantCount)
	}
}

func TestAdd_BoolUnsupported(t *testing.T) {
	s, _ := New([]bool{true, true, false, true})
	_, err := s.AddConst(1)
	if err == nil {
		t.Error("Returned nil error when adding constant to Bool, want error")
	}
}
