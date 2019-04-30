package values

import (
	"fmt"
	"log"
)

// Data Type
// ------------------------------------------------
type StringValues []StringValue
type StringValue struct {
	V    string
	Null bool
}

// Convenience methods
// ------------------------------------------------
func (vals StringValues) count() int {
	var count int
	for _, val := range vals {
		if !val.Null {
			count++
		}
	}
	return count
}

func (vals StringValues) Len() int {
	return len(vals)
}

func (vals StringValues) Describe() string {
	offset := 7
	l := len(vals)
	len := fmt.Sprintf("%-*s %d\n", offset, "len", l)
	return fmt.Sprint(len)
}

func (vals StringValues) In(positions []int) interface{} {
	var ret StringValues
	for _, position := range positions {
		if position >= len(vals) {
			log.Panicf("Unable to get value: index out of range: %d", position)
		}
		ret = append(ret, vals[position])
	}
	return ret
}
