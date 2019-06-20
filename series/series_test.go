package series

import (
	"reflect"
	"testing"

	"github.com/ptiger10/pd/options"
)

func TestElement(t *testing.T) {
	s, err := New([]string{"", "valid"}, Config{MultiIndex: []interface{}{[]string{"A", "B"}, []int{1, 2}}})
	if err != nil {
		t.Error(err)
	}
	var tests = []struct {
		position int
		wantVal  interface{}
		wantNull bool
		wantIdx  []interface{}
	}{
		{0, "NaN", true, []interface{}{"A", int64(1)}},
		{1, "valid", false, []interface{}{"B", int64(2)}},
	}
	wantIdxTypes := []options.DataType{options.String, options.Int64}
	for _, test := range tests {
		got := s.Element(test.position)
		if got.Value != test.wantVal {
			t.Errorf("Element returned value %v, want %v", got.Value, test.wantVal)
		}
		if got.Null != test.wantNull {
			t.Errorf("Element returned bool %v, want %v", got.Null, test.wantNull)
		}
		if !reflect.DeepEqual(got.Labels, test.wantIdx) {
			t.Errorf("Element returned index %#v, want %#v", got.Labels, test.wantIdx)
		}
		if !reflect.DeepEqual(got.LabelTypes, wantIdxTypes) {
			t.Errorf("Element returned kind %v, want %v", got.LabelTypes, wantIdxTypes)
		}
	}
}
func TestDatatype(t *testing.T) {
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

func Test_Copy(t *testing.T) {
	s, _ := New("foo")
	s.name = "foo"
	sOrig, _ := New("foo")
	sOrig.name = "foo"
	sCopy := s.Copy()
	sCopy.index.Levels[0].Labels.Set(0, 5)
	sCopy.values.Set(0, "bar")
	sCopy.name = "foobar"
	sCopy.index.Refresh()
	want, _ := New("bar", Config{Index: 5, Name: "foobar"})
	if !Equal(sCopy, want) {
		t.Errorf("s.Copy() returned %v, want %v", sCopy.index, want.index)
	}
	if !Equal(s, sOrig) || Equal(s, sCopy) {
		t.Errorf("s.copy() retained references to original, want fresh copy")
	}
}

func Test_Equals(t *testing.T) {
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

func TestReplaceNil(t *testing.T) {
	s := MustNew(nil)
	s2 := MustNew([]int{1, 2})
	s.replace(s2)
	if !Equal(s, s2) {
		t.Errorf("Series.replace() returned %v, want %v", s, s2)
	}
}

func TestMaxWidth(t *testing.T) {
	s := MustNew([]string{"foo", "quux", "grault"}, Config{Name: "grapply"})
	got := s.MaxWidth()
	want := 6
	if got != want {
		t.Errorf("s.MaxWidth got %v, want %v", got, want)
	}
}
