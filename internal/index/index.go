package index

import (
	"fmt"
	"sort"

	"github.com/ptiger10/pd/options"
)

// An Index is a collection of levels, plus label mappings
type Index struct {
	Levels  []Level
	NameMap LabelMap
}

// New receives one or more Levels and returns a new Index.
// Expects that Levels already have .LabelMap and .Longest set.
func New(levels ...Level) Index {
	idx := Index{
		Levels: levels,
	}
	idx.UpdateNameMap()
	return idx
}

// Default creates an index with one unnamed index level and range labels (0, 1, 2, ...n)
func Default(length int) Index {
	level := DefaultLevel(length, "")
	return New(level)
}

// In returns a copy of the index with only those levels located at specified integer positions
func (idx Index) In(positions []int) (Index, error) {
	idx = idx.Copy()
	var lvls []Level
	for _, pos := range positions {
		if pos >= len(idx.Levels) {
			return Index{}, fmt.Errorf("error indexing index levels: level %d is out of range", pos)
		}
		lvls = append(lvls, idx.Levels[pos])
	}
	newIdx := New(lvls...)
	return newIdx, nil
}

// Copy returns a deep copy of each index level
func (idx Index) Copy() Index {
	idxCopy := Index{NameMap: LabelMap{}}
	for k, v := range idx.NameMap {
		idxCopy.NameMap[k] = v
	}
	for i := 0; i < len(idx.Levels); i++ {
		idxCopy.Levels = append(idxCopy.Levels, idx.Levels[i].Copy())
	}
	return idxCopy
}

// Drop drops an index level and modifies the Index in-place. If there one or fewer levels, does nothing.
func (idx *Index) Drop(level int) error {
	if idx.Len() <= 1 {
		return nil
	}
	if level >= idx.Len() {
		return fmt.Errorf("invalid level: %v (max: %v)", level, idx.Len())
	}
	idx.Levels = append(idx.Levels[:level], idx.Levels[level+1:]...)
	idx.Refresh()
	return nil
}

// dropLevels drops multiple rows
func (idx *Index) dropLevels(positions []int) error {
	sort.IntSlice(positions).Sort()
	for i, position := range positions {
		err := idx.Drop(position - i)
		if err != nil {
			return err
		}
	}
	return nil
}

// Len returns the number of levels in the index.
func (idx Index) Len() int {
	return len(idx.Levels)
}

// Aligned ensures that all index levels have the same length.
func (idx Index) Aligned() bool {
	lvl0 := idx.Levels[0].Len()
	for i := 1; i < idx.Len(); i++ {
		if lvl0 != idx.Levels[i].Len() {
			return false
		}
	}
	return true
}

// DataTypes returns a slice of the DataTypes at each level of the index
func (idx Index) DataTypes() []options.DataType {
	var idxDataTypes []options.DataType
	for _, lvl := range idx.Levels {
		idxDataTypes = append(idxDataTypes, lvl.DataType)
	}
	return idxDataTypes
}

// Elements returns all the index elements at an integer position.
func (idx Index) Elements(position int) Elements {
	var labels []interface{}
	var datatypes []options.DataType
	for _, lvl := range idx.Levels {
		elem := lvl.Element(position)
		labels = append(labels, elem.Label)
		datatypes = append(datatypes, elem.DataType)
	}
	return Elements{labels, datatypes}
}

type Elements struct {
	Labels    []interface{}
	DataTypes []options.DataType
}
