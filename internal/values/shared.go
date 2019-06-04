package values

import (
	"fmt"

	"github.com/cheekybits/genny/generic"
	"github.com/ptiger10/pd/opt"
)

//go:generate genny -in=$GOFILE -out=shared_autogen.go gen "valueType=float64,int64,string,bool,time.Time,interface{}"

// [START] valueTypeValues

// valueType is the generic ValueType that will be replaced by specific types on `make generate`
type valueType generic.Type

// valueTypeValues is a slice of valueType-typed value/null structs
type valueTypeValues []valueTypeValue

// valueTypeValue is a valueType-typed value/null struct
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
func (vals valueTypeValues) In(rowPositions []int) (Values, error) {
	var ret valueTypeValues
	for _, position := range rowPositions {
		if position >= len(vals) {
			return nil, fmt.Errorf("%d is not a valid integer position (len: %v)", position, len(vals))
		}
		ret = append(ret, vals[position])
	}
	return ret, nil
}

// All returns only the Value fields for the collection of Value/Null structs as an interface slice.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals valueTypeValues) All() []interface{} {
	var ret []interface{}
	for _, val := range vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Vals returns only the Value fields for the collection of Value/Null structs as an empty interface.
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
func (vals valueTypeValues) Element(position int) Elem {
	return Elem{vals[position].v, vals[position].null}
}

// ToString converts the values to stringValues
func (vals valueTypeValues) ToString() Values {
	var ret stringValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, stringVal(opt.GetDisplayStringNullFiller(), true))
		} else {
			ret = append(ret, stringVal(fmt.Sprint(val.v), false))
		}
	}
	return ret
}

// ToInterface converts the values to interfaceValues
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
// ---------------------------------------------------------------------------
var placeholder = true

// the placeholder and this comment are overwritten on `make generate`, but are included so that the [END] comment survives
