package series

import (
	"reflect"
	"testing"
)

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

func TestValueCounts(t *testing.T) {
	s, _ := New([]string{"low", "", "high", "high", "high"})
	got, _ := s.ValueCounts()
	want := map[string]int{"high": 3, "low": 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("s.ValueCounts() returned %v, want %v", got, want)
	}
}
