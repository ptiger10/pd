package values

import (
	"math"
	"testing"
)

func TestValid(t *testing.T) {
	vals := FloatValues([]FloatValue{
		Float(1, false), Float(math.NaN(), true),
	})
	vals.In(vals.Valid()).All().(float64)
}
