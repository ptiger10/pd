package values

import (
	"log"

	"github.com/ptiger10/pd/options"
)

// warning is the standardized part of any warning message shown to a caller due to bad parameters.
var warning = "Warning! %v. Returning %v."

// Warn prints a warning to the caller when they have triggered a fatal error,
// and a description of what is being returned instead of what they expected.
func Warn(err error, returnDesc string) {
	if options.GetLogWarnings() {
		log.Printf(warning, err, returnDesc)
		return
	}
	return
}

// MakeRange returns a sequential series of numbers, for use in the default Series index constructor
func MakeRange(min, max int) []int64 {
	a := make([]int64, max-min)
	for i := range a {
		a[i] = int64(min + i)
	}
	return a
}

// MakeIntRange returns a sequential series of numbers, for use in creating a list of integer positions
func MakeIntRange(min, max int) []int {
	a := make([]int, max-min)
	for i := range a {
		a[i] = min + i
	}
	return a
}

// NewDefault returns a factory of []int64 values {0, 1, 2, ... n} for use in a default index.
func NewDefault(n int) Values {
	defaultRange := MakeRange(0, n)
	v := newSliceInt64(defaultRange)
	return v.Values
}
