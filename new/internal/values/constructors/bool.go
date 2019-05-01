package constructors

import "github.com/ptiger10/pd/new/internal/values"

// [START Constructor Functions]

// SliceBool converts []bool -> values.BoolValues
func SliceBool(data interface{}) values.BoolValues {
	var vals values.BoolValues
	d := data.([]bool)
	for i := 0; i < len(d); i++ {
		val := d[i]
		vals = append(vals, values.Bool(val, false))

	}
	return vals
}

// [END Constructor Functions]
