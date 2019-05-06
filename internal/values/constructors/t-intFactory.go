package constructors

import (
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/kinds"
)

// [START Utilities]

// [END Utilities]

// [START Constructor Functions]

// SliceInt converts []int (of any variety) -> ValuesFactory with values.IntValues
func SliceInt(vals []int64) ValuesFactory {
	var v values.IntValues
	for _, val := range vals {
		v = append(v, values.Int(val, false))
	}
	return ValuesFactory{v, kinds.Int}
}

// [END Constructor Functions]
