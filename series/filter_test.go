package series

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestApply(t *testing.T) {
	s := MustNew([]float64{1, 2, 3})
	sArchive := s.Copy()

	s.InPlace.Apply(func(val interface{}) interface{} {
		v, ok := val.(float64)
		if !ok {
			return ""
		}
		return ((v - s.Mean()) / s.Std())
	})
	want := MustNew([]float64{-1.224744871391589, 0, 1.224744871391589})
	if !Equal(s, want) {
		t.Errorf("InPlace.Apply() returned %v, want %v", s, want)
	}

	sCopy := sArchive.Apply(func(val interface{}) interface{} {
		v, ok := val.(float64)
		if !ok {
			return ""
		}
		return ((v - sArchive.Mean()) / sArchive.Std())
	})
	if !Equal(sCopy, want) {
		t.Errorf("Apply() returned %v, want %v", sCopy, want)
	}
	if Equal(sArchive, sCopy) {
		t.Errorf("Apply() retained access to original, want copy")
	}
}

func TestApply_riskier(t *testing.T) {
	s := MustNew([]float64{1, 2, 3})
	got := s.Apply(func(val interface{}) interface{} {
		return (val.(float64) - s.Mean()) / s.Std()
	})
	want := MustNew([]float64{-1.224744871391589, 0, 1.224744871391589})
	if !Equal(got, want) {
		t.Errorf("Apply() returned %v, want %v", got, want)
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

func TestFilter_Contains(t *testing.T) {
	s := MustNew([]string{"foo", "bar", "baz"})
	got := s.Contains("ba")
	want := []int{1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("s.Contains() got %v, want %v", got, want)
	}

	got = s.InList([]string{"foo", "bar"})
	want = []int{0, 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("s.In() got %v, want %v", got, want)
	}
}

func TestFilter_float(t *testing.T) {
	s := MustNew([]float64{1, 2, 3})
	got := s.Filter(func(val interface{}) bool {
		v, ok := val.(float64)
		if !ok {
			return false
		}
		if v > 2 {
			return true
		}
		return false
	})
	want := []int{2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("s.Filter() got %v, want %v", got, want)
	}
}

func TestFilter_string(t *testing.T) {
	s := MustNew([]string{"bamboo", "leaves", "taboo"})
	got := s.Filter(func(val interface{}) bool {
		v, ok := val.(string)
		if !ok {
			return false
		}
		if strings.HasSuffix(v, "boo") {
			return true
		}
		return false
	})
	want := []int{0, 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("s.Filter() got %v, want %v", got, want)
	}
}

func TestFilter_string_riskier(t *testing.T) {
	s := MustNew([]string{"bamboo", "leaves", "taboo"})
	got := s.Filter(func(val interface{}) bool {
		if strings.HasSuffix(val.(string), "boo") {
			return true
		}
		return false
	})
	want := []int{0, 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("s.Filter() got %v, want %v", got, want)
	}
}
