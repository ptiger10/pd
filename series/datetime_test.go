package series

import (
	"math"
	"testing"
	"time"
)

func TestConstructor_DateTime(t *testing.T) {
	s, err := New([]time.Time{
		time.Date(2019, 4, 18, 15, 0, 0, 0, time.UTC),
		time.Date(2019, 4, 19, 15, 0, 0, 0, time.UTC),
		time.Time{}})
	if err != nil {
		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
	}
	_, err = s.Sum()
	if err == nil {
		t.Errorf("Returned nil error when when summing datetime series, want error")
	}

	gotCount := s.Count()
	wantCount := 2
	if gotCount != wantCount {
		t.Errorf("Count() returned %v, want %v", gotCount, wantCount)
	}

}

func TestConstructor_InterfaceDateTime(t *testing.T) {
	s, err := New(
		[]interface{}{
			time.Date(2019, 4, 18, 15, 0, 0, 0, time.UTC), "1/1/18", "Jan 1, 2018", "January 1 2018",
			"1pm", "1", // times cannot be parsed
			[]string{"1", "2"}, // slice cannot be parsed
			time.Location{},    // struct other than time.Time cannot be parsed
			nil, complex64(1), "", "n/a", "N/A", "nan", "NaN", math.NaN()},
		SeriesType(DateTime))
	if err != nil {
		t.Errorf("Unable to create new series for %v: %v", t.Name(), err)
		return
	}

	gotCount := s.Count()
	wantCount := 4
	if gotCount != wantCount {
		t.Errorf("Count() returned %v, want %v", gotCount, wantCount)
	}
}
