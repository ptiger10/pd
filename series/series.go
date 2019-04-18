package series

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Series struct {
	Index  Index
	Values Values
	Kind   reflect.Kind
}

func (s Series) Sum() (float64, error) {
	switch s.Kind {
	case Float:
		i := s.Values.(floatValues)
		return i.Sum()
	case Int:
		i := s.Values.(intValues)
		return i.Sum()
	case Bool:
		i := s.Values.(boolValues)
		return i.Sum()
	default:
		return math.NaN(), fmt.Errorf("Sum not supported for type %v", s.Kind)
	}
}

func (s Series) Count() int {
	var count int
	switch s.Kind {
	case Float:
		vals := s.Values.(floatValues)
		for _, val := range vals {
			if !val.null {
				count++
			}
		}
	case Int:
		vals := s.Values.(intValues)
		for _, val := range vals {
			if !val.null {
				count++
			}
		}
	case String:
		vals := s.Values.(stringValues)
		for _, val := range vals {
			if !val.null {
				count++
			}
		}
	case Bool:
		vals := s.Values.(boolValues)
		for _, val := range vals {
			if !val.null {
				count++
			}
		}
	default:
		return 0
	}
	return count
}

type Index struct {
	Levels []IndexLevel
}

type IndexLevel struct {
	Type   reflect.Kind
	Values Values
}

type Values interface {
}

type floatValues []floatValue
type floatValue struct {
	v    float64
	null bool
}

func (vals floatValues) Sum() (float64, error) {
	var sum float64
	for _, val := range vals {
		if !val.null {
			sum += val.v
		}
	}
	return sum, nil
}

type intValues []intValue
type intValue struct {
	v    int64
	null bool
}

func (vals intValues) Sum() (float64, error) {
	var sum float64
	for _, val := range vals {
		if !val.null {
			sum += float64(val.v)
		}
	}
	return sum, nil
}

type stringValues []stringValue
type stringValue struct {
	v    string
	null bool
}

func isNullString(s string) bool {
	nullStrings := []string{"nan", "n/a", ""}
	for _, ns := range nullStrings {
		if strings.ToLower(s) == ns {
			return true
		}
	}
	return false
}

type boolValues []boolValue
type boolValue struct {
	v    bool
	null bool
}

func (vals boolValues) Sum() (float64, error) {
	var sum float64
	for _, val := range vals {
		if val.v && !val.null {
			sum++
		}
	}
	return sum, nil
}

type timeValues []timeValue
type timeValue struct {
	v    time.Time
	null bool
}

func (vals floatValues) Count() int {
	var count int
	for _, val := range vals {
		if !val.null {
			count++
		}
	}
	return count
}

func (vals intValues) Describe() string {
	return ""
}

func (vals stringValues) Describe() string {
	return ""
}

func (vals boolValues) Describe() string {
	return ""
}

func (vals timeValues) Describe() string {
	return ""
}

type newSeriesOption func(*newSeriesConfig)
type newSeriesConfig struct {
	kind  reflect.Kind
	index interface{}
}

func SeriesType(t reflect.Kind) newSeriesOption {
	return func(c *newSeriesConfig) {
		c.kind = t
	}
}

const (
	Float  = reflect.Float64
	Int    = reflect.Int64
	String = reflect.String
	Bool   = reflect.Bool
	Time   = reflect.Struct
	None   = reflect.UnsafePointer // pseudo-nil value for type reflect.Kind
)

func New(data interface{}, options ...newSeriesOption) (Series, error) {
	advanced := newSeriesConfig{kind: None}
	for _, option := range options {
		option(&advanced)
	}
	s := Series{
		Kind: advanced.kind,
	}

	switch data.(type) {
	case []float32, []float64:
		vals := floatToFloatValues(data)
		s.Values = vals
		s.Kind = Float

	case []int, []int8, []int16, []int32, []int64:
		vals := intToIntValues(data)
		s.Values = vals
		s.Kind = Int

	case []uint, []uint8, []uint16, []uint32, []uint64:
		vals := uIntToIntValues(data)
		s.Values = vals
		s.Kind = Int

	case []string:
		vals := stringToStringValues(data)
		s.Values = vals
		s.Kind = String

	case []bool:
		vals := boolToBoolValues(data)
		s.Values = vals
		s.Kind = Bool

	case []time.Time:
		vals := timeToTimeValues(data)
		s.Values = vals
		s.Kind = Time

	case []interface{}:
		d := reflect.ValueOf(data)
		switch advanced.kind {
		case None: // this checks for the pseduo-nil type
			return Series{}, fmt.Errorf("Must supply a SeriesType to decode interface")
		case Float:
			vals := interfaceToFloatValues(d)
			s.Values = vals
		case Int:
			vals := interfaceToIntValues(d)
			s.Values = vals
		case String:
			vals := interfaceToStringValues(d)
			s.Values = vals
		case Bool:
			vals := interfaceToBoolValues(d)
			s.Values = vals
		default:
			return s, fmt.Errorf("Type not supported for conversion from []interface: %v", advanced.kind)
		}

	default:
		return s, fmt.Errorf("Type not supported: %T", data)
	}
	return s, nil
}

