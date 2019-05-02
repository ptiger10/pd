package constructors

import (
	"strings"

	"github.com/ptiger10/pd/new/internal/values"
)

// [START Convenience Functions]

func isNullString(s string) bool {
	nullStrings := []string{"NaN", "n/a", "N/A", "", "nil"}
	for _, ns := range nullStrings {
		if strings.TrimSpace(s) == ns {
			return true
		}
	}
	return false
}

// [END Convenience Functions]

// [START Constructor Functions]

// SliceString converts []string -> values.StringValues
func SliceString(vals []string) values.StringValues {
	var v values.StringValues
	for _, val := range vals {
		if isNullString(val) {
			v = append(v, values.String("", true))
		} else {
			v = append(v, values.String(val, false))
		}

	}
	return v
}

// [END Constructor Functions]
