package values

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
	"github.com/ptiger10/pd/options"
)

// InterfaceFactory converts interface{} to Container
func InterfaceFactory(data interface{}) (Container, error) {
	var container Container
	var err error
	if data == nil {
		container = Container{Values: emptyValues(), DataType: options.None}
	} else {
		switch reflect.ValueOf(data).Kind() {
		case reflect.Float32, reflect.Float64,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.String,
			reflect.Bool,
			reflect.Struct:
			container, err = ScalarFactory(data)

		case reflect.Slice:
			container, err = SliceFactory(data)

		default:
			return Container{Values: emptyValues()}, fmt.Errorf("internal.values.InterfaceFactory(): type not supported: %T", data)
		}
	}
	return container, err
}

// InterfaceSliceFactory converts []interface{} to []Container. If manualMode is true, []interface{} columns will not be interpolated.
func InterfaceSliceFactory(data []interface{}, manualMode bool, dataType options.DataType) ([]Container, error) {
	vals := make([]Container, len(data))
	for i := 0; i < len(data); i++ {
		container, err := InterfaceFactory(data[i])
		if err != nil {
			return nil, fmt.Errorf("dataframe.New(): %v", err)
		}
		_, isInterfaceSlice := data[i].([]interface{})
		if !manualMode && isInterfaceSlice {
			vals, _ := data[i].([]interface{})
			interpolateAs := Interpolate(vals)
			if interpolateAs != options.Interface {
				// ducks error because interpolation is controlled
				container.Values, _ = Convert(container.Values, interpolateAs)
			}
			container.DataType = interpolateAs
		}

		// optional DataType conversion
		if dataType != options.None {
			container.Values, err = Convert(container.Values, dataType)
			if err != nil {
				return nil, fmt.Errorf("dataframe.New(): %v", err)
			}
			container.DataType = dataType
		}
		vals[i] = container
	}
	return vals, nil
}

// MustCreateValuesFromInterface returns a container that satisfies the Values interface or panics.
func MustCreateValuesFromInterface(data interface{}) Container {
	container, err := InterfaceFactory(data)
	if err != nil {
		if options.GetLogWarnings() {
			log.Printf("MustCreateValuesFromInterface(): %v", err)
		}
	}
	return container
}

// ScalarFactory creates a Container from an interface{} that has been determined elsewhere to be a scalar.
// Be careful modifying this constructor and its children,
// as reflection is inherently unsafe and each function expects to receive specific types only.
func ScalarFactory(data interface{}) (Container, error) {
	var ret Container
	switch data.(type) {
	case float32, float64:
		val := scalarFloatToFloat64(data)
		f := newFloat64(val)
		ret = Container{&float64Values{f}, options.Float64}

	case int, int8, int16, int32, int64:
		val := scalarIntToInt64(data)
		i := newInt64(val)
		ret = Container{&int64Values{i}, options.Int64}

	case uint, uint8, uint16, uint32, uint64:
		val := scalarUIntToInt64(data)
		i := newInt64(val)
		ret = Container{&int64Values{i}, options.Int64}

	case string:
		val := newString(data.(string))
		ret = Container{&stringValues{val}, options.String}

	case bool:
		val := newBool(data.(bool))
		ret = Container{&boolValues{val}, options.Bool}

	case time.Time:
		val := newDateTime(data.(time.Time))
		ret = Container{&dateTimeValues{val}, options.DateTime}

	default:
		ret = Container{}
		return ret, fmt.Errorf("Type %T not supported", data)
	}

	return ret, nil
}

// scalarFloatToFloat64 converts a known float interface{} -> float64
func scalarFloatToFloat64(data interface{}) float64 {
	d := reflect.ValueOf(data)
	return d.Float()
}

// scalarIntToInt64 converts a known int interface{} -> int64
func scalarIntToInt64(data interface{}) int64 {
	d := reflect.ValueOf(data)
	return d.Int()
}

// scalarUIntToInt64 converts a known uint interface{} -> int64
func scalarUIntToInt64(data interface{}) int64 {
	d := reflect.ValueOf(data)
	return int64(d.Uint())
}

// SliceFactory creates a Container from an interface{} that has been determined elsewhere to be a Slice.
// Be careful modifying this constructor and its children,
// as reflection is inherently unsafe and each function expects to receive specific types only.
func SliceFactory(data interface{}) (Container, error) {
	var ret Container

	switch data.(type) {
	case []float32:
		vals := sliceFloatToSliceFloat64(data)
		ret = newSliceFloat64(vals)

	case []float64:
		vals := data.([]float64)
		ret = newSliceFloat64(vals)

	case []int, []int8, []int16, []int32:
		vals := sliceIntToSliceInt64(data)
		ret = newSliceInt64(vals)
	case []int64:
		vals := data.([]int64)
		ret = newSliceInt64(vals)

	case []uint, []uint8, []uint16, []uint32, []uint64:
		// converts into []int64 so it can be passed into SliceInt
		vals := sliceUIntToSliceInt64(data)
		ret = newSliceInt64(vals)

	case []string:
		vals := data.([]string)
		ret = newSliceString(vals)

	case []bool:
		vals := data.([]bool)
		ret = newSliceBool(vals)

	case []time.Time:
		vals := data.([]time.Time)
		ret = newSliceDateTime(vals)

	case []interface{}:
		vals := data.([]interface{})
		ret = newSliceInterface(vals)

	default:
		ret = Container{nil, options.None}
		return ret, fmt.Errorf("Type %T not supported", data)
	}

	return ret, nil
}

