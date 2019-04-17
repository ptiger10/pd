package series

import (
	"fmt"
	"testing"
)

func TestFloat32(t *testing.T) {
	s, _ := New([]float32{1, -2, 3.5, 0})
	got, _ := s.Sum()
	want := 2.5
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}

func TestFloat64(t *testing.T) {
	s, _ := New([]float64{1, -2, 3.5, 0})
	got, _ := s.Sum()
	want := 2.5
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}

func TestFloatNaN(t *testing.T) {
	s, err := New(
		[]interface{}{1, nil, "", 3.5, "na", "nan", "NaN", "3.5"},
		SeriesType(Float))
	if err != nil {
		t.Errorf("%v returned err, nil expected: %v", t.Name(), err)
		return
	}
	fmt.Println(s)
	got, _ := s.Sum()
	want := 8.0
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}

func TestInt(t *testing.T) {
	s, _ := New([]int{1, -2, 3, 0})
	got, _ := s.Sum()
	want := 2.0
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}
func TestInt32(t *testing.T) {
	s, _ := New([]int32{1, -2, 3, 0})
	got, _ := s.Sum()
	want := 2.0
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}
func TestInt64(t *testing.T) {
	s, _ := New([]int64{1, -2, 3, 0})
	got, _ := s.Sum()
	want := 2.0
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}

func TestUInt(t *testing.T) {
	s, _ := New([]uint{1, 2, 3, 0})
	got, _ := s.Sum()
	want := 6.0
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}
func TestUInt32(t *testing.T) {
	s, _ := New([]uint32{1, 2, 3, 0})
	got, _ := s.Sum()
	want := 6.0
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}
func TestUInt64(t *testing.T) {
	s, _ := New([]uint64{1, 2, 3, 0})
	got, _ := s.Sum()
	want := 6.0
	if got != want {
		t.Errorf("Sum() returned %v, want %v", got, want)
	}
}

func TestUnsupported(t *testing.T) {
	_, err := New([]bool{true, true})
	if err == nil {
		t.Errorf("Returned nil, want error")
	}
}
