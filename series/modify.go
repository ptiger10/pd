package series

import (
	"fmt"
	"sort"
)

// [START Sort interface]
func (s Series) Swap(i, j int) {
	s.values.Swap(i, j)
	for lvl := 0; lvl < s.index.Len(); lvl++ {
		s.index.Levels[lvl].Labels.Swap(i, j)
		s.index.Levels[lvl].Refresh()
	}
}

func (s Series) Less(i, j int) bool {
	return s.values.Less(i, j)
}

// [END Sort interface]

// [START Return New Series]

// Sort sorts the series by its values and returns a new Series.
func (s Series) Sort(asc bool) Series {
	s = s.copy()
	s.InPlace.Sort(asc)
	return s
}

// Insert inserts a new row into the Series immediately before the specified integer position and returns a new Series.
func (s Series) Insert(pos int, val interface{}, idx []interface{}) (Series, error) {
	s = s.copy()
	s.InPlace.Insert(pos, val, idx)
	return s, nil
}

// dropRows drops multiple rows and returns a new Series
func (s Series) dropRows(positions []int) (Series, error) {
	s = s.copy()
	s.InPlace.dropRows(positions)
	return s, nil
}

// Drop drops the row at the specified integer position and returns a new Series.
func (s Series) Drop(pos int) (Series, error) {
	s = s.copy()
	s.InPlace.Drop(pos)
	return s, nil
}

// Append adds a row at the end and returns a new Series.
func (s Series) Append(val interface{}, idx []interface{}) Series {
	s, _ = s.Insert(s.Len(), val, idx)
	return s
}

// [END Return New Series]

// [START Modify InPlace]

// InPlace contains methods for modifying a Series in place.
type InPlace struct {
	s *Series
}

// Sort sorts the series by its values and modifies the Series in place.
func (c InPlace) Sort(asc bool) {
	if asc {
		sort.Stable(c.s)
	} else {
		sort.Stable(sort.Reverse(c.s))
	}
}

// Insert inserts a new row into the Series immediately before the specified integer position and modifies the Series in place.
func (c InPlace) Insert(pos int, val interface{}, idx []interface{}) error {
	if len(idx) != c.s.index.Len() {
		return fmt.Errorf("Series.Insert() len(idx) must equal number of index levels: supplied %v want %v",
			len(idx), c.s.index.Len())
	}
	for i := 0; i < c.s.index.Len(); i++ {
		err := c.s.index.Levels[i].Labels.Insert(pos, idx[i])
		if err != nil {
			return fmt.Errorf("Series.Insert() with idx val %v at idx level %v: %v", val, i, err)
		}
		c.s.index.Levels[i].Refresh()
	}
	if err := c.s.values.Insert(pos, val); err != nil {
		return fmt.Errorf("Series.Insert() with val %v: %v", val, err)
	}
	return nil
}

// dropRows drops multiple rows
func (c InPlace) dropRows(positions []int) error {
	sort.IntSlice(positions).Sort()
	for i, position := range positions {
		err := c.s.InPlace.Drop(position - i)
		if err != nil {
			return err
		}
	}
	return nil
}

// Drop drops a row at a specified integer position and modifies the Series in place.
func (c InPlace) Drop(pos int) error {
	for i := 0; i < c.s.index.Len(); i++ {
		err := c.s.index.Levels[i].Labels.Drop(pos)
		if err != nil {
			return fmt.Errorf("Series.Drop(): %v", err)
		}
		c.s.index.Levels[i].Refresh()
	}
	if err := c.s.values.Drop(pos); err != nil {
		return fmt.Errorf("Series.Drop(): %v", err)
	}
	return nil
}

// Append adds a row at a specified integer position and modifies the Series in place.
func (c InPlace) Append(val interface{}, idx []interface{}) {
	_ = c.s.InPlace.Insert(c.s.Len(), val, idx)
	return
}

// [END Modify InPlace]

// [START Index modifications]

// Index contains index selection and conversion
type Index struct {
	s  *Series
	To To
}
