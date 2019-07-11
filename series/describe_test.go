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

func TestSeries_Describe(t *testing.T) {
	type want struct {
		len          int
		numIdxLevels int
		maxWidth     int
		values       []interface{}
		vals         interface{}
		datatype     string
		name         string
		valid        []int
		null         []int
	}
	tests := []struct {
		name  string
		input *Series
		want  want
	}{
		{"empty",
			newEmptySeries(),
			want{len: 0, numIdxLevels: 0, maxWidth: 0,
				values: []interface{}{}, vals: []interface{}{}, datatype: "none", name: "",
				valid: []int{}, null: []int{}}},
		{name: "default index",
			input: MustNew([]string{"foo", "", "bar", ""}),
			want: want{len: 4, numIdxLevels: 1, maxWidth: 3,
				values:   []interface{}{"foo", "NaN", "bar", "NaN"},
				vals:     []string{"foo", "NaN", "bar", "NaN"},
				datatype: "string", name: "",
				valid: []int{0, 2}, null: []int{1, 3}}},
		{"multi index",
			MustNew(
				1.0,
				Config{MultiIndex: []interface{}{"baz", "qux"}, Name: "foo"},
			),
			want{len: 1, numIdxLevels: 2, maxWidth: 4,
				values: []interface{}{1.0}, vals: []float64{1},
				datatype: "float64", name: "foo",
				valid: []int{0}, null: []int{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.input.Copy()
			gotLen := s.Len()
			if gotLen != tt.want.len {
				t.Errorf("s.Len(): got %v, want %v", gotLen, tt.want.len)
			}
			gotNumIdxLevels := s.NumLevels()
			if gotNumIdxLevels != tt.want.numIdxLevels {
				t.Errorf("s.NumLevels(): got %v, want %v", gotNumIdxLevels, tt.want.numIdxLevels)
			}
			gotMaxWidth := s.MaxWidth()
			if gotMaxWidth != tt.want.maxWidth {
				t.Errorf("s.MaxWidth(): got %v, want %v", gotMaxWidth, tt.want.maxWidth)
			}
			gotValues := s.Values()
			if !reflect.DeepEqual(gotValues, tt.want.values) {
				t.Errorf("s.Values(): got %v, want %v", gotValues, tt.want.values)
			}
			gotVals := s.Vals()
			if !reflect.DeepEqual(gotVals, tt.want.vals) {
				t.Errorf("s.Vals(): got %#v, want %v", gotVals, tt.want.vals)
			}
			gotDatatype := s.DataType()
			if gotDatatype != tt.want.datatype {
				t.Errorf("s.Datatype(): got %v, want %v", gotDatatype, tt.want.datatype)
			}
			gotName := s.Name()
			if gotName != tt.want.name {
				t.Errorf("s.Name(): got %v, want %v", gotName, tt.want.name)
			}
			gotValid := s.valid()
			if !reflect.DeepEqual(gotValid, tt.want.valid) {
				t.Errorf("s.valid(): got %v, want %v", gotValid, tt.want.valid)
			}
			gotNull := s.null()
			if !reflect.DeepEqual(gotNull, tt.want.null) {
				t.Errorf("s.null(): got %v, want %v", gotNull, tt.want.null)
			}
		})
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
