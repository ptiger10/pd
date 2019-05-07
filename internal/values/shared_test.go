package values

import (
	"fmt"
	"math"
	"testing"
)

func TestValid(t *testing.T) {
	vals := float64Values([]float64Value{
		float64Val(1, false), float64Val(math.NaN(), true),
	})
	at, _ := vals.In(vals.Valid())
	fmt.Println(at.Vals().([]float64))
}