func floatToFloatValues(data interface{}) floatValues {
	var vals []floatValue
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		val := d.Index(i).Float()
		vals = append(vals, floatValue{v: val})
		if math.IsNaN(val) {
			vals[i].null = true
		}
	}
	return floatValues(vals)
}

func intToIntValues(data interface{}) intValues {
	var vals []intValue
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		vals = append(vals, intValue{v: d.Index(i).Int()})
	}
	return intValues(vals)
}

func uIntToIntValues(data interface{}) intValues {
	var vals []intValue
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		val := int64(d.Index(i).Uint())
		vals = append(vals, intValue{v: val})
	}
	return intValues(vals)
}

func stringToStringValues(data interface{}) stringValues {
	var vals []stringValue
	d := data.([]string)
	for i := 0; i < len(d); i++ {
		val := d[i]
		if isNullString(val) {
			vals = append(vals, stringValue{null: true})
		} else {
			vals = append(vals, stringValue{v: val})
		}
	}
	return stringValues(vals)
}

func boolToBoolValues(data interface{}) boolValues {
	var vals []boolValue
	d := data.([]bool)
	for i := 0; i < len(d); i++ {
		vals = append(vals, boolValue{v: d[i]})
	}
	return boolValues(vals)
}

func timeToTimeValues(data interface{}) timeValues {
	var vals []timeValue
	d := data.([]time.Time)
	for i := 0; i < len(d); i++ {
		val := d[i]
		vals = append(vals, timeValue{v: val})
		if (time.Time{}) == val {
			vals[i].null = true
		}
	}
	return timeValues(vals)
}

func interfaceToFloatValues(d reflect.Value) floatValues {
	var vals []floatValue
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Elem()
		switch d.Index(i).Elem().Kind() {
		case reflect.Invalid:
			vals = append(vals, floatValue{null: true})
		case reflect.String:
			vals = append(vals, stringToFloat(v.String()))
		case reflect.Float32, reflect.Float64:
			val := v.Float()
			vals = append(vals, floatValue{v: val})
			if math.IsNaN(val) {
				vals[i].null = true
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vals = append(vals, floatValue{v: float64(v.Int())})
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vals = append(vals, floatValue{v: float64(v.Uint())})
		}
	}
	return floatValues(vals)
}

func stringToFloat(v string) floatValue {
	val, err := strconv.ParseFloat(v, 64)
	if err != nil || math.IsNaN(val) {
		return floatValue{null: true, v: math.NaN()}
	}
	return floatValue{v: val}
}

func interfaceToIntValues(d reflect.Value) intValues {
	var vals []intValue
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Elem()
		switch v.Kind() {
		case reflect.Invalid:
			vals = append(vals, intValue{null: true})
		case reflect.String:
			val, err := strconv.ParseFloat(v.String(), 64)
			if err != nil || math.IsNaN(val) {
				vals = append(vals, intValue{null: true, v: int64(math.NaN())})
			} else {
				vals = append(vals, intValue{v: int64(val)})
			}
		case reflect.Float32, reflect.Float64:
			val := v.Float()
			vals = append(vals, intValue{v: int64(val)})
			if math.IsNaN(val) {
				vals[i].null = true
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vals = append(vals, intValue{v: int64(v.Int())})
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vals = append(vals, intValue{v: int64(v.Uint())})
		}
	}
	return intValues(vals)
}

func interfaceToStringValues(d reflect.Value) stringValues {
	var vals []stringValue
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Elem()
		switch v.Kind() {
		case reflect.Invalid:
			vals = append(vals, stringValue{null: true})
		case reflect.String:
			val := v.String()
			if isNullString(val) {
				vals = append(vals, stringValue{null: true})
			} else {
				vals = append(vals, stringValue{v: val})
			}
		case reflect.Float32, reflect.Float64:
			vals = append(vals, floatToString(v.Float()))

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vals = append(vals, stringValue{v: fmt.Sprint(v.Int())})
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vals = append(vals, stringValue{v: fmt.Sprint(v.Uint())})
		}
	}
	return stringValues(vals)
}

func floatToString(v float64) stringValue {
	if math.IsNaN(v) {
		return stringValue{null: true}
	}
	return stringValue{v: fmt.Sprint(v)}
}

func interfaceToBoolValues(d reflect.Value) boolValues {
	var vals []boolValue
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Elem()
		switch v.Kind() {
		case reflect.Invalid:
			vals = append(vals, boolValue{null: true})
		case reflect.String:
			if isNullString(v.String()) {
				vals = append(vals, boolValue{null: true})
			} else {
				vals = append(vals, boolValue{v: true})
			}
		case reflect.Float32, reflect.Float64:
			val := v.Float()
			if math.IsNaN(val) {
				vals = append(vals, boolValue{null: true})
			} else {
				vals = append(vals, boolValue{v: true})
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vals = append(vals, boolValue{v: true})
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vals = append(vals, boolValue{v: true})
		}
	}
	return boolValues(vals)
}
