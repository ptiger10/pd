package values

import (
	"fmt"
	"math"
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

func TestSet(t *testing.T) {
	vals := float64Values([]float64Value{
		float64Value{1, false},
	})
	vals.Set(0, "foo")
	fmt.Println(vals)
}
