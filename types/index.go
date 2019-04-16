package types

import (
	"fmt"
	"reflect"
)

// An Index contains integer and string index positions corresponding to Values.
// After a sort, only the idx-value maps should be updated
type Index struct {
	IntIdx       []int
	IntValMap    map[int]int
	StringIdx    []string
	StringValMap map[string]int
}

func (i Index) String() string {
	if i.hasStringIdx() {
		return fmt.Sprintf("%v dtype: %v", i.StringIdx, reflect.String)
	}
	return fmt.Sprintf("%v dtype: %v", i.IntIdx, reflect.Int)
}

func (i Index) hasStringIdx() bool {
	if len(i.StringIdx) > 0 {
		return true
	}
	return false
}

func (i Index) idxMax() int {
	var max int
	for _, elem := range i.IntIdx {
		if elem > max {
			max = elem
		}
	}
	return max
}

// count the length of the longest string or the number of digits in the largest int
func (i Index) longestIdx() int {
	var max int
	switch i.hasStringIdx() {
	case true:
		for _, elem := range i.StringIdx {
			if len(elem) > max {
				max = len(elem)
			}
		}
	case false:
		longest := i.idxMax()
		for longest != 0 {
			longest /= 10
			max = max + 1
		}
	}
	return max
}
