package series

import (
	"reflect"
	"testing"

	"github.com/ptiger10/pd/options"
)

func TestElement(t *testing.T) {
	s, err := New([]string{"", "valid"}, Idx([]string{"A", "B"}), Idx([]int{1, 2}))
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
	wantIdxKinds := []options.DataType{options.String, options.Int64}
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
		if !reflect.DeepEqual(got.LabelKinds, wantIdxKinds) {
			t.Errorf("Element returned kind %v, want %v", got.LabelKinds, wantIdxKinds)
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
	s.Name = "foo"
	origS, _ := New("foo")
	origS.Name = "foo"
	copyS := s.Copy()
	copyS.index.Levels[0].Labels.Set(0, "5")
	copyS.values.Set(0, "bar")
	copyS.Name = "bar"
	copyS.datatype = options.Bool
	if !seriesEquals(s, origS) || seriesEquals(s, copyS) {
		t.Errorf("s.copy() retained references to original, want fresh copy")
	}
}

// func Test_Equals(t *testing.T) {
// 	s, err := New("foo", Idx("bar"), options.Name("baz"))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	s2, _ := New("foo", Idx("bar"), options.Name("baz"))
// 	if !seriesEquals(s, s2) {
// 		t.Errorf("seriesEquals() returned false, want true")
// 	}
// 	s2.datatype = options.Bool
// 	if seriesEquals(s, s2) {
// 		t.Errorf("seriesEquals() returned true for different kind, want false")
// 	}

// 	s2, _ = New("quux", Idx("bar"), options.Name("baz"))
// 	if seriesEquals(s, s2) {
// 		t.Errorf("seriesEquals() returned true for different values, want false")
// 	}
// 	s2, _ = New("foo", Idx("corge"), options.Name("baz"))
// 	if seriesEquals(s, s2) {
// 		t.Errorf("seriesEquals() returned true for different index, want false")
// 	}
// 	s2, _ = New("foo", Idx("bar"), options.Name("qux"))
// 	if seriesEquals(s, s2) {
// 		t.Errorf("seriesEquals() returned true for different name, want false")
// 	}
// }

func TestReplaceNil(t *testing.T) {
	s := MustNew(nil)
	s2 := MustNew([]int{1, 2})
	s.replace(s2)
	if !seriesEquals(s, s2) {
		t.Errorf("Series.replace() returned %v, want %v", s, s2)
	}
}
