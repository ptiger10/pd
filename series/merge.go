package series

import "github.com/ptiger10/pd/options"

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
