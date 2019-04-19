package series

import (
	"math"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestFilter_Float(t *testing.T) {
	s, _ := New([]float64{1, 2, 3, 4, 5, math.NaN()})
	s, err := s.FilterFloat(func(v float64) bool { return v > 3 })
	if err != nil {
		t.Errorf("Unable to filter Float Series: %v", err)
	}
	gotSum, _ := s.Sum()
	wantSum := 9.0
	if gotSum != wantSum {
		t.Errorf("Sum() returned %v after filtering, want %v", gotSum, wantSum)
	}

	gotCount := s.Count()
	wantCount := 2
	if gotCount != wantCount {
		t.Errorf("Count() returned %v after filtering, want %v", gotCount, wantCount)
	}
}

func TestFilter_Int(t *testing.T) {
	s, _ := New([]interface{}{1, 2, 3, 4, 5, "", ""}, SeriesType(Int))
	s, err := s.FilterInt(func(v int64) bool { return v > 2 && v < 5 })
	if err != nil {
		t.Errorf("Unable to filter Int Series: %v", err)
	}
	gotSum, _ := s.Sum()
	wantSum := 7.0
	if gotSum != wantSum {
		t.Errorf("Sum() returned %v after filtering, want %v", gotSum, wantSum)
	}

	gotCount := s.Count()
	wantCount := 2
	if gotCount != wantCount {
		t.Errorf("Count() returned %v after filtering, want %v", gotCount, wantCount)
	}
}

func TestFilter_String(t *testing.T) {
	s, _ := New([]string{"Highest", "High", "", "Medium", "Medium", "Low"})
	s, err := s.FilterString(func(v string) bool { return !strings.Contains(v, "High") })
	if err != nil {
		t.Errorf("Unable to filter String Series: %v", err)
	}
	got, _ := s.ValueCounts()
	want := map[string]int{"Medium": 2, "Low": 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ValueCount() returned %v after filtering, want %v", got, want)
	}

	gotCount := s.Count()
	wantCount := 3
	if gotCount != wantCount {
		t.Errorf("Count() returned %v, want %v", gotCount, wantCount)
	}
}

func TestFilter_DateTime(t *testing.T) {
	s, _ := New([]time.Time{time.Date(2018, 12, 31, 0, 0, 0, 0, time.UTC), time.Date(2019, 1, 2, 0, 0, 0, 0, time.UTC)})
	s, err := s.FilterDateTime(func(v time.Time) bool { return v.After(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)) })
	if err != nil {
		t.Errorf("Unable to filter String Series: %v", err)
	}

	gotCount := s.Count()
	wantCount := 1
	if gotCount != wantCount {
		t.Errorf("Count() returned %v, want %v", gotCount, wantCount)
	}
}

func TestFilter_Unsupported(t *testing.T) {
	s, _ := New([]bool{true, true})
	_, err := s.FilterFloat(func(v float64) bool { return v > 3 })
	if err == nil {
		t.Errorf("Returned nil error when calling unsupported FilterFloat on Bool series, want error")
	}
	_, err = s.FilterInt(func(v int64) bool { return v > 3 })
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
}
