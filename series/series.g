package series

import (
	"fmt"
	"os"
	"reflect"
)

// A Series is a 1-D container of data
type Series struct {
	Values []interface{}
	Index  Index
	Type   reflect.Kind
}

func (s *Series) String() string {
	var list string
	if s.Index.hasStringIdx() {
		for _, pos := range s.Index.StringIdx {
			list += fmt.Sprintf("%-*v    %v\n", s.Index.longestIdx(), pos, s.Values[s.Index.StringValMap[pos]])
		}
	} else {
		for _, pos := range s.Index.IntIdx {
			list += fmt.Sprintf("%-*v    %v\n", s.Index.longestIdx(), pos, s.Values[s.Index.IntValMap[pos]])
		}
	}
	return fmt.Sprint(list, "dtype: ", s.Type)
}

// Sum of a Series' values
func (s *Series) Sum() int {
	var sum int
	for _, val := range s.Values {
		sum += val.(int)
	}
	return sum
}

// At returns the value at Index position i
// and accepts either string or int values
func (s *Series) At(i interface{}) interface{} {
	switch i.(type) {
	case string:
		return s.Values[s.Index.StringValMap[i.(string)]]
	case int:
		return s.Values[s.Index.IntValMap[i.(int)]]
	default:
		fmt.Fprintf(os.Stderr, "At accepts string or int value; you provided %T", i)
		return nil
	}
}

func isHomogenous(vals []interface{}) error {
	firstVal := vals[0]
	anchorType := reflect.TypeOf(firstVal)
	for i, val := range vals {
		thisType := reflect.TypeOf(val)
		if thisType != anchorType {
			return fmt.Errorf("Heterogeneous values in slice: %v (idx: 0, val: %v) and %v (idx: %d, val: %v)",
				anchorType, firstVal, thisType, i, val)
		}
	}
	return nil
}

func interfaceToIntSlice(vals []interface{}) []int {
	ret := make([]int, len(vals))
	for i, val := range vals {
		ret[i] = val.(int)
	}
	return ret
}

func interfaceToStringSlice(vals []interface{}) []string {
	ret := make([]string, len(vals))
	for i, val := range vals {
		ret[i] = val.(string)
	}
	return ret
}

// New series constructor with an optional index
func New(data interface{}, index ...interface{}) *Series {
	var idx []interface{}
	intIdx := make([]int, len(data))
	stringIdx := make([]string, 0)

	var priorIntIdx bool

	switch len(index) {
	case 0:
		idx = make([]interface{}, len(data))
	case 1:
		idx = index[0]
		if err := isHomogenous(idx); err != nil {
			fmt.Fprintf(os.Stderr, "Index must be of homogenous type: %v", err)
			return nil
		}
		switch idx[0].(type) {
		case string:
			stringIdx = interfaceToStringSlice(idx)
		case int:
			intIdx = interfaceToIntSlice(idx)
			priorIntIdx = true
		}

	default:
		fmt.Fprintf(os.Stderr, "Index parameter accepts at most 1 value. You provided %d: %v", len(index), index)
		return nil
	}

	vals := make([]interface{}, len(data))

	intValMap := make(map[int]int)
	stringValMap := make(map[string]int)
	for i, val := range data {
		if !priorIntIdx {
			intIdx[i] = i
		}
		intValMap[intIdx[i]] = i
		if len(stringIdx) > 0 {
			stringValMap[stringIdx[i]] = i
		}
		vals[i] = val
	}
	return &Series{
		Index:  Index{IntIdx: intIdx, IntValMap: intValMap, StringIdx: stringIdx, StringValMap: stringValMap},
		Values: vals,
		Type:   reflect.Int,
	}
}
