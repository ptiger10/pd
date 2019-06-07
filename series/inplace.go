package series

import (
	"fmt"
	"sort"
)

// InPlace contains methods for modifying a Series in-place.
type InPlace struct {
	s *Series
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
