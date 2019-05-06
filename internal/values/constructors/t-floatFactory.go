package constructors

import (
	"math"

	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/kinds"
)

// [START Constructor Functions]

// SliceFloat converts []float64  -> ValuesFactory with FloatValues
func SliceFloat(vals []float64) ValuesFactory {
	var v values.FloatValues
	for _, val := range vals {
		if math.IsNaN(val) {
			v = append(v, values.Float(val, true))
			continue
		}
		v = append(v, values.Float(val, false))
	}
	return ValuesFactory{v, kinds.Float}
}

// [END Constructor Functions]
