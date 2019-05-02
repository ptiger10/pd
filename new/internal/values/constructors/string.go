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
func SliceString(data interface{}) values.StringValues {
	var vals values.StringValues
	d := data.([]string)
	for i := 0; i < len(d); i++ {
		val := d[i]
		if isNullString(val) {
			vals = append(vals, values.String("", true))
		} else {
			vals = append(vals, values.String(val, false))
		}

	}
	return vals
}

// [END Constructor Functions]
