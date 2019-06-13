package series

import (
	"fmt"
	"sort"

	"github.com/ptiger10/pd/options"
)

// [START Sort interface]
func (s *Series) Swap(i, j int) {
	s.values.Swap(i, j)
	for lvl := 0; lvl < s.index.Len(); lvl++ {
		s.index.Levels[lvl].Labels.Swap(i, j)
		s.index.Levels[lvl].Refresh()
	}
}

func (s *Series) Less(i, j int) bool {
	return s.values.Less(i, j)
}

// [END Sort interface]

// [START return new Series]

// Sort sorts the series by its values and returns a new Series.
func (s *Series) Sort(asc bool) *Series {
	s = s.copy()
	s.InPlace.Sort(asc)
	return s
}

// Insert inserts a new row into the Series immediately before the specified integer position and returns a new Series.
func (s *Series) Insert(pos int, val interface{}, idx []interface{}) (*Series, error) {
	s = s.copy()
	s.InPlace.Insert(pos, val, idx)
	return s, nil
}

// dropRows drops multiple rows and returns a new Series
func (s *Series) dropRows(positions []int) (*Series, error) {
	s = s.copy()
	s.InPlace.dropRows(positions)
	return s, nil
}

// Drop drops the row at the specified integer position and returns a new Series.
func (s *Series) Drop(pos int) (*Series, error) {
	s = s.copy()
	s.InPlace.Drop(pos)
	return s, nil
}

// Append adds a row at the end and returns a new Series.
func (s *Series) Append(val interface{}, idx []interface{}) *Series {
	s, _ = s.Insert(s.Len(), val, idx)
	return s
}

// Join converts s2 to the same type as the base Series (s), appends s2 to the end, and returns a new Series.
func (s *Series) Join(s2 *Series) *Series {
	s = s.copy()
	s.InPlace.Join(s2)
	return s
}

// [END return new Series]

// [START modify in place]

// InPlace contains methods for modifying a Series in place.
type InPlace struct {
	s *Series
}

// Sort sorts the series by its values and modifies the Series in place.
func (ip InPlace) Sort(asc bool) {
	if asc {
		sort.Stable(ip.s)
	} else {
		sort.Stable(sort.Reverse(ip.s))
	}
}

// Insert inserts a new row into the Series immediately before the specified integer position and modifies the Series in place.
func (ip InPlace) Insert(pos int, val interface{}, idx []interface{}) error {
	if len(idx) != ip.s.index.Len() {
		return fmt.Errorf("Series.Insert() len(idx) must equal number of index levels: supplied %v want %v",
			len(idx), ip.s.index.Len())
	}
	for i := 0; i < ip.s.index.Len(); i++ {
		err := ip.s.index.Levels[i].Labels.Insert(pos, idx[i])
		if err != nil {
			return fmt.Errorf("Series.Insert() with idx val %v at idx level %v: %v", val, i, err)
		}
		ip.s.index.Levels[i].Refresh()
	}
	if err := ip.s.values.Insert(pos, val); err != nil {
		return fmt.Errorf("Series.Insert() with val %v: %v", val, err)
	}
	return nil
}

// dropRows drops multiple rows
func (ip InPlace) dropRows(positions []int) error {
	sort.IntSlice(positions).Sort()
	for i, position := range positions {
		err := ip.s.InPlace.Drop(position - i)
		if err != nil {
			return err
		}
	}
	return nil
}

// Drop drops a row at a specified integer position and modifies the Series in place.
func (ip InPlace) Drop(pos int) error {
	for i := 0; i < ip.s.index.Len(); i++ {
		err := ip.s.index.Levels[i].Labels.Drop(pos)
		if err != nil {
			return fmt.Errorf("Series.Drop(): %v", err)
		}
		ip.s.index.Levels[i].Refresh()
	}
	if err := ip.s.values.Drop(pos); err != nil {
		return fmt.Errorf("Series.Drop(): %v", err)
	}
	return nil
}

