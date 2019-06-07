package series

import (
	"reflect"
	"testing"
	"time"
)

func filterNumeric(vals []float64, fn func(float64) bool) []bool {
	var ret []bool
	for _, val := range vals {
		ret = append(ret, fn(val))
	}
	return ret
}

func filterString(vals []string, fn func(string) bool) []bool {
	var ret []bool
	for _, val := range vals {
		ret = append(ret, fn(val))
	}
	return ret
}

func filterBool(vals []bool, fn func(bool) bool) []bool {
	var ret []bool
	for _, val := range vals {
		ret = append(ret, fn(val))
	}
	return ret
}

func filterTime(vals []time.Time, fn func(time.Time) bool) []bool {
	var ret []bool
	for _, val := range vals {
		ret = append(ret, fn(val))
	}
	return ret
}
func TestNumerics(t *testing.T) {
	var tests = []struct {
		input   float64
		wantGt  []bool
		wantGte []bool
		wantLt  []bool
		wantLte []bool
		wantEq  []bool
	}{
		{
			3, []bool{false, false, false}, []bool{false, false, true},
			[]bool{true, true, false}, []bool{true, true, true},
			[]bool{false, false, true},
		},
	}
	for _, test := range tests {
		vals := []float64{1, 2, 3}
		gotGt := filterNumeric(vals, Gt(test.input))
		if !reflect.DeepEqual(gotGt, test.wantGt) {
			t.Errorf("Gt() returned %v for input %v, want %v", gotGt, vals, test.wantGt)
		}
		gotGte := filterNumeric(vals, Gte(test.input))
		if !reflect.DeepEqual(gotGte, test.wantGte) {
			t.Errorf("Gte() returned %v for input %v, want %v", gotGte, vals, test.wantGte)
		}
		gotLt := filterNumeric(vals, Lt(test.input))
		if !reflect.DeepEqual(gotLt, test.wantLt) {
			t.Errorf("Lt() returned %v for input %v, want %v", gotLt, vals, test.wantLt)
		}
		gotLte := filterNumeric(vals, Lte(test.input))
		if !reflect.DeepEqual(gotLte, test.wantLte) {
			t.Errorf("Lte() returned %v for input %v, want %v", gotLte, vals, test.wantLte)
		}
		gotEq := filterNumeric(vals, Eq(test.input))
		if !reflect.DeepEqual(gotEq, test.wantEq) {
			t.Errorf("Eq() returned %v for input %v, want %v", gotEq, vals, test.wantEq)
		}

	}
}

func TestGt(t *testing.T) {
	s, _ := New([]float64{1, 2, 3})
	newS, err := s.Filter.Gt(2)
	if err != nil {
		t.Error(err)
	}
	want, _ := New([]float64{3})
	if seriesEquals(newS, want) {
		t.Errorf("s.Filter.Gt() returned %v, want %v", newS, want)
	}
}

func TestGt_Null(t *testing.T) {
	s, _ := New([]string{"", "1", "", "2", "", "3"})
	s = s.To.Float()
	newS, err := s.Filter.Gt(2)
	if err != nil {
		t.Error(err)
	}
	want, _ := New([]float64{3})
	if seriesEquals(newS, want) {
		t.Errorf("s.Filter.Gt() returned %v, want %v", newS, want)
	}
}
