package series

import (
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
		s = s.As(test.convertTo)
		if s.Kind != test.convertTo {
			t.Errorf("Conversion of []interface Series to %v returned %v, want %v", test.convertTo, s.Kind, test.convertTo)
		}
	}
}

func TestConvertIndex(t *testing.T) {
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
		s = s.IndexAs(test.convertTo)
		if s.index.Levels[0].Kind != test.convertTo {
			t.Errorf("Conversion of default []int64 Series to %v returned %v, want %v", test.convertTo, s.Kind, test.convertTo)
		}
	}
}
