package series

import (
	"testing"
)

func TestAt_int_multi_int_int(t *testing.T) {
	s, _ := New([]string{"1", "3", "5"}, Index([]int{0, 1, 2}, Index([]int{100, 101, 102})))
	var tests = []struct {
		input      int
		wantLength int
		wantVal    int64
		wantIdx    []int64
	}{
		{0, 1, 1, []int64{0, 100}},
		{1, 1, 3, []int64{1, 101}},
	}
	for _, test := range tests {
		got := s.At(test.input)
		gotLength := got.Len()
		if gotLength != test.wantLength {
			t.Errorf("Returned Series of length %d, want %d", gotLength, test.wantLength)
		}
		gotVal := got.Elem().Value.(int64)
		if gotVal != test.wantVal {
			t.Errorf("Returned value %d, want %d", gotVal, test.wantVal)
		}

	}

}
