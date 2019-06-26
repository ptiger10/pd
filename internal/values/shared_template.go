package values

import (
	"fmt"

	"github.com/cheekybits/genny/generic"
	"github.com/ptiger10/pd/options"
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

// newSlicevalueType converts []valueType -> Factory with valueTypeValues
func newSlicevalueType(vals []valueType) Factory {
	var ret valueTypeValues
	for _, val := range vals {
		ret = append(ret, newvalueType(val))
	}
	return Factory{&ret, options.PlaceholdervalueType}
}

// Len returns the number of value/null structs in the container.
func (vals *valueTypeValues) Len() int {
	return len(*vals)
}

func (vals *valueTypeValues) Swap(i, j int) {
	(*vals)[i], (*vals)[j] = (*vals)[j], (*vals)[i]
}

// Subset returns the values located at specific index positions.
func (vals *valueTypeValues) Subset(rowPositions []int) (Values, error) {
	var ret valueTypeValues
	for _, position := range rowPositions {
		if position >= len(*vals) {
			return nil, fmt.Errorf("invalid row position %d (max: %v)", position, len(*vals)-1)
		}
		ret = append(ret, (*vals)[position])
	}
	return &ret, nil
}

// Vals returns only the Value fields for the collection of Value/Null structs as an empty interface.
//
// Caution: This operation excludes the Null field but retains any null values.
func (vals *valueTypeValues) Vals() interface{} {
	var ret []valueType
	for _, val := range *vals {
		ret = append(ret, val.v)
	}
	return ret
}

// Element returns a Value/Null pair at an integer position.
func (vals *valueTypeValues) Element(position int) Elem {
	return Elem{(*vals)[position].v, (*vals)[position].null}
}

// Copy transfers every value from the current valueTypeValues container into a new Values container
func (vals *valueTypeValues) Copy() Values {
	newValues := valueTypeValues{}
	for _, val := range *vals {
		newValues = append(newValues, val)
	}
	return &newValues
}

// Set overwrites a Value/Null pair at an integer position.
func (vals *valueTypeValues) Set(position int, newVal interface{}) error {
	if position >= len(*vals) {
		return fmt.Errorf("unable to set value at position %v: index out of range", position)
	}
	v := interfaceValue{newVal, false}
	(*vals)[position] = v.tovalueType()
	return nil
}

// Drop drops the Value/Null pair at an integer position.
func (vals *valueTypeValues) Drop(pos int) error {
	if pos >= len(*vals) {
		return fmt.Errorf("unable to drop value at position %v: index out of range", pos)
	}
	*vals = append((*vals)[:pos], (*vals)[pos+1:]...)
	return nil
}

// Insert inserts a new Value/Null pair at an integer position.
func (vals *valueTypeValues) Insert(pos int, val interface{}) error {
	if pos > len(*vals) {
		return fmt.Errorf("unable to insert value at position %v: index out of range", pos)
	}
	v := interfaceValue{val, false}
	*vals = append((*vals)[:pos], append([]valueTypeValue{v.tovalueType()}, (*vals)[pos:]...)...)
	return nil
}

// ToFloat converts valueTypeValues to floatValues.
func (vals *valueTypeValues) ToFloat64() Values {
	var ret float64Values
	for _, val := range *vals {
		ret = append(ret, val.toFloat64())
	}
	return &ret
}

// ToInt converts valueTypeValues to intValues.
func (vals *valueTypeValues) ToInt64() Values {
	var ret int64Values
	for _, val := range *vals {
		ret = append(ret, val.toInt64())
	}
	return &ret
}

func (val *valueTypeValue) toString() stringValue {
	if val.null {
		return stringValue{options.GetDisplayStringNullFiller(), true}
	}
	return stringValue{fmt.Sprint(val.v), false}
}

// ToString converts valueTypeValues to stringValues.
func (vals *valueTypeValues) ToString() Values {
	var ret stringValues
	for _, val := range *vals {
		ret = append(ret, val.toString())
	}
	return &ret
}

// ToBool converts valueTypeValues to boolValues.
func (vals *valueTypeValues) ToBool() Values {
	var ret boolValues
	for _, val := range *vals {
		ret = append(ret, val.toBool())
	}
	return &ret
}

// ToBool converts valueTypeValues to dateTimeValues.
func (vals *valueTypeValues) ToDateTime() Values {
	var ret dateTimeValues
	for _, val := range *vals {
		ret = append(ret, val.toDateTime())
	}
	return &ret
}

// ToInterface converts valueTypeValues to interfaceValues.
func (vals *valueTypeValues) ToInterface() Values {
	var ret interfaceValues
	for _, val := range *vals {
		if val.null {
			ret = append(ret, interfaceValue{val.v, true})
		} else {
			ret = append(ret, interfaceValue{val.v, false})
		}
	}
	return &ret
}

// [END] valueTypeValues
// ---------------------------------------------------------------------------
var placeholder = true

// the placeholder and this comment are overwritten on `make generate`, but are included so that the [END] comment survives
