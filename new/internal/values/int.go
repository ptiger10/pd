package values

import (
	"fmt"
	"log"
	"reflect"

	"github.com/ptiger10/pd/new/kinds"
)

// Data Type
type IntValues []IntValue
type IntValue struct {
	V    int64
	Null bool
}

func Int(v int64, null bool) IntValue {
	return IntValue{
		V:    v,
		Null: null,
	}
}

// Convenience methods
// ------------------------------------------------
func (vals IntValues) count() int {
	var count int
	for _, val := range vals {
		if !val.Null {
			count++
		}
	}
	return count
}

func (vals IntValues) Len() int {
	return len(vals)
}

func (vals IntValues) Describe() string {
	offset := 7
	l := len(vals)
	len := fmt.Sprintf("%-*s %d\n", offset, "len", l)
	return fmt.Sprint(len)
}

func (vals IntValues) In(positions []int) interface{} {
	var ret IntValues
	for _, position := range positions {
		if position >= len(vals) {
			log.Panicf("Unable to get value: index out of range: %d", position)
		}
		ret = append(ret, vals[position])
	}
	return ret
}

func (vals IntValues) Kind() reflect.Kind {
	return kinds.String
}

func (vals IntValues) At(position int) interface{} {
	return vals[position].V
}

func (vals IntValues) All() []interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.V)
	}
	return ret
}
