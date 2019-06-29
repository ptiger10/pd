package series

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ptiger10/pd/options"
)

func TestSeries_DataType(t *testing.T) {
	var tests = []struct {
		datatype options.DataType
		expected string
	}{

		{options.None, "none"},
		{options.Float64, "float64"},
		{options.Int64, "int64"},
		{options.String, "string"},
		{options.Bool, "bool"},
		{options.DateTime, "dateTime"},
		{options.Interface, "interface"},
		{options.Unsupported, "unsupported"},
		{-1, "unknown"},
		{100, "unknown"},
	}
	for _, test := range tests {
		s, _ := New([]int{1, 2, 3})
		s.datatype = test.datatype
		if s.DataType() != test.expected {
			t.Errorf("s.Datatype() for datatype %v returned %v, want %v", test.datatype, test.datatype.String(), test.expected)
		}
	}
}

func TestSeries_Equal(t *testing.T) {
	s, err := New("foo", Config{Index: "bar", Name: "baz"})
	if err != nil {
		t.Error(err)
	}
	s2, _ := New("foo", Config{Index: "bar", Name: "baz"})
	if !Equal(s, s2) {
		t.Errorf("Equal() returned false, want true")
	}
	s2.datatype = options.Bool
	if Equal(s, s2) {
		t.Errorf("Equal() returned true for different kind, want false")
	}

	s2, _ = New("quux", Config{Index: "bar", Name: "baz"})
	if Equal(s, s2) {
		t.Errorf("Equal() returned true for different values, want false")
	}
	s2, _ = New("foo", Config{Index: "corge", Name: "baz"})
	if Equal(s, s2) {
		t.Errorf("Equal() returned true for different index, want false")
	}
	s2, _ = New("foo", Config{Index: "bar", Name: "qux"})
	if Equal(s, s2) {
		t.Errorf("Equal() returned true for different name, want false")
	}
}

func TestSeries_ReplaceNil(t *testing.T) {
	s := MustNew(nil)
	s2 := MustNew([]int{1, 2})
	s.replace(s2)
	if !Equal(s, s2) {
		t.Errorf("Series.replace() returned %v, want %v", s, s2)
	}
}

func TestSeries_Values(t *testing.T) {
	s := MustNew([]string{"foo", "bar", "baz"})
	got := s.Values()
	want := []interface{}{"foo", "bar", "baz"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("s.Values() got %v, want %v", got, want)
	}
}

func TestSeries_MaxWidth(t *testing.T) {
	s := MustNew([]string{"foo", "quux", "grault"}, Config{Name: "grapply"})
	got := s.MaxWidth()
	want := 6
	if got != want {
		t.Errorf("s.MaxWidth got %v, want %v", got, want)
	}
}

func TestSeries_Describe_unsupported(t *testing.T) {
	s := MustNew([]float64{1, 2, 3})
	tm := s.Earliest()
	if (time.Time{}) != tm {
		t.Errorf("Earliest() got %v, want time.Time{} for unsupported type", tm)
	}
	tm = s.Latest()
	if (time.Time{}) != tm {
		t.Errorf("Latest() got %v, want time.Time{} for unsupported type", tm)
	}
}

// [START ensure tests]
func TestSeries_EnsureTypes_fail(t *testing.T) {
	defer log.SetOutput(os.Stderr)
	vals := []interface{}{1, 2, 3}

	var buf bytes.Buffer
	log.SetOutput(&buf)
	ensureFloatFromNumerics(vals)
	if buf.String() == "" {
		t.Errorf("ensureNumerics() returned no log message, want log due to fail")
	}
	buf.Reset()

	ensureDateTime(vals)
	if buf.String() == "" {
		t.Errorf("ensureDateTime() returned no log message, want log due to fail")
	}
	buf.Reset()

	ensureBools(vals)
	if buf.String() == "" {
		t.Errorf("ensureBools() returned no log message, want log due to fail")
	}
	buf.Reset()
}

// [END ensure tests]
