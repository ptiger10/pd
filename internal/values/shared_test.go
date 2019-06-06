package values

import (
	"math"
	"reflect"
	"testing"
)

func TestDrop_float(t *testing.T) {
	var tests = []struct {
		dropPosition int
		want         float64Values
	}{
		{0, float64Values([]float64Value{float64Value{2, false}, float64Value{3, false}})},
		{1, float64Values([]float64Value{float64Value{1, false}, float64Value{3, false}})},
		{2, float64Values([]float64Value{float64Value{1, false}, float64Value{2, false}})},
	}
	for _, test := range tests {
		vals := float64Values([]float64Value{float64Value{1, false}, float64Value{2, false}, float64Value{3, false}})
		err := vals.Drop(test.dropPosition)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(vals, test.want) {
			t.Errorf("float64Values.Drop() returned %v, want %v", vals, test.want)
		}
	}
}

func TestInsert_float(t *testing.T) {
	var tests = []struct {
		insertPosition int
		val            interface{}
		want           float64Values
	}{
		{0, 10, float64Values([]float64Value{float64Value{10, false}, float64Value{1, false}, float64Value{2, false}})},
		{1, 10, float64Values([]float64Value{float64Value{1, false}, float64Value{10, false}, float64Value{2, false}})},
		{2, 10, float64Values([]float64Value{float64Value{1, false}, float64Value{2, false}, float64Value{10, false}})},
	}
	for _, test := range tests {
		vals := float64Values([]float64Value{float64Value{1, false}, float64Value{2, false}})
		err := vals.Insert(test.insertPosition, test.val)
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(vals, test.want) {
			t.Errorf("float64Values.Insert() returned %v, want %v", vals, test.want)
		}
	}
}

func TestSet_float(t *testing.T) {
	vals := float64Values([]float64Value{float64Value{1, false}})
	err := vals.Set(0, "foo")
	if err != nil {
		t.Errorf("values.Set() returned err: %v", err)
	}
	got := vals.Element(0)
	want := float64Value{math.NaN(), true}
	if !math.IsNaN(got.Value.(float64)) {
		t.Errorf("values.Set() returned %v for value, want NaN", got.Value)
	}
	if !reflect.DeepEqual(got.Null, want.null) {
		t.Errorf("values.Set() returned %v for null value, want %v", got.Null, want.null)
	}

}

func TestCopy_float(t *testing.T) {
	vals := float64Values([]float64Value{float64Value{1, false}})
	copyVals := vals.Copy()
	if reflect.ValueOf(vals).Pointer() == reflect.ValueOf(copyVals).Pointer() {
		t.Errorf("values.Copy() returned original Values, want fresh copy")
	}
}
