package constructors

import (
	"github.com/ptiger10/pd/internal/values"
)

// [START Utilities]

// [END Utilities]

// [START Constructor Functions]

// SliceInt converts []int (of any variety) -> values.IntValues
func SliceInt(vals []int64) values.IntValues {
	var v values.IntValues
	for _, val := range vals {
		v = append(v, values.Int(val, false))
	}
	return v
}

// [END Constructor Functions]
