package values

import (
	"log"

	"github.com/ptiger10/pd/options"
)

// warning is the standardized part of any warning message shown to a caller due to bad parameters.
var warning = "Warning: caller has triggered a fatal error: %v. Returning %v instead."

// Warn prints a warning to the caller when they have triggered a fatal error,
// and a description of what is being returned instead of what they expected.
func Warn(err error, returnDesc string) {
	if options.GetLogWarnings() {
		log.Printf(warning, err, returnDesc)
		return
	}
	return
}

// MakeRange returns a sequential series of numbers, for use in default constructors
func MakeRange(min, max int) []int64 {
	a := make([]int64, max-min)
	for i := range a {
		a[i] = int64(min + i)
	}
	return a
}

// MakeIntRange returns a sequential series of numbers, for use in position trackers
func MakeIntRange(min, max int) []int {
	a := make([]int, max-min)
	for i := range a {
		a[i] = min + i
	}
	return a
}
