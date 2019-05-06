package constructors

import (
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/kinds"
)

// [START Constructor Functions]

// SliceBool converts []bool -> ValuesFactory with values.BoolValues
func SliceBool(vals []bool) ValuesFactory {
	var v values.BoolValues
	for _, val := range vals {
		v = append(v, values.Bool(val, false))
	}
	return ValuesFactory{v, kinds.Bool}
}

// [END Constructor Functions]
