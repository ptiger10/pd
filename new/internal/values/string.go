package values

import (
	"fmt"
	"log"
)

// Data Type
type StringValues []StringValue
type StringValue struct {
	v    string
	null bool
}

func String(v string, null bool) StringValue {
	return StringValue{
		v:    v,
		null: null,
	}
}

// Convenience methods
// ------------------------------------------------
func (vals StringValues) count() int {
	var count int
	for _, val := range vals {
		if !val.null {
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
