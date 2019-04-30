package series

import (
	"reflect"
	"testing"
)

func TestAt_int_multi_int_int(t *testing.T) {

	var tests = []struct {
		input      int
		wantLength int
		wantVal    interface{}
		wantIdx    interface{}
		wantLevels int
	}{
		{0, 1, "hot", []interface{}{int64(5), int64(100)}, 2},
		{1, 1, "dog", []interface{}{int64(6), int64(101)}, 2},
	}
	for _, test := range tests {
		s, _ := New([]string{"hot", "dog", "log"}, Index([]int{5, 6, 7}), Index([]int{100, 101, 102}))
		got := s.At(test.input)
		gotLength := got.Len()

		if gotLength != test.wantLength {
			t.Errorf("Returned Series of length %d, want %d", gotLength, test.wantLength)
		}
		gotVal := got.Elem().Value
		if gotVal != test.wantVal {
			t.Errorf("Returned value %s, want %s", gotVal, test.wantVal)
		}
		gotIdx := got.Elem().Index
		if !reflect.DeepEqual(gotIdx, test.wantIdx) {
			t.Errorf("Returned index %#v, want %#v", gotIdx, test.wantIdx)
		}
		gotLevels := len(got.Index.Levels)
		if gotLevels != test.wantLevels {
			t.Errorf("Returned %d index level(s), want %d", gotLevels, test.wantLevels)
		}

	}

}
