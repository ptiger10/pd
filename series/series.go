package series

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
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
	return count(s.Values)
}

func (s Series) String() string {
	switch s.Kind {
	case DateTime:
		var printer string
		vals := s.Values.(dateTimeValues)
		for _, val := range vals {
			printer += fmt.Sprintln(val.v.Format("01/02/2006"))
		}
		return printer
	default:
		return print(s.Values)
	}
}

// expects to receive a slice of typed value structs (eg pd.floatValues)
// each struct must contain a boolean field called "null"
func count(data interface{}) int {
	var count int
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		if !d.Index(i).FieldByName("null").Bool() {
			count++
		}
	}
	return count
}

// expects to receive a slice of typed value structs (eg pd.floatValues)
// each struct must contain a boolean field called "v"
func print(data interface{}) string {
	d := reflect.ValueOf(data)
	var printer string
	for i := 0; i < d.Len(); i++ {
		printer += fmt.Sprintln(d.Index(i).FieldByName("v"))
	}
	return printer
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
		if strings.TrimSpace(strings.ToLower(s)) == ns {
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

type dateTimeValues []dateTimeValue
type dateTimeValue struct {
	v    time.Time
	null bool
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

func (vals dateTimeValues) Describe() string {
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
	Float    = reflect.Float64
	Int      = reflect.Int64
	String   = reflect.String
	Bool     = reflect.Bool
	DateTime = reflect.Struct        // time.Time{} are the only structs accepted by constructor
	None     = reflect.UnsafePointer // pseudo-nil value for type reflect.Kind
)

// New Series constructor
// that expects to receive a slice of values.
// If passing []interface{}, must supply a type expectation for the Series.
// Options: Float, Int, String, Bool, DateTime
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
		vals := timeToDateTimeValues(data)
		s.Values = vals
		s.Kind = DateTime

	case []interface{}:
		d := reflect.ValueOf(data)
		switch advanced.kind {
		case None: // this checks for the pseduo-nil type
			return Series{}, fmt.Errorf("Must supply a SeriesType to decode interface")
		case Float:
			vals, err := interfaceToFloatValues(d)
			if err != nil {
				return Series{}, fmt.Errorf("Unable to convert data to Float-based Series: %v", err)
			}
			s.Values = vals
		case Int:
			vals, err := interfaceToIntValues(d)
			if err != nil {
				return Series{}, fmt.Errorf("Unable to convert data to Int-based Series: %v", err)
			}
			s.Values = vals
		case String:
			vals, err := interfaceToStringValues(d)
			if err != nil {
				return Series{}, fmt.Errorf("Unable to convert data to String-based Series: %v", err)
			}
			s.Values = vals
		case Bool:
			vals, err := interfaceToBoolValues(d)
			if err != nil {
				return Series{}, fmt.Errorf("Unable to convert data to Bool-based Series: %v", err)
			}
			s.Values = vals
		case DateTime:
			vals, err := interfaceToDateTimeValues(d)
			if err != nil {
				return Series{}, fmt.Errorf("Unable to convert data to time.Time-based Series: %v", err)
			}
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

func timeToDateTimeValues(data interface{}) dateTimeValues {
	var vals []dateTimeValue
	d := data.([]time.Time)
	for i := 0; i < len(d); i++ {
		val := d[i]
		vals = append(vals, dateTimeValue{v: val})
		if (time.Time{}) == val {
			vals[i].null = true
		}
	}
	return dateTimeValues(vals)
}

func interfaceToFloatValues(d reflect.Value) (floatValues, error) {
	var vals []floatValue
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Elem()
		switch v.Kind() {
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
		default:
			return nil, fmt.Errorf("Unsupported datatype (%T) passed in interface as val %v", v, v)
		}

	}
	return floatValues(vals), nil
}

func stringToFloat(v string) floatValue {
	val, err := strconv.ParseFloat(v, 64)
	if err != nil || math.IsNaN(val) {
		return floatValue{null: true, v: math.NaN()}
	}
	return floatValue{v: val}
}

func interfaceToIntValues(d reflect.Value) (intValues, error) {
	var vals []intValue
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Elem()
		switch v.Kind() {
		case reflect.Invalid:
			vals = append(vals, intValue{null: true})
		case reflect.String:
			val, err := strconv.ParseFloat(v.String(), 64)
			if err != nil || math.IsNaN(val) {
				vals = append(vals, intValue{null: true})
			} else {
				vals = append(vals, intValue{v: int64(val)})
			}
		case reflect.Float32, reflect.Float64:
			val := v.Float()
			if math.IsNaN(val) {
				vals = append(vals, intValue{null: true})
			} else {
				vals = append(vals, intValue{v: int64(val)})
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vals = append(vals, intValue{v: int64(v.Int())})
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vals = append(vals, intValue{v: int64(v.Uint())})
		default:
			return nil, fmt.Errorf("Unsupported datatype (%T) passed in interface as val %v", v, v)
		}
	}
	return intValues(vals), nil
}

func interfaceToStringValues(d reflect.Value) (stringValues, error) {
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
		default:
			return nil, fmt.Errorf("Unsupported datatype (%T) passed in interface as val %v", v, v)
		}
	}
	return stringValues(vals), nil
}

func floatToString(v float64) stringValue {
	if math.IsNaN(v) {
		return stringValue{null: true}
	}
	return stringValue{v: fmt.Sprint(v)}
}

func interfaceToBoolValues(d reflect.Value) (boolValues, error) {
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
		default:
			return nil, fmt.Errorf("Unsupported datatype (%T) passed in interface as val %v", v, v)
		}
	}
	return boolValues(vals), nil
}

