package constructors

import "github.com/ptiger10/pd/new/internal/values"

// [START Constructor Functions]

// SliceBool converts []bool -> values.BoolValues
func SliceBool(vals []bool) values.BoolValues {
	var v values.BoolValues
	for _, val := range vals {
		v = append(v, values.Bool(val, false))
	}
	return v
}

// [END Constructor Functions]
