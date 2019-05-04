package series

import (
	"reflect"
	"testing"

	"github.com/ptiger10/pd/kinds"
)

func TestConvert(t *testing.T) {
	var tests = []struct {
		convertTo kinds.Kind
	}{
		{kinds.Float},
		{kinds.Int},
		{kinds.String},
		{kinds.Bool},
		{kinds.DateTime},
	}
	for _, test := range tests {
		s, err := New([]interface{}{1, 2, 3})
		if err != nil {
			t.Error(err)
		}
		newS := s.As(test.convertTo)
		if newS.kind != test.convertTo {
			t.Errorf("Conversion of Series with []interface values to %v returned %v, want %v", test.convertTo, s.kind, test.convertTo)
		}
		if newS.Kind() == s.Kind() {
			t.Errorf("Conversion to %v occurred in place, want copy only", test.convertTo)
		}
	}
}

func TestConvertIndexDefault(t *testing.T) {
	var tests = []struct {
		convertTo kinds.Kind
	}{
		{kinds.Float},
		{kinds.Int},
		{kinds.String},
		{kinds.Bool},
		{kinds.DateTime},
	}
	for _, test := range tests {
		s, err := New([]interface{}{1, 2, 3})
		if err != nil {
			t.Error(err)
		}
		newS, err := s.IndexAs(test.convertTo)
		if err != nil {
			t.Error(err)
		}
		if newS.index.Levels[0].Kind != test.convertTo {
			t.Errorf("Conversion of Series with default []int64 index to %v returned %v, want %v", test.convertTo, s.index.Levels[0].Kind, test.convertTo)
		}
		if test.convertTo != kinds.Int {
			if s.index.Levels[0].Kind == newS.index.Levels[0].Kind {
				t.Errorf("Conversion to %v occurred in place, want copy only", test.convertTo)
			}
		}
	}
}

func TestConvertIndexMulti(t *testing.T) {
	var tests = []struct {
		convertTo kinds.Kind
		lvl       int
	}{
		{kinds.Float, 0},
		{kinds.Float, 1},
		{kinds.Int, 0},
		{kinds.Int, 1},
		{kinds.String, 0},
		{kinds.String, 1},
		{kinds.Bool, 0},
		{kinds.Bool, 1},
		{kinds.DateTime, 0},
		{kinds.DateTime, 1},
	}
	for _, test := range tests {
		s, err := New([]interface{}{1, 2, 3}, Index([]int{1, 2, 3}), Index([]int{10, 20, 30}))
		if err != nil {
			t.Error(err)
		}
		newS, err := s.IndexLevelAs(test.lvl, test.convertTo)
		if err != nil {
			t.Error(err)
		}
		if newS.index.Levels[test.lvl].Kind != test.convertTo {
			t.Errorf("Conversion of Series with multiIndex level %v to %v returned %v, want %v", test.lvl, test.convertTo, newS.index.Levels[test.lvl].Kind, test.convertTo)
		}
		// excludes Int because the original test Index is int
		if test.convertTo != kinds.Int {
			if s.index.Levels[test.lvl].Kind == newS.index.Levels[test.lvl].Kind {
				t.Errorf("Conversion to %v occurred in place, want copy only", test.convertTo)
			}
		}
	}
}

func TestConvertUnsupported(t *testing.T) {
	s, err := New([]interface{}{1, 2, 3})
	if err != nil {
		t.Error(err)
	}
	newS := s.As(kinds.Unsupported)
	if !reflect.DeepEqual(s, newS) {
		t.Errorf("Unsupported conversion returned %v, want %v", newS, s)
	}
}

func TestConvertUnsupportedIndex(t *testing.T) {
	s, err := New([]interface{}{1, 2, 3})
	if err != nil {
		t.Error(err)
	}
	_, err = s.IndexAs(kinds.Unsupported)
	if err == nil {
		t.Errorf("Returned nil error, expected error due to unsupported conversion")
	}
}

func TestConvertUnsupportedIndexLevel(t *testing.T) {
	s, err := New([]interface{}{1, 2, 3}, Index([]int{1, 2, 3}), Index([]int{10, 20, 30}))
	if err != nil {
		t.Error(err)
	}
	_, err = s.IndexLevelAs(0, kinds.Unsupported)
	if err == nil {
		t.Errorf("Returned nil error, expected error due to unsupported conversion")
	}
	_, err = s.IndexLevelAs(5, kinds.Float)
	if err == nil {
		t.Errorf("Returned nil error, expected error due to out-of-range conversion")
	}
}
