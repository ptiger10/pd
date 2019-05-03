package constructors

import (
	"math"

	"github.com/ptiger10/pd/internal/values"
)

// [START Constructor Functions]

// SliceFloat converts []float64  -> values.FloatValues
func SliceFloat(vals []float64) values.FloatValues {
	var v values.FloatValues
	for _, val := range vals {
		if math.IsNaN(val) {
			v = append(v, values.Float(val, true))
			continue
		}
		v = append(v, values.Float(val, false))
	}
	return v
}

// [END Constructor Functions]
