package series

import (
	"reflect"
	"testing"

	"github.com/ptiger10/pd/kinds"
	"github.com/ptiger10/pd/opt"
)

func TestKind(t *testing.T) {
	var tests = []struct {
		kind     kinds.Kind
		expected string
	}{

		{kinds.None, "none"},
		{kinds.Float, "float64"},
		{kinds.Int, "int64"},
		{kinds.String, "string"},
		{kinds.Bool, "bool"},
		{kinds.DateTime, "time.Time"},
		{kinds.Interface, "interface"},
		{kinds.Unsupported, "unsupported"},
		{-1, "unknown"},
		{100, "unknown"},
	}
	for _, test := range tests {
		s, _ := New([]int{1, 2, 3})
		s.kind = test.kind
		if s.Kind() != test.expected {
			t.Errorf("s.Kind() for kind %v returned %v, want %v", test.kind, test.kind.String(), test.expected)
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
	copyS.kind = kinds.Bool
	if !reflect.DeepEqual(s, origS) {
		t.Errorf("s.copy() returned original, want fresh copy")
	}
}

func Test_Equals(t *testing.T) {
	s, err := New("foo", Index("bar"), opt.Name("baz"))
	if err != nil {
		t.Error(err)
	}
	s2, _ := New("foo", Index("bar"), opt.Name("baz"))
	if !seriesEquals(s, s2) {
		t.Errorf("seriesEquals() returned false, want true")
	}
	s2.kind = kinds.Bool
	if seriesEquals(s, s2) {
		t.Errorf("seriesEquals() returned true for different kind, want false")
	}

	s2, _ = New("quux", Index("bar"), opt.Name("baz"))
	if seriesEquals(s, s2) {
		t.Errorf("seriesEquals() returned true for different values, want false")
	}
	s2, _ = New("foo", Index("corge"), opt.Name("baz"))
	if seriesEquals(s, s2) {
		t.Errorf("seriesEquals() returned true for different index, want false")
	}
	s2, _ = New("foo", Index("bar"), opt.Name("qux"))
	if seriesEquals(s, s2) {
		t.Errorf("seriesEquals() returned true for different name, want false")
	}
}
