package constructors

import (
	"math"
	"reflect"

	"github.com/ptiger10/pd/new/internal/values"
)

// [START Constructor Functions]

// SliceFloat converts []float (of any variety) -> values.FloatValues
func SliceFloat(data interface{}) values.FloatValues {
	var vals values.FloatValues
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		val := d.Index(i).Float()
		if math.IsNaN(val) {
			vals = append(vals, values.Float(val, true))
			continue
		}
		vals = append(vals, values.Float(val, false))
	}
	return vals
}

// [END Constructor Functions]
