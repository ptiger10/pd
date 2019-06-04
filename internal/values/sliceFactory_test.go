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
			wantVals: SliceFloat([]float64{0, 1, 2}),
			wantKind: kinds.Float,
		},
		{
			data:     []float64{0, 1, 2},
			wantVals: SliceFloat([]float64{0, 1, 2}),
			wantKind: kinds.Float,
		},
		{
			data:     []float64{},
			wantVals: SliceFloat([]float64{}),
			wantKind: kinds.Float,
		},
		{
			data:     []int{0, 1, 2},
			wantVals: SliceInt([]int64{0, 1, 2}),
			wantKind: kinds.Int,
		},
		{
			data:     []int64{0, 1, 2},
			wantVals: SliceInt([]int64{0, 1, 2}),
			wantKind: kinds.Int,
		},
		{
			data:     []int{},
			wantVals: SliceInt([]int64{}),
			wantKind: kinds.Int,
		},
		{
			data:     []uint{0, 1, 2},
			wantVals: SliceInt([]int64{0, 1, 2}),
			wantKind: kinds.Int,
		},
		{
			data:     []string{"0", "1", ""},
			wantVals: SliceString([]string{"0", "1", "NaN"}),
			wantKind: kinds.String,
		},
		{
			data:     []string{},
			wantVals: SliceString([]string{}),
			wantKind: kinds.String,
		},
		{
			data:     []bool{true, true, false},
			wantVals: SliceBool([]bool{true, true, false}),
			wantKind: kinds.Bool,
		},
		{
			data:     []bool{},
			wantVals: SliceBool([]bool{}),
			wantKind: kinds.Bool,
		},
		{
			data:     []time.Time{testDate, time.Time{}},
			wantVals: SliceDateTime([]time.Time{testDate, time.Time{}}),
			wantKind: kinds.DateTime,
		},
		{
			data:     []time.Time{},
			wantVals: SliceDateTime([]time.Time{}),
			wantKind: kinds.DateTime,
		},
		{
			data:     []interface{}{1.5, 1, "", false, testDate},
			wantVals: SliceInterface([]interface{}{1.5, 1, "", false, testDate}),
			wantKind: kinds.Interface,
		},
		{
			data:     []interface{}{},
			wantVals: SliceInterface([]interface{}{}),
			wantKind: kinds.Interface,
		},
	}
	for _, test := range tests {
		vals, err := SliceFactory(test.data)
		if err != nil {
			t.Errorf("Unable to construct values from %v: %v", test.data, err)
		}
		if !reflect.DeepEqual(vals.V, test.wantVals.V) {
			t.Errorf("%T test returned values %#v, want %#v", test.data, vals, test.wantVals.V)
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
	val := vals.V.Element(0).Value.(float64)
	if !math.IsNaN(val) {
		t.Errorf("Returned %v, want NaN", val)
	}
}

func TestSliceConstructor_NullFloatInterface(t *testing.T) {
	vals, err := SliceFactory([]interface{}{math.NaN()})
	if err != nil {
		t.Errorf("Unable to construct values from null float: %v", err)
	}
	val := vals.V.Element(0).Value.(float64)
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