func interfaceToDateTimeValues(d reflect.Value) (dateTimeValues, error) {
	var vals []dateTimeValue
	for i := 0; i < d.Len(); i++ {
		v := d.Index(i).Elem()
		switch v.Kind() {
		case reflect.Invalid:
			vals = append(vals, dateTimeValue{null: true})
		case reflect.String:
			s := v.String()
			if isNullString(s) {
				vals = append(vals, dateTimeValue{null: true})
			} else {
				val, err := dateparse.ParseAny(s)
				if err != nil {
					vals = append(vals, dateTimeValue{null: true})
				} else {
					vals = append(vals, dateTimeValue{v: val.UTC()})
				}
			}
		case reflect.Struct:
			val, ok := v.Interface().(time.Time)
			if !ok {
				return nil, fmt.Errorf("Unsupported datatype (%T) passed in interface as val %v", v, v)
			}
			vals = append(vals, dateTimeValue{v: val})
		case reflect.Float32, reflect.Float64:
			vals = append(vals, dateTimeValue{null: true})

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vals = append(vals, dateTimeValue{null: true})
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vals = append(vals, dateTimeValue{null: true})

		default:
			return nil, fmt.Errorf("Unsupported datatype (%T) passed in interface as val %v", v, v)
		}
	}
	return dateTimeValues(vals), nil
}

// Acceptable DateTime formats
/*
	"May 8, 2009 5:57:51 PM",
	"oct 7, 1970",
	"oct 7, '70",
	"oct. 7, 1970",
	"oct. 7, 70",
	"Mon Jan  2 15:04:05 2006",
	"Mon Jan  2 15:04:05 MST 2006",
	"Mon Jan 02 15:04:05 -0700 2006",
	"Monday, 02-Jan-06 15:04:05 MST",
	"Mon, 02 Jan 2006 15:04:05 MST",
	"Tue, 11 Jul 2017 16:28:13 +0200 (CEST)",
	"Mon, 02 Jan 2006 15:04:05 -0700",
	"Thu, 4 Jan 2018 17:53:36 +0000",
	"Mon Aug 10 15:44:11 UTC+0100 2015",
	"Fri Jul 03 2015 18:04:07 GMT+0100 (GMT Daylight Time)",
	"September 17, 2012 10:09am",
	"September 17, 2012 at 10:09am PST-08",
	"September 17, 2012, 10:10:09",
	"October 7, 1970",
	"October 7th, 1970",
	"12 Feb 2006, 19:17",
	"12 Feb 2006 19:17",
	"7 oct 70",
	"7 oct 1970",
	"03 February 2013",
	"1 July 2013",
	"2013-Feb-03",
	//   mm/dd/yy
	"3/31/2014",
	"03/31/2014",
	"08/21/71",
	"8/1/71",
	"4/8/2014 22:05",
	"04/08/2014 22:05",
	"4/8/14 22:05",
	"04/2/2014 03:00:51",
	"8/8/1965 12:00:00 AM",
	"8/8/1965 01:00:01 PM",
	"8/8/1965 01:00 PM",
	"8/8/1965 1:00 PM",
	"8/8/1965 12:00 AM",
	"4/02/2014 03:00:51",
	"03/19/2012 10:11:59",
	"03/19/2012 10:11:59.3186369",
	// yyyy/mm/dd
	"2014/3/31",
	"2014/03/31",
	"2014/4/8 22:05",
	"2014/04/08 22:05",
	"2014/04/2 03:00:51",
	"2014/4/02 03:00:51",
	"2012/03/19 10:11:59",
	"2012/03/19 10:11:59.3186369",
	// Chinese
	"2014年04月08日",
	//   yyyy-mm-ddThh
	"2006-01-02T15:04:05+0000",
	"2009-08-12T22:15:09-07:00",
	"2009-08-12T22:15:09",
	"2009-08-12T22:15:09Z",
	//   yyyy-mm-dd hh:mm:ss
	"2014-04-26 17:24:37.3186369",
	"2012-08-03 18:31:59.257000000",
	"2014-04-26 17:24:37.123",
	"2013-04-01 22:43",
	"2013-04-01 22:43:22",
	"2014-12-16 06:20:00 UTC",
	"2014-12-16 06:20:00 GMT",
	"2014-04-26 05:24:37 PM",
	"2014-04-26 13:13:43 +0800",
	"2014-04-26 13:13:43 +0800 +08",
	"2014-04-26 13:13:44 +09:00",
	"2012-08-03 18:31:59.257000000 +0000 UTC",
	"2015-09-30 18:48:56.35272715 +0000 UTC",
	"2015-02-18 00:12:00 +0000 GMT",
	"2015-02-18 00:12:00 +0000 UTC",
	"2015-02-08 03:02:00 +0300 MSK m=+0.000000001",
	"2015-02-08 03:02:00.001 +0300 MSK m=+0.000000001",
	"2017-07-19 03:21:51+00:00",
	"2014-04-26",
	"2014-04",
	"2014",
	"2014-05-11 08:20:13,787",
	// mm.dd.yy
	"3.31.2014",
	"03.31.2014",
	"08.21.71",
	"2014.03",
	"2014.03.30",
	//  yyyymmdd and similar
	"20140601",
	"20140722105203",
	// unix seconds, ms, micro, nano
	"1332151919",
	"1384216367189",
	"1384216367111222",
	"1384216367111222333",
}
*/
