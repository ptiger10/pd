package values

import (
	"fmt"
	"math"
	"testing"
)

func TestValid(t *testing.T) {
	vals := FloatValues([]FloatValue{
		Float(1, false), Float(math.NaN(), true),
	})
	at := vals.In(vals.Valid()).Vals().([]float64)
	fmt.Println(at)
}
