package index_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/ptiger10/pd/new/internal/index"

	constructIdx "github.com/ptiger10/pd/new/internal/index/constructors"
	"github.com/ptiger10/pd/new/kinds"
)

func TestConvertIndex_int(t *testing.T) {
	var tests = []struct {
		lvl       index.Level
		convertTo reflect.Kind
	}{

		{constructIdx.SliceFloat([]float64{1, 2, 3}, ""), kinds.Float},
		{constructIdx.SliceFloat([]float64{1, 2, 3}, ""), kinds.Int},
		{constructIdx.SliceFloat([]float64{1, 2, 3}, ""), kinds.String},
		{constructIdx.SliceFloat([]float64{1, 2, 3}, ""), kinds.Bool},
		{constructIdx.SliceFloat([]float64{1, 2, 3}, ""), kinds.DateTime},

		{constructIdx.SliceInt([]int{1, 2, 3}, ""), kinds.Float},
		{constructIdx.SliceInt([]int{1, 2, 3}, ""), kinds.Int},
		{constructIdx.SliceInt([]int{1, 2, 3}, ""), kinds.String},
		{constructIdx.SliceInt([]int{1, 2, 3}, ""), kinds.Bool},
		{constructIdx.SliceInt([]int{1, 2, 3}, ""), kinds.DateTime},

		{constructIdx.SliceString([]string{"1", "2", "3"}, ""), kinds.Float},
		{constructIdx.SliceString([]string{"1", "2", "3"}, ""), kinds.Int},
		{constructIdx.SliceString([]string{"1", "2", "3"}, ""), kinds.String},
		{constructIdx.SliceString([]string{"1", "2", "3"}, ""), kinds.Bool},
		{constructIdx.SliceString([]string{"1", "2", "3"}, ""), kinds.DateTime},

		{constructIdx.SliceBool([]bool{true, false, false}, ""), kinds.Float},
		{constructIdx.SliceBool([]bool{true, false, false}, ""), kinds.Int},
		{constructIdx.SliceBool([]bool{true, false, false}, ""), kinds.String},
		{constructIdx.SliceBool([]bool{true, false, false}, ""), kinds.Bool},
		{constructIdx.SliceBool([]bool{true, false, false}, ""), kinds.DateTime},

		{constructIdx.SliceDateTime([]time.Time{time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)}, ""), kinds.Float},
		{constructIdx.SliceDateTime([]time.Time{time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)}, ""), kinds.Int},
		{constructIdx.SliceDateTime([]time.Time{time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)}, ""), kinds.String},
		{constructIdx.SliceDateTime([]time.Time{time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)}, ""), kinds.Bool},
		{constructIdx.SliceDateTime([]time.Time{time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)}, ""), kinds.DateTime},

		{constructIdx.SliceInterface([]interface{}{1, "2", true}, ""), kinds.Float},
		{constructIdx.SliceInterface([]interface{}{1, "2", true}, ""), kinds.Int},
		{constructIdx.SliceInterface([]interface{}{1, "2", true}, ""), kinds.String},
		{constructIdx.SliceInterface([]interface{}{1, "2", true}, ""), kinds.Bool},
		{constructIdx.SliceInterface([]interface{}{1, "2", true}, ""), kinds.DateTime},
	}
	for _, test := range tests {
		lvl, err := test.lvl.Convert(test.convertTo)
		if err != nil {
			t.Error(err)
		}
		if lvl.Kind != test.convertTo {
			t.Errorf("Attempted conversion to %v returned %v", test.convertTo, lvl.Kind)
		}
	}
}

func TestConvert_Numeric_Datetime(t *testing.T) {
	n := 1556668800000000000
	wantVal := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	var tests = []struct {
		lvl index.Level
	}{
		{constructIdx.SliceInt([]int{n}, "")},
		{constructIdx.SliceFloat([]float64{float64(n)}, "")},
	}
	for _, test := range tests {
		lvl, _ := test.lvl.Convert(kinds.DateTime)
		elem := lvl.Labels.Element(0)
		gotVal := elem[0].(time.Time)
		if gotVal != wantVal {
			t.Errorf("Error converting %v to datetime: returned %v, want %v", test.lvl, gotVal, wantVal)
		}
	}
}

func TestConvert_Unsupported(t *testing.T) {
	var tests = []struct {
		kind reflect.Kind
	}{
		{kinds.None},
		{reflect.Complex64},
	}
	for _, test := range tests {
		lvl := constructIdx.SliceFloat([]float64{1, 2, 3}, "")
		_, err := lvl.Convert(test.kind)
		if err == nil {
			t.Errorf("Returned nil error, expected error due to unsupported type %v", test.kind)
		}
	}
}
