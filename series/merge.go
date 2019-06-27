package series

import (
	"fmt"
	"reflect"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/options"
)

// Join converts s2 to the same type as the base Series (s), appends s2 to the end, and modifies s in place.
func (ip InPlace) Join(s2 *Series) {
	if ip.s == nil || ip.s.datatype == options.None {
		ip.s.replace(s2)
		return
	}
	for i := 0; i < s2.Len(); i++ {
		elem := s2.Element(i)
		ip.s.InPlace.Append(elem.Value, elem.Labels)
	}
}

// Join converts s2 to the same type as the base Series (s), appends s2 to the end, and returns a new Series.
func (s *Series) Join(s2 *Series) *Series {
	s = s.Copy()
	s.InPlace.Join(s2)
	return s
}

// match returns the row position of the first match of index Elements within a Series, or -1 if no match exists.
func (s *Series) match(idx index.Elements) int {
	for i := 0; i < s.Len(); i++ {
		if reflect.DeepEqual(idx, s.index.Elements(i)) {
			return i
		}
	}
	return -1
}

// LookupSeries performs a vlookup of each values in one Series against another Series.
func (s *Series) LookupSeries(s2 *Series) *Series {
	vals := make([]interface{}, 0)
	idx := make([]interface{}, 0)
	for i := 0; i < s.Len(); i++ {
		elems := s.index.Elements(i)
		pos := s2.match(elems)
		idx = append(idx, elems.Labels)
		if pos != -1 {
			vals = append(vals, s2.At(pos))
		} else {
			vals = append(vals, "")
		}
	}
	ret, err := New(vals, Config{MultiIndex: idx, DataType: s2.datatype})
	if err != nil {
		fmt.Println(err)
	}
	return ret
}
