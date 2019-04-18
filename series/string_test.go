package series

import (
	"math"
	"testing"
)

func TestConstructor_SliceString(t *testing.T) {
	s, err := New([]string{"low", "", "high"})
	if err != nil {
		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
	}
	_, err = s.Sum()
	if err == nil {
		t.Errorf("Returned nil error when summing string series, want error")
	}
	_, err = s.Median()
	if err == nil {
		t.Errorf("Returned nil error when finding median of string series, want error")
	}
	_, err = s.Mean()
	if err == nil {
		t.Errorf("Returned nil error when finding mean of string series, want error")
	}
}

func TestConstructor_InterfaceString(t *testing.T) {
	s, err := New(
		[]interface{}{float32(1), float64(1.5), 0.5, int32(1), 1, uint64(2), "0.5", "1", nil, complex64(1), "", "n/a", "N/A", "nan", "NaN", math.NaN()},
		SeriesType(String))
	if err != nil {
		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
		return
	}
	gotCount := s.Count()
	wantCount := 8
	if gotCount != wantCount {
		t.Errorf("Count() returned %v, want %v", gotCount, wantCount)
	}
}
