package types

import (
	"reflect"
	"testing"
)

func TestSum_series(t *testing.T) {
	s := &Series{
		Values: []interface{}{1, 3, 5},
		Index: Index{
			IntIdx: []int{0, 1, 2}, IntValMap: nil, StringIdx: nil, StringValMap: nil},
		Type: reflect.Int,
	}
	got := s.Sum()
	want := 9
	if got != want {
		t.Errorf("%s returned %v, want %v", t.Name(), got, want)
	}

}
