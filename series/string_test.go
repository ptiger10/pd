package series

import (
	"math"
	"reflect"
	"testing"
)

func TestConstructor_SliceString(t *testing.T) {
	_, err := New([]string{"low", "", "high"})
	if err != nil {
		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
	}
}

func TestString_Unsupported(t *testing.T) {
	s, _ := New([]string{"low", "", "high"})
	_, err := s.Sum()
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

	_, err = s.Min()
	if err == nil {
		t.Errorf("Returned nil error when finding min of string series, want error")
	}

	_, err = s.Max()
	if err == nil {
		t.Errorf("Returned nil error when finding max of string series, want error")
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

func TestValueCounts(t *testing.T) {
	s, _ := New([]string{"low", "", "high", "high", "high"})
	got, _ := s.ValueCounts()
	want := map[string]int{"high": 3, "low": 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("s.ValueCounts() returned %v, want %v", got, want)
	}
}
