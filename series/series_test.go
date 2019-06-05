package series

import (
	"reflect"
	"testing"

	"github.com/ptiger10/pd/kinds"
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
