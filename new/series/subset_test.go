package series

import (
	"reflect"
	"testing"
)

func TestAt_singleindex(t *testing.T) {
	var tests = []struct {
		input int
		want  Series
	}{
		{0, mustNew([]string{"hot"}, Index([]int{10}))},
		{1, mustNew([]string{"dog"}, Index([]int{100}))},
		{2, mustNew([]string{"log"}, Index([]int{1000}))},
	}
	for _, test := range tests {
		s, _ := New([]string{"hot", "dog", "log"}, Index([]int{10, 100, 1000}))
		got := s.At(test.input)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Returned Series %v, want %v", got, test.want)
		}
	}
}

func TestAt_multiindex(t *testing.T) {
	var tests = []struct {
		input int
		want  Series
	}{
		{0, mustNew([]string{"hot"}, Index([]int{10}), Index([]string{"A"}))},
		{1, mustNew([]string{"dog"}, Index([]int{100}), Index([]string{"B"}))},
		{2, mustNew([]string{"log"}, Index([]int{1000}), Index([]string{"C"}))},
	}
	for _, test := range tests {
		s, _ := New(
			[]string{"hot", "dog", "log"},
			Index([]int{10, 100, 1000}),
			Index([]string{"A", "B", "C"}),
		)
		got := s.At(test.input)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Returned Series %v, want %v", got, test.want)
		}
	}
}
