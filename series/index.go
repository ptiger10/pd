package series

import (
	"fmt"
	"reflect"
)

// An Index contains integer and string index positions corresponding to Values.
// After a sort, only the idx-value maps should be updated
// type Index struct {
// 	Loc  map[string]interface{}
// 	ILoc map[int]interface{}

// 	IntIdx       []int
// 	IntValMap    map[int]int
// 	StringIdx    []string
// 	StringValMap map[string]int
// }

// func (i Index) String() string {
// 	if i.hasStringIdx() {
// 		return fmt.Sprintf("%v dtype: %v", i.StringIdx, reflect.String)
// 	}
// 	return fmt.Sprintf("%v dtype: %v", i.IntIdx, reflect.Int)
// }

// func (i Index) hasStringIdx() bool {
// 	if len(i.StringIdx) > 0 {
// 		return true
// 	}
// 	return false
// }

// func (i Index) idxMax() int {
// 	var max int
// 	for _, elem := range i.IntIdx {
// 		if elem > max {
// 			max = elem
// 		}
// 	}
// 	return max
// }

// // count the length of the longest string or the number of digits in the largest int
// func (i Index) longestIdx() int {
// 	var max int
// 	switch i.hasStringIdx() {
// 	case true:
// 		for _, elem := range i.StringIdx {
// 			if len(elem) > max {
// 				max = len(elem)
// 			}
// 		}
// 	case false:
// 		longest := i.idxMax()
// 		for longest != 0 {
// 			longest /= 10
// 			max = max + 1
// 		}
// 	}
// 	return max
// }

type Index struct {
	Levels []IndexLevel
}

type IndexLevel struct {
	Type   reflect.Kind
	Labels Labels
}

// Supported: []int, []string, []datetime

type Labels interface {
	At(int) (interface{}, error)
	From(int, int) (interface{}, error)
	Loc(string) (interface{}, error)
	ILoc(int) (interface{}, error)
}

type intLabels struct {
	l    []int64
	iloc map[int]int
}

func (labels intLabels) At(i int) (interface{}, error) {
	if i >= len(labels.l) {
		return nil, fmt.Errorf("Unable to get index value at %d: out of range", i)
	}
	return labels.l[i], nil
}

func (labels intLabels) From(i, j int) (interface{}, error) {
	if i >= len(labels.l) {
		return nil, fmt.Errorf("Unable to get index value at %d: out of range", i)
	} else if j >= len(labels.l) {
		return nil, fmt.Errorf("Unable to get index value at %d: out of range", j)
	}
	return labels.l[i : j+1], nil
}

func (labels intLabels) Loc(s string) (interface{}, error) {
	return nil, fmt.Errorf("Loc is supported for String Index only")
}

func (labels intLabels) ILoc(i int) (interface{}, error) {
	return nil, fmt.Errorf("ILoc not yet supported")
}
