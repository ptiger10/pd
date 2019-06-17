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
	if levels == nil {
		emptyLevel, _ := NewLevel(nil, "")
		levels = append(levels, emptyLevel)
	}
	idx := Index{
		Levels: levels,
	}
	idx.updateNameMap()
	return idx
}

// NewDefault creates a new index with one unnamed index level and range labels (0, 1, 2, ...n)
func NewDefault(length int) Index {
	level := NewDefaultLevel(length, "")
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
func (idx Index) Aligned() error {
	lvl0 := idx.Levels[0].Len()
	for i := 1; i < idx.Len(); i++ {
		if cmpLvl := idx.Levels[i].Len(); lvl0 != cmpLvl {
			return fmt.Errorf("index.Aligned(): index level %v must have same number of labels as level 0, %d != %d",
				i, cmpLvl, lvl0)
		}
	}
	return nil
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

// Elements refer to all the elements at the same position across all levels of an index.
type Elements struct {
	Labels    []interface{}
	DataTypes []options.DataType
}

// updateNameMap updates the holistic index map of {index level names: [index level positions]}
func (idx *Index) updateNameMap() {
	nameMap := make(LabelMap)
	for i, lvl := range idx.Levels {
		nameMap[lvl.Name] = append(nameMap[lvl.Name], i)
	}
	idx.NameMap = nameMap
}

// Refresh updates the global name map and the label mappings at every level.
// Should be called after Series selection or index modification
func (idx *Index) Refresh() {
	if idx.Len() == 0 {
		return
	}
	idx.updateNameMap()
	for i := 0; i < len(idx.Levels); i++ {
		idx.Levels[i].Refresh()
	}
}

// NewFromConfig returns a new Index with length n using a config struct.
func NewFromConfig(n int, config Config) (Index, error) {
	var index Index
	// both nil: return default index
	if config.Index == nil && config.MultiIndex == nil {
		return NewDefault(n), nil
	}
	// both not nil: return error
	if config.Index != nil && config.MultiIndex != nil {
		return Index{}, fmt.Errorf("internal index.NewFromConfig(): supplying both config.Index and config.MultiIndex is ambiguous; supply one or the other")
	}
	// single index
	if config.Index != nil {
		newLevel, err := NewLevel(config.Index, config.IndexName)
		if err != nil {
			return Index{}, fmt.Errorf("internal index.NewFromConfig(): %v", err)
		}
		return New(newLevel), nil
	}
	// multi index
	if config.MultiIndex != nil {
		// name misalignment
		if config.MultiIndexNames != nil && len(config.MultiIndexNames) != len(config.MultiIndex) {
			return Index{}, fmt.Errorf(
				"internal index.NewFromConfig(): if MultiIndexNames is not nil, it must must have same length as MultiIndex: %d != %d",
				len(config.MultiIndexNames), len(config.MultiIndex))
		}
		var newLevels []Level
		for i := 0; i < len(config.MultiIndex); i++ {
			var levelName string
			if i < len(config.MultiIndexNames) {
				levelName = config.MultiIndexNames[i]
			}
			newLevel, err := NewLevel(config.MultiIndex[i], levelName)
			if err != nil {
				return Index{}, fmt.Errorf("internal index.NewFromConfig(): %v", err)
			}
			newLevels = append(newLevels, newLevel)
		}
		return New(newLevels...), nil
	}
	return index, nil
}

// NewColumnFromConfig returns new Columns with length n using a config struct.
func NewColumnFromConfig(n int, config Config) (Columns, error) {
	var columns Columns
	// Handling columns
	if config.Cols != nil && config.MultiCols != nil {
		return Columns{}, fmt.Errorf("columnFactory(): supplying both config.Index and config.MultiIndex is ambiguous; supply one or the other")
	}
	if config.Cols != nil {
		newLevel := NewColLevel(config.Cols, config.ColsName)
		columns = NewColumns(newLevel)
	}
	if config.MultiIndex != nil {
		if config.MultiColsNames != nil && len(config.MultiColsNames) != len(config.MultiCols) {
			return Columns{}, fmt.Errorf(
				"columnFactory(): if MultiColsNames is not nil, it must must have same length as MultiCols: %d != %d",
				len(config.MultiColsNames), len(config.MultiCols))
		}
		var newLevels []ColLevel
		for i := 0; i < len(config.MultiCols); i++ {
			var levelName string
			if i < len(config.MultiColsNames) {
				levelName = config.MultiColsNames[i]
			}
			newLevel := NewColLevel(config.MultiCols[i], levelName)
			newLevels = append(newLevels, newLevel)
		}
		columns = NewColumns(newLevels...)
	}
	return columns, nil
}

// A Config customizes the construction of an Index or Columns object.
type Config struct {
	Index           interface{}
	IndexName       string
	MultiIndex      []interface{}
	MultiIndexNames []string
	Cols            []interface{}
	ColsName        string
	MultiCols       [][]interface{}
	MultiColsNames  []string
}
