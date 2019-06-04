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

// valueTypeValues is a slice of valueType-typed value/null structs.
type valueTypeValues []valueTypeValue

// valueTypeValue is a valueType-typed value/null struct.
type valueTypeValue struct {
	v    valueType
	null bool
}

// Len returns the number of value/null structs in the container.
func (vals valueTypeValues) Len() int {
	return len(vals)
}

// In returns the values located at specific index positions.
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

// Valid returns the integer position of all valid (i.e., non-null) values in the collection.
func (vals valueTypeValues) Valid() []int {
	var ret []int
	for i, val := range vals {
		if !val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Null returns the integer position of all null values in the collection.
func (vals valueTypeValues) Null() []int {
	var ret []int
	for i, val := range vals {
		if val.null {
			ret = append(ret, i)
		}
	}
	return ret
}

// Element returns a Value/Null pair at an integer position.
func (vals valueTypeValues) Element(position int) Elem {
	return Elem{vals[position].v, vals[position].null}
}

// Set overwrites a Value/Null pair at an integer position.
func (vals valueTypeValues) Set(position int, val interface{}) error {
	v := interfaceValue{val, false}
	if position >= vals.Len() {
		return fmt.Errorf("unable to set value at position %v: index out of range", position)
	}
	vals[position] = v.tovalueType()
	return nil
}

// ToFloat converts valueTypeValues to floatValues.
func (vals valueTypeValues) ToFloat() float64Values {
	var ret float64Values
	for _, val := range vals {
		ret = append(ret, val.toFloat64())
	}
	return ret
}

// ToInt converts valueTypeValues to intValues.
func (vals valueTypeValues) ToInt() int64Values {
	var ret int64Values
	for _, val := range vals {
		ret = append(ret, val.toInt64())
	}
	return ret
}

func (val valueTypeValue) toString() stringValue {
	if val.null {
		return stringValue{opt.GetDisplayStringNullFiller(), true}
	}
	return stringValue{fmt.Sprint(val.v), false}
}

// ToString converts valueTypeValues to stringValues.
func (vals valueTypeValues) ToString() stringValues {
	var ret stringValues
	for _, val := range vals {
		ret = append(ret, val.toString())
	}
	return ret
}

// ToBool converts valueTypeValues to boolValues.
func (vals valueTypeValues) ToBool() boolValues {
	var ret boolValues
	for _, val := range vals {
		ret = append(ret, val.toBool())
	}
	return ret
}

// ToBool converts valueTypeValues to dateTimeValues.
func (vals valueTypeValues) ToDateTime() dateTimeValues {
	var ret dateTimeValues
	for _, val := range vals {
		ret = append(ret, val.toDateTime())
	}
	return ret
}

// ToInterface converts valueTypeValues to interfaceValues.
func (vals valueTypeValues) ToInterface() interfaceValues {
	var ret interfaceValues
	for _, val := range vals {
		if val.null {
			ret = append(ret, interfaceValue{val.v, true})
		} else {
			ret = append(ret, interfaceValue{val.v, false})
		}
	}
	return ret
}

// [END] valueTypeValues
// ---------------------------------------------------------------------------
var placeholder = true

// the placeholder and this comment are overwritten on `make generate`, but are included so that the [END] comment survives
