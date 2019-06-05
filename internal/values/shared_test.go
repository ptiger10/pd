package values

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestValid(t *testing.T) {
	vals := float64Values([]float64Value{
		float64Value{1, false}, float64Value{math.NaN(), true},
	})
	at, _ := vals.In(vals.Valid())
	fmt.Println(at.Vals().([]float64))
}

// Set

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
