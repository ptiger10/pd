package pd

import (
	"fmt"
	"os"
	"reflect"

	"github.com/ptiger10/pd/types"
)

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

// Series constructor with an optional index
func Series(data []int, index ...[]interface{}) *types.Series {
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
	return &types.Series{
		Index:  types.Index{IntIdx: intIdx, IntValMap: intValMap, StringIdx: stringIdx, StringValMap: stringValMap},
		Values: vals,
		Type:   reflect.Int,
	}
}