// Append adds a row at a specified integer position and modifies the Series in place.
func (ip InPlace) Append(val interface{}, idx []interface{}) {
	_ = ip.s.InPlace.Insert(ip.s.Len(), val, idx)
	return
}

// Join converts s2 to the same type as the base Series (s), appends s2 to the end, and modifies s in place.
func (ip InPlace) Join(s2 *Series) {
	if ip.s.values == nil {
		ip.s.replace(s2)
		return
	}
	for i := 0; i < s2.Len(); i++ {
		elem := s2.Element(i)
		ip.s.InPlace.Append(elem.Value, elem.Labels)
	}
}

// [END modify in place]

// [START type conversion methods]

// To contains conversion methods
type To struct {
	s   *Series
	idx bool
}

// Float converts Series values to float64.
func (t To) Float() *Series {
	if t.idx {
		t.s.index = t.s.index.Copy()
		t.s.index.Levels[0] = t.s.index.Levels[0].ToFloat()
	} else {
		t.s.values = t.s.values.ToFloat()
		t.s.datatype = options.Float64
	}
	return t.s
}

// Int converts Series values to int64.
func (t To) Int() *Series {
	if t.idx {
		t.s.index = t.s.index.Copy()
		t.s.index.Levels[0] = t.s.index.Levels[0].ToInt()
	} else {
		t.s.values = t.s.values.ToInt()
		t.s.datatype = options.Int64
	}
	return t.s
}

// String converts Series values to string.
func (t To) String() *Series {
	if t.idx {
		t.s.index = t.s.index.Copy()
		t.s.index.Levels[0] = t.s.index.Levels[0].ToString()
	} else {
		t.s.values = t.s.values.ToString()
		t.s.datatype = options.String
	}
	return t.s
}

// Bool converts Series values to bool.
func (t To) Bool() *Series {
	if t.idx {
		t.s.index = t.s.index.Copy()
		t.s.index.Levels[0] = t.s.index.Levels[0].ToBool()
	} else {
		t.s.values = t.s.values.ToBool()
		t.s.datatype = options.Bool
	}
	return t.s
}

// DateTime converts Series values to time.Time.
func (t To) DateTime() *Series {
	if t.idx {
		t.s.index = t.s.index.Copy()
		t.s.index.Levels[0] = t.s.index.Levels[0].ToDateTime()
	} else {
		t.s.values = t.s.values.ToDateTime()
		t.s.datatype = options.DateTime
	}
	return t.s
}

// Interface converts Series values to interface.
func (t To) Interface() *Series {
	if t.idx {
		t.s.index = t.s.index.Copy()
		t.s.index.Levels[0] = t.s.index.Levels[0].ToInterface()
	} else {
		t.s.values = t.s.values.ToInterface()
		t.s.datatype = options.Interface
	}
	return t.s
}

// [END type conversion methods]

// [START Index modifications]

// Index contains index selection and conversion
type Index struct {
	s  *Series
	To To
}

// Levels returns the number of levels in the index
func (idx Index) Levels() int {
	return idx.s.index.Len()
}

// Len returns the number of items in each level of the index.
func (idx Index) Len() int {
	if len(idx.s.index.Levels) == 0 {
		return 0
	}
	return idx.s.index.Levels[0].Len()
}

// Swap swaps two labels at index level 0 and modifies the index in place.
func (idx Index) Swap(i, j int) {
	idx.s.Swap(i, j)
}

// Less compares two elements and returns true if the first is less than the second.
func (idx Index) Less(i, j int) bool {
	return idx.s.index.Levels[0].Labels.Less(i, j)
}

// Sort sorts the index by index level 0 and modifies the Series in place.
func (idx Index) Sort(asc bool) {
	if asc {
		sort.Stable(idx)
	} else {
		sort.Stable(sort.Reverse(idx))
	}
}

func (s *Series) rename(name string) {
	s = s.copy()
	s.index.Levels[0].Name = name
}