// MapSplitter splits map[string]interface{} into []interface (for use in dataframe constructor) and []string (for use in Column config)
func MapSplitter(data []interface{}) (isSplit bool, extractedData []interface{}, extractedColumns []string) {
	if len(data) == 0 {
		return
	}
	if reflect.TypeOf(data[0]).String() == "map[string]interface {}" {
		for k, v := range data[0].(map[string]interface{}) {
			extractedColumns = append(extractedColumns, k)
			extractedData = append(extractedData, v)
		}
		isSplit = true
	}
	return
}

// [START interface converters]

// sliceFloatToSliceFloat64 converts known []float interface{} -> []float64
func sliceFloatToSliceFloat64(data interface{}) []float64 {
	var vals []float64
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Float()
		vals = append(vals, v)
	}
	return vals
}

// sliceIntToSliceInt64 converts known []int interface{} -> []int64
func sliceIntToSliceInt64(data interface{}) []int64 {
	var vals []int64
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Int()
		vals = append(vals, v)
	}
	return vals
}

// sliceUIntToSliceInt64 converts knonw []uint interface{} -> []int64
func sliceUIntToSliceInt64(data interface{}) []int64 {
	var vals []int64
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Uint()
		vals = append(vals, int64(v))
	}
	return vals
}

// Interpolate counts the number of instances of each dataType option within data, which must be []interface.
// If any ratio exceeeds the Interpoliation Threshold ratio, returns the dataType with the highest ratio.
func Interpolate(data interface{}) options.DataType {
	count := make(map[options.DataType]float64)
	vals := data.([]interface{})
	if n := GetInterpolationMaximum(); len(vals) > n {
		vals = vals[:n]
	}
	for _, val := range vals {
		switch val.(type) {
		case float32, float64:
			count[options.Float64]++
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			count[options.Int64]++
		case string:
			count[options.String]++
		case bool:
			count[options.Bool]++
		case time.Time:
			count[options.DateTime]++
		}
	}
	type ratios struct {
		Float64  float64
		Int64    float64
		String   float64
		Bool     float64
		DateTime float64
	}
	r := ratios{
		Float64:  count[options.Float64] / float64(len(vals)),
		Int64:    count[options.Int64] / float64(len(vals)),
		String:   count[options.String] / float64(len(vals)),
		Bool:     count[options.Bool] / float64(len(vals)),
		DateTime: count[options.DateTime] / float64(len(vals)),
	}
	var max float64
	var maxType string
	v := reflect.ValueOf(r)
	for i := 0; i < v.NumField(); i++ {
		if ratio := v.Field(i).Float(); ratio > max {
			max = ratio
			maxType = v.Type().Field(i).Name
		}
	}
	if max >= GetInterpolationThreshold() {
		return options.DT(maxType)
	}
	if r.Float64+r.Int64 >= GetInterpolationThreshold() {
		return options.Float64
	}
	return options.Interface
}

// InterpolateString converts a string into another datatype if possible, or retains as string otherwise,
// then creates a Container from the new datatype. Primary use is translating column data into index data or reading from [][]string.
func InterpolateString(s string) interface{} {
	if intVal, err := strconv.Atoi(s); err == nil {
		return intVal
	} else if floatVal, err := strconv.ParseFloat(s, 64); err == nil {
		return floatVal
	} else if boolVal, err := strconv.ParseBool(s); err == nil {
		return boolVal
	} else if dateTimeVal, err := dateparse.ParseAny(s); err == nil {
		return dateTimeVal
	}
	return s
}

// [END interface converters]

// [START utility slices]

// MakeIntRange returns a sequential series of numbers, for use in making default index labels.
// Includes min and excludes max.
func MakeIntRange(min, max int) []int {
	a := make([]int, max-min)
	for i := range a {
		a[i] = min + i
	}
	return a
}

// MakeNullRange returns a sequential series of null values, for use in stacking and unstacking columns.
// Includes min and excludes max.
func MakeNullRange(n int) []interface{} {
	a := make([]interface{}, n)
	for i := range a {
		a[i] = ""
	}
	return a
}

// MakeIntRangeInclusive returns a sequential series of numbers, for use in making default index labels.
// Includes start and end.
func MakeIntRangeInclusive(start, end int) []int {
	var ret []int
	if start <= end {
		ret = make([]int, (end-start)+1)
		for i := range ret {
			ret[i] = start + i
		}
	} else {
		ret = make([]int, (start-end)+1)
		for i := range ret {
			ret[i] = start - i
		}
	}

	return ret
}

// MakeStringRange returns a sequential series of numbers as string values, for use in making default column labels.
func MakeStringRange(min, max int) []string {
	a := make([]string, max-min)
	for i := range a {
		a[i] = strconv.Itoa(min + i)
	}
	return a
}

// [END utility slices]

// Copy copies a Values Container
func (vc Container) Copy() Container {
	return Container{
		Values:   vc.Values.Copy(),
		DataType: vc.DataType,
	}
}

// TransposeValues pivots [][]interface{}{row1{col1, col2, col3}} to []interface{}{col1{row1}, col2{row1}, col3{row1}}
func TransposeValues(data [][]interface{}) []interface{} {
	var transposedData [][]interface{}
	if len(data) > 0 {
		transposedData = make([][]interface{}, len(data[0]))
		for m := 0; m < len(data[0]); m++ {
			transposedData[m] = make([]interface{}, len(data))
		}
		for i := 0; i < len(data); i++ {
			for m := 0; m < len(data[0]); m++ {
				transposedData[m][i] = data[i][m]
			}
		}
	}
	var ret []interface{}
	for _, col := range transposedData {
		ret = append(ret, col)
	}
	return ret
}
