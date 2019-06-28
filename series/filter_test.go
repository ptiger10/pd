package series

import (
	"reflect"
	"testing"
	"time"
)

func TestSubset(t *testing.T) {
	tests := []struct {
		args    []int
		want    *Series
		wantErr bool
	}{
		{[]int{0}, MustNew("foo"), false},
		{[]int{1}, MustNew("bar", Config{Index: 1}), false},
		{[]int{0, 1}, MustNew([]string{"foo", "bar"}), false},
		{[]int{1, 0}, MustNew([]string{"bar", "foo"}, Config{Index: []int{1, 0}}), false},
		{[]int{}, newEmptySeries(), true},
		{[]int{3}, newEmptySeries(), true},
	}
	for _, tt := range tests {
		s := MustNew([]string{"foo", "bar", "baz"})
		got, err := s.Subset(tt.args)
		if (err != nil) != tt.wantErr {
			t.Errorf("s.Subset() error = %v, want %v for args %v", err, tt.wantErr, tt.args)
		}
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
		{"GT", (*Series).GT, 2, []int{2}},
		{"GTE", (*Series).GTE, 2, []int{1, 2}},
		{"LT", (*Series).LT, 2, []int{0}},
		{"LTE", (*Series).LTE, 2, []int{0, 1}},
		{"EQ", (*Series).EQ, 2, []int{1}},
		{"NEQ", (*Series).NEQ, 2, []int{0, 2}},
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

func TestFilterBool(t *testing.T) {
	tests := []struct {
		name string
		fn   func(*Series) []int
		want []int
	}{
		{"True", (*Series).True, []int{1}},
		{"False", (*Series).False, []int{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MustNew([]bool{false, true})
			got := tt.fn(s)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("s.Filter() got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterDateTime(t *testing.T) {
	tests := []struct {
		name string
		fn   func(*Series, time.Time) []int
		arg  time.Time
		want []int
	}{
		{"Before", (*Series).Before, time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC), []int{0}},
		{"After", (*Series).After, time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC), []int{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MustNew([]time.Time{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2019, 3, 1, 0, 0, 0, 0, time.UTC)})
			got := tt.fn(s, tt.arg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("s.Filter() got %v, want %v", got, tt.want)
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
