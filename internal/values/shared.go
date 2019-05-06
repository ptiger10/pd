package values

import (
	"fmt"
	"log"

	"github.com/cheekybits/genny/generic"
	"github.com/ptiger10/pd/options"
)

//go:generate genny -in=$GOFILE -out=shared_autogen.go gen "valueType=float64,int64,string,bool,time.Time,interface{}"

// [START] valueTypeValues

// valueType is the generic ValueType that will be replaced by specific types on `make generate`
type valueType generic.Type

// valueTypeValues is a slice of ValueType-typed value/null structs
type valueTypeValues []valueTypeValue

// valueTypeValue is a ValueType-typed value/null struct
type valueTypeValue struct {
	v    valueType
	null bool
}

// valueTypeVal constructs a valueTypeValue value/null struct
func valueTypeVal(v valueType, null bool) valueTypeValue {
	return valueTypeValue{
		v:    v,
		null: null,
	}
}

// Len returns the number of value/null structs in the container
func (vals valueTypeValues) Len() int {
	return len(vals)
}

// In returns the values located at specific index positions
func (vals valueTypeValues) In(positions []int) Values {
	var ret valueTypeValues
	for _, position := range positions {
		if position >= len(vals) {
			log.Panicf("Unable to get value: index out of range: %d", position)
		}
		ret = append(ret, vals[position])
	}
	return ret
}

// All returns all of the values in the collection (including null values) as an interface slice
func (vals valueTypeValues) All() []interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Vals returns only the Value fields for the collection of Value/Null structs.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals valueTypeValues) Vals() interface{} {
	var ret []valueType
	for _, val := range vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Valid returns the integer position of all valid (i.e., non-null) values in the collection
func (vals valueTypeValues) Valid() []int {
	var ret []int
	for i, val := range vals {
		if !val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Null returns the integer position of all null values in the collection
func (vals valueTypeValues) Null() []int {
	var ret []int
	for i, val := range vals {
		if val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Element returns a value element at an integer position in form
// []interface{val, null}
func (vals valueTypeValues) Element(position int) []interface{} {
	return []interface{}{vals[position].v, vals[position].null}
}

// ToString converts the values to stringValues
func (vals valueTypeValues) ToString() Values {
	var ret stringValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, stringVal(options.DisplayStringNullFiller, true))
		} else {
			ret = append(ret, stringVal(fmt.Sprint(val.v), false))
		}
	}
	return ret
}

// ToInterface converts teh values to interfaceValues
func (vals valueTypeValues) ToInterface() Values {
	var ret interfaceValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, interfaceVal(val.v, true))
		} else {
			ret = append(ret, interfaceVal(val.v, false))
		}
	}
	return ret
}

// [END] valueTypeValues
var placeholder = true
