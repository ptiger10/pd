package series

import (
	"reflect"
	"testing"

	"github.com/ptiger10/pd/new/kinds"
)

func TestElement(t *testing.T) {
	s, err := New([]string{"", "valid"}, Index([]string{"A", "B"}), Index([]int{1, 2}))
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
	wantKind := kinds.String
	wantIdxKinds := []kinds.Kind{kinds.String, kinds.Int}
	for _, test := range tests {
		got := s.Elem(test.position)
		if got.Value != test.wantVal {
			t.Errorf("Element returned value %v, want %v", got.Value, test.wantVal)
		}
		if got.Null != test.wantNull {
			t.Errorf("Element returned bool %v, want %v", got.Null, test.wantNull)
		}
		if !reflect.DeepEqual(got.Index, test.wantIdx) {
			t.Errorf("Element returned index %#v, want %#v", got.Index, test.wantIdx)
		}
		if got.Kind != wantKind {
			t.Errorf("Element returned kind %v, want %v", got.Kind, wantKind)
		}
		if !reflect.DeepEqual(got.IndexKinds, wantIdxKinds) {
			t.Errorf("Element returned kind %v, want %v", got.IndexKinds, wantIdxKinds)
		}
	}
}
