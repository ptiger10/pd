package values

import (
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/ptiger10/pd/kinds"
)

func TestSliceConstructor(t *testing.T) {
	testDate := time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)
	var tests = []struct {
		data     interface{}
		wantVals Factory
		wantKind kinds.Kind
	}{
		{
			data:     []float32{0, 1, 2},
			wantVals: newSliceFloat64([]float64{0, 1, 2}),
			wantKind: kinds.Float64,
		},
		{
			data:     []float64{0, 1, 2},
			wantVals: newSliceFloat64([]float64{0, 1, 2}),
			wantKind: kinds.Float64,
		},
		{
			data:     []float64{},
			wantVals: newSliceFloat64([]float64{}),
			wantKind: kinds.Float64,
		},
		{
			data:     []int{0, 1, 2},
			wantVals: newSliceInt64([]int64{0, 1, 2}),
			wantKind: kinds.Int64,
		},
		{
			data:     []int64{0, 1, 2},
			wantVals: newSliceInt64([]int64{0, 1, 2}),
			wantKind: kinds.Int64,
		},
		{
			data:     []int{},
			wantVals: newSliceInt64([]int64{}),
			wantKind: kinds.Int64,
		},
		{
			data:     []uint{0, 1, 2},
			wantVals: newSliceInt64([]int64{0, 1, 2}),
			wantKind: kinds.Int64,
		},
		{
			data:     []string{"0", "1", ""},
			wantVals: newSliceString([]string{"0", "1", "NaN"}),
			wantKind: kinds.String,
		},
		{
			data:     []string{},
			wantVals: newSliceString([]string{}),
			wantKind: kinds.String,
		},
		{
			data:     []bool{true, true, false},
			wantVals: newSliceBool([]bool{true, true, false}),
			wantKind: kinds.Bool,
		},
		{
			data:     []bool{},
			wantVals: newSliceBool([]bool{}),
			wantKind: kinds.Bool,
		},
		{
			data:     []time.Time{testDate, time.Time{}},
			wantVals: newSliceDateTime([]time.Time{testDate, time.Time{}}),
			wantKind: kinds.DateTime,
		},
		{
			data:     []time.Time{},
			wantVals: newSliceDateTime([]time.Time{}),
			wantKind: kinds.DateTime,
		},
		{
			data:     []interface{}{1.5, 1, "", false, testDate},
			wantVals: newSliceInterface([]interface{}{1.5, 1, "", false, testDate}),
			wantKind: kinds.Interface,
		},
		{
			data:     []interface{}{},
			wantVals: newSliceInterface([]interface{}{}),
			wantKind: kinds.Interface,
		},
	}
	for _, test := range tests {
		vals, err := SliceFactory(test.data)
		if err != nil {
			t.Errorf("Unable to construct values from %v: %v", test.data, err)
		}
		if !reflect.DeepEqual(vals.Values, test.wantVals.Values) {
			t.Errorf("%T test returned values %#v, want %#v", test.data, vals, test.wantVals.Values)
		}
		if vals.Kind != test.wantKind {
			t.Errorf("%T test returned value %v, want %v", test.data, vals.Kind, test.wantKind)
		}
	}
}

func TestSliceConstructor_NullFloat(t *testing.T) {
	vals, err := SliceFactory([]float64{math.NaN()})
	if err != nil {
		t.Errorf("Unable to construct values from null float: %v", err)
	}
	val := vals.Values.Element(0).Value.(float64)
	if !math.IsNaN(val) {
		t.Errorf("Returned %v, want NaN", val)
	}
}

func TestSliceConstructor_NullFloatInterface(t *testing.T) {
	vals, err := SliceFactory([]interface{}{math.NaN()})
	if err != nil {
		t.Errorf("Unable to construct values from null float: %v", err)
	}
	val := vals.Values.Element(0).Value.(float64)
	if !math.IsNaN(val) {
		t.Errorf("Returned %v, want NaN", val)
	}
}

func TestSliceConstructor_Unsupported(t *testing.T) {
	data := []complex64{1, 2, 3}
	_, err := SliceFactory(data)
	if err == nil {
		t.Errorf("Returned nil error, expected error due to unsupported type %T", data)
	}
}
