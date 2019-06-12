package series

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ptiger10/pd/datatypes"
	"github.com/ptiger10/pd/opt"
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
	wantIdxKinds := []datatypes.DataType{datatypes.String, datatypes.Int64}
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
func TestKind(t *testing.T) {
	var tests = []struct {
		datatype datatypes.DataType
		expected string
	}{

		{datatypes.None, "none"},
		{datatypes.Float64, "float64"},
		{datatypes.Int64, "int64"},
		{datatypes.String, "string"},
		{datatypes.Bool, "bool"},
		{datatypes.DateTime, "time.Time"},
		{datatypes.Interface, "interface"},
		{datatypes.Unsupported, "unsupported"},
		{-1, "unknown"},
		{100, "unknown"},
	}
	for _, test := range tests {
		s, _ := New([]int{1, 2, 3})
		s.datatype = test.datatype
		if s.DataType() != test.expected {
			t.Errorf("s.Kind() for kind %v returned %v, want %v", test.datatype, test.datatype.String(), test.expected)
		}
	}
}

func Test_Copy(t *testing.T) {
	s, _ := New("foo")
	s.Name = "foo"
	origS, _ := New("foo")
	origS.Name = "foo"
	copyS := s.copy()
	copyS.index.Levels[0].Labels.Set(0, "5")
	copyS.values.Set(0, "bar")
	copyS.Name = "bar"
	copyS.datatype = datatypes.Bool
	MathPtr := fmt.Sprintf("%p", copyS.Math.s)
	ToPtr := fmt.Sprintf("%p", copyS.To.s)
	IndexToPtr := fmt.Sprintf("%p", copyS.Index.s)
	if !seriesEquals(s, origS) || seriesEquals(s, copyS) || fmt.Sprintf("%p", &s) == MathPtr {
		t.Errorf("s.copy() retained references to original, want fresh copy")
	}
	if copyS.Math.s == nil || copyS.To.s == nil || copyS.Index.s == nil {
		t.Errorf("s.copy() did not instantiate new pointers for embedded structs")
	}
	if MathPtr != ToPtr || MathPtr != IndexToPtr {
		t.Errorf("s.copy() did not instantiate pointers for embedded structs correctly")
	}
}

func Test_Equals(t *testing.T) {
	s, err := New("foo", Idx("bar"), opt.Name("baz"))
	if err != nil {
		t.Error(err)
	}
	s2, _ := New("foo", Idx("bar"), opt.Name("baz"))
	if !seriesEquals(s, s2) {
		t.Errorf("seriesEquals() returned false, want true")
	}
	s2.datatype = datatypes.Bool
	if seriesEquals(s, s2) {
		t.Errorf("seriesEquals() returned true for different kind, want false")
	}

	s2, _ = New("quux", Idx("bar"), opt.Name("baz"))
	if seriesEquals(s, s2) {
		t.Errorf("seriesEquals() returned true for different values, want false")
	}
	s2, _ = New("foo", Idx("corge"), opt.Name("baz"))
	if seriesEquals(s, s2) {
		t.Errorf("seriesEquals() returned true for different index, want false")
	}
	s2, _ = New("foo", Idx("bar"), opt.Name("qux"))
	if seriesEquals(s, s2) {
		t.Errorf("seriesEquals() returned true for different name, want false")
	}
}

func TestReplaceNil(t *testing.T) {
	s := mustNew(nil)
	s2 := mustNew([]int{1, 2})
	s.replace(&s2)
	if !seriesEquals(s, s2) {
		t.Errorf("Series.replace() returned %v, want %v", s, s2)
	}
}
