package series

import (
	"fmt"
	"log"
	"reflect"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/options"
)

// Join converts s2 to the same type as the base Series (s), appends s2 to the end, and modifies s in place.
func (ip InPlace) Join(s2 *Series) error {
	if ip.s == nil || ip.s.datatype == options.None {
		ip.s.replace(s2)
		return nil
	}

	if s2.index.NumLevels() != ip.s.NumLevels() {
		return fmt.Errorf("Series.Join(): s2 must have same number of index levels as s (%d != %d)", s2.index.NumLevels(), ip.s.NumLevels())
	}
	for i := 0; i < s2.Len(); i++ {
		elem := s2.Element(i)
		ip.s.InPlace.Append(elem.Value, elem.Labels...)
	}
	return nil
}

// Join converts s2 to the same type as the base Series (s), appends s2 to the end, and returns a new Series.
func (s *Series) Join(s2 *Series) (*Series, error) {
	s = s.Copy()
	err := s.InPlace.Join(s2)
	return s, err
}

// LookupSeries performs a vlookup of each values in one Series against another Series.
func (s *Series) LookupSeries(s2 *Series) *Series {
	if s2.index.NumLevels() != s.NumLevels() {
		if options.GetLogWarnings() {
			log.Printf("Series.LookupSeries(): s2 must have same number of index levels as s (%d != %d)\n", s2.index.NumLevels(), s.NumLevels())
		}
		return newEmptySeries()
	}

	matchShallow := func(s *Series, idx index.Elements) int {
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(idx, s.index.Elements(i)) {
				return i
			}
		}
		return -1
	}

	vals := make([]interface{}, 0)
	for i := 0; i < s.Len(); i++ {
		elems := s.index.Elements(i)
		pos := matchShallow(s2, elems)
		if pos != -1 {
			vals = append(vals, s2.At(pos))
		} else {
			vals = append(vals, "")
		}
	}
	// ducks error because there will be no unsupported values coming from an existing series
	ret, _ := New(vals, Config{DataType: s2.datatype})
	ret.index = s.index

	return ret
}
