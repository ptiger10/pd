package series_test

import (
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/ptiger10/pd/cmp"
	"github.com/ptiger10/pd/series"
)

func TestFilter_Float(t *testing.T) {
	var tests = []struct {
		name      string
		filterFn  func(float64) bool
		wantSum   float64
		wantCount int
	}{
		{"manual", func(v float64) bool { return v > 2 && v < 5 }, 7, 2},
		{"gt", cmp.Gt(8), 9, 1},
		{"gte", cmp.Gte(8), 17, 2},
		{"lt", cmp.Lt(3), 3, 2},
		{"lte", cmp.Lte(3), 6, 3},
		{"eq", cmp.Eq(5), 5, 1},
		{"neq", cmp.Neq(10), 45, 9},
	}
	s, _ := series.New([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, math.NaN()})
	for _, test := range tests {
		sFilt, err := s.FilterFloat(test.filterFn)
		if err != nil {
			t.Errorf("Unable to filter Float Series for test %v: %v", test.name, err)
		}
		gotSum, _ := sFilt.Sum()
		if gotSum != test.wantSum {
			t.Errorf("Test %v: Sum() returned %v, want %v", test.name, gotSum, test.wantSum)
		}

		gotCount := sFilt.Count()
		if gotCount != test.wantCount {
			t.Errorf("Test %v: Count() returned %v, want %v", test.name, gotCount, test.wantCount)
		}
	}

}

func TestFilter_Int(t *testing.T) {
	var tests = []struct {
		name      string
		filterFn  func(float64) bool
		wantSum   float64
		wantCount int
	}{
		{"manual", func(v float64) bool { return v > 2 && v < 5 }, 7, 2},
		{"gt", cmp.Gt(8), 9, 1},
		{"gte", cmp.Gte(8), 17, 2},
		{"lt", cmp.Lt(3), 3, 2},
		{"lte", cmp.Lte(3), 6, 3},
		{"eq", cmp.Eq(5), 5, 1},
		{"neq", cmp.Neq(10), 45, 9},
	}
	s, _ := series.New([]interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, "", ""}, series.Type(series.Int))
	for _, test := range tests {
		sFilt, err := s.FilterInt(test.filterFn)
		if err != nil {
			t.Errorf("Unable to filter Float Series for test %v: %v", test.name, err)
		}
		gotSum, _ := sFilt.Sum()
		if gotSum != test.wantSum {
			t.Errorf("Test %v: Sum() returned %v, want %v", test.name, gotSum, test.wantSum)
		}

		gotCount := sFilt.Count()
		if gotCount != test.wantCount {
			t.Errorf("Test %v: Count() returned %v, want %v", test.name, gotCount, test.wantCount)
		}
	}
}

func TestFilter_String(t *testing.T) {
	var tests = []struct {
		name           string
		filterFn       func(string) bool
		wantValueCount map[string]int
		wantCount      int
	}{
		{"manual", func(v string) bool { return strings.HasSuffix(v, "est") }, map[string]int{"Highest": 1}, 1},
		{"in", cmp.In([]string{"Medium"}), map[string]int{"Medium": 2}, 2},
		{"in-2", cmp.In([]string{"High"}), map[string]int{"High": 1}, 1},
		{"nin", cmp.Nin([]string{"Highest", "High", "Medium"}), map[string]int{"Low": 1}, 1},
	}

	s, _ := series.New([]string{"Highest", "High", "", "Medium", "Medium", "Low"})
	for _, test := range tests {
		sFilt, err := s.FilterString(test.filterFn)
		if err != nil {
			t.Errorf("Unable to filter String Series for test %v: %v", test.name, err)
		}
		gotValueCount, _ := sFilt.ValueCounts()
		if !reflect.DeepEqual(gotValueCount, test.wantValueCount) {
			t.Errorf("Test %v: ValueCount() returned %v, want %v", test.name, gotValueCount, test.wantValueCount)
		}

		gotCount := sFilt.Count()
		if gotCount != test.wantCount {
			t.Errorf("Test %v: Count() returned %v, want %v", test.name, gotCount, test.wantCount)
		}

	}
}

func TestFilter_DateTime(t *testing.T) {
	var tests = []struct {
		name      string
		filterFn  func(time.Time) bool
		wantCount int
	}{
		{"manual", func(v time.Time) bool { return v.Equal(time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)) }, 1},
		{"before", cmp.Before(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)), 1},
		{"after", cmp.Before(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)), 1},
	}

	s, _ := series.New([]time.Time{
		time.Date(2018, 12, 31, 0, 0, 0, 0, time.UTC),
		time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC),
		time.Time{}})
	for _, test := range tests {
		sFilt, err := s.FilterDateTime(test.filterFn)
		if err != nil {
			t.Errorf("Unable to filter DateTime Series: %v", err)
		}
		gotCount := sFilt.Count()
		if gotCount != test.wantCount {
			t.Errorf("Test %v: Count() returned %v, want %v", test.name, gotCount, test.wantCount)
		}
	}

}

func TestFilter_Unsupported(t *testing.T) {
	s, _ := series.New([]bool{true, true})
	_, err := s.FilterFloat(func(v float64) bool { return v > 3 })
	if err == nil {
		t.Errorf("Returned nil error when calling unsupported FilterFloat on Bool series, want error")
	}
	_, err = s.FilterInt(func(v float64) bool { return v > 3 })
	if err == nil {
		t.Errorf("Returned nil error when calling unsupported FilterInt on Bool series, want error")
	}
	_, err = s.FilterString(func(v string) bool { return strings.Contains(v, "cats") })
	if err == nil {
		t.Errorf("Returned nil error when calling unsupported FilterString on Bool series, want error")
	}
	_, err = s.FilterDateTime(func(v time.Time) bool { return v.After(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)) })
	if err == nil {
		t.Errorf("Returned nil error when calling unsupported FilterDateTime on Bool series, want error")
	}

	s, _ = series.New([]int{1, 2})
	_, err = s.FilterFloat(func(v float64) bool { return v > 3 })
	if err == nil {
		t.Errorf("Returned nil error when calling unsupported FilterFloat on Int series, want error")
	}

	s, _ = series.New([]float64{1, 2})
	_, err = s.FilterInt(func(v float64) bool { return v > 3 })
	if err == nil {
		t.Errorf("Returned nil error when calling unsupported FilterInt on Float series, want error")
	}
}
