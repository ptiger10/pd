package series

import (
	"reflect"
	"testing"
)

func TestSubset(t *testing.T) {
	tests := []struct {
		args []int
		want *Series
	}{
		{[]int{0}, MustNew("foo")},
		{[]int{1}, MustNew("bar", Config{Index: 1})},
		{[]int{0, 1}, MustNew([]string{"foo", "bar"})},
		{[]int{1, 0}, MustNew([]string{"bar", "foo"}, Config{Index: []int{1, 0}})},
	}
	for _, tt := range tests {
		s := MustNew([]string{"foo", "bar", "baz"})
		got := s.Subset(tt.args)
		if !Equal(got, tt.want) {
			t.Errorf("s.Subset() got %v, want %v for args %v", got, tt.want, tt.args)
		}
	}
}

func TestFilterFloat64(t *testing.T) {
	tests := []struct {
		name string
		fn   func(*Series, float64) []int
		arg  float64
		want []int
	}{
		{"Gt", (*Series).Gt, 2, []int{2}},
		{"Gte", (*Series).Gte, 2, []int{1, 2}},
		{"Lt", (*Series).Lt, 2, []int{0}},
		{"Lte", (*Series).Lte, 2, []int{0, 1}},
		{"Eq", (*Series).Eq, 2, []int{1}},
		{"Neq", (*Series).Neq, 2, []int{0, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MustNew([]float64{1, 2, 3})
			got := tt.fn(s, tt.arg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("s.Filter() got %v, want %v for arg %v", got, tt.want, tt.arg)
			}
		})
	}
}

func TestFilterString(t *testing.T) {
	s := MustNew([]string{"foo", "bar", "baz"})
	got := s.Contains("ba")
	want := []int{1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("s.Contains() got %v, want %v", got, want)
	}
}
