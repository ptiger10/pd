package index

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
)

// Columns is a collection of column levels, plus name mappings.
type Columns struct {
	NameMap LabelMap
	Levels  []ColLevel
}

// NewColumns returns a new Columns collection from a slice of column levels.
func NewColumns(levels ...ColLevel) Columns {
	if levels == nil {
		return Columns{Levels: make([]ColLevel, 0), NameMap: make(LabelMap)}
	}
	cols := Columns{
		Levels: levels,
	}
	cols.updateNameMap()
	return cols
}

// NewColumnsFromConfig returns new Columns with default length n using a config struct.
func NewColumnsFromConfig(config Config, n int) (Columns, error) {
	var columns Columns

	// both nil: return default index
	if config.Col == nil && config.MultiCol == nil {
		cols := NewDefaultColLevel(n, config.ColName)
		return NewColumns(cols), nil

	}
	// both not nil: return error
	if config.Col != nil && config.MultiCol != nil {
		return Columns{}, fmt.Errorf("columnFactory(): supplying both config.Col and config.MultiCol is ambiguous; supply one or the other")
	}
	// single-level Columns
	if config.Col != nil {
		newLevel := NewColLevel(config.Col, config.ColName)
		columns = NewColumns(newLevel)
	}

	// multi-level Columns
	if config.MultiCol != nil {
		if config.MultiColNames != nil && len(config.MultiColNames) != len(config.MultiCol) {
			return Columns{}, fmt.Errorf(
				"columnFactory(): if MultiColNames is not nil, it must must have same length as MultiCol: %d != %d",
				len(config.MultiColNames), len(config.MultiCol))
		}
		columns = CreateMultiCol(config.MultiCol, config.MultiColNames)
	}
	return columns, nil
}

// CreateMultiCol returns a MultiCol from [][]string
func CreateMultiCol(cols [][]string, colNames []string) Columns {
	var newLevels []ColLevel
	for i := 0; i < len(cols); i++ {
		var levelName string
		if i < len(colNames) {
			levelName = colNames[i]
		}
		newLevel := NewColLevel(cols[i], levelName)
		newLevels = append(newLevels, newLevel)
	}
	return NewColumns(newLevels...)
}

// NewDefaultColumns returns a new Columns collection with default range labels (0, 1, 2, ... n).
func NewDefaultColumns(n int) Columns {
	return NewColumns(NewDefaultColLevel(n, ""))
}

// Len returns the number of labels in every level of the column.
func (cols Columns) Len() int {
	if cols.NumLevels() == 0 {
		return 0
	}
	return cols.Levels[0].Len()
}

// NumLevels returns the number of column levels.
func (cols Columns) NumLevels() int {
	return len(cols.Levels)
}

// MultiNames returns a slice of the names of the column at every level for every column position.
func (cols Columns) MultiNames() [][]string {
	names := make([][]string, cols.Len())
	for k := 0; k < cols.Len(); k++ {
		nameSlice := make([]string, cols.NumLevels())
		for j := 0; j < cols.NumLevels(); j++ {
			nameSlice[j] = cols.Levels[j].Labels[k]
		}
		names[k] = nameSlice
	}
	return names
}

// Names returns the name of every column by concatenating the labels across every level.
func (cols Columns) Names() []string {
	names := make([]string, cols.Len())
	for i, mn := range cols.MultiNames() {
		name := strings.Join(mn, values.GetMultiColNameSeparator())
		names[i] = name
	}
	return names
}

// Name returns the name of the column at col
func (cols Columns) Name(col int) string {
	if cols.NumLevels() == 0 {
		return ""
	}
	return cols.Names()[col]
}

// MultiName returns the names of the column levels at col
func (cols Columns) MultiName(col int) []string {
	if cols.NumLevels() == 0 {
		return []string{}
	}
	return cols.MultiNames()[col]
}

// MaxNameWidth returns the number of characters in the column name with the most characters.
func (cols Columns) MaxNameWidth() int {
	var max int
	for k := range cols.NameMap {
		if length := len(fmt.Sprint(k)); length > max {
			max = length
		}
	}
	return max
}

// UpdateNameMap updates the holistic index map of {index level names: [index level positions]}
func (cols *Columns) updateNameMap() {
	nameMap := make(LabelMap)
	for i, lvl := range cols.Levels {
		nameMap[lvl.Name] = append(nameMap[lvl.Name], i)
	}
	cols.NameMap = nameMap
}

// Refresh updates the global name map and the label mappings at every level.
// Should be called after Series selection or index modification
func (cols *Columns) Refresh() {
	cols.updateNameMap()
	for i := 0; i < len(cols.Levels); i++ {
		cols.Levels[i].Refresh()
	}
}

// A ColLevel is a single collection of column labels within a Columns collection, plus label mappings and metadata.
// It is identical to an index Level except for the Labels, which are a simple []interface{} that do not satisfy the values.Values interface.
type ColLevel struct {
	Name      string
	Labels    []string
	LabelMap  LabelMap
	DataType  options.DataType
	IsDefault bool
}

// NewDefaultColLevel creates a column level with range labels (0, 1, 2, ...n) and optional name.
func NewDefaultColLevel(n int, name string) ColLevel {
	colsInt := values.MakeStringRange(0, n)
	lvl := ColLevel{Labels: colsInt, DataType: options.Int64, Name: name, IsDefault: true}
	lvl.Refresh()
	return lvl
}

// NewColLevel returns a Columns level with updated label map.
func NewColLevel(labels []string, name string) ColLevel {
	if len(labels) == 0 {
		return ColLevel{}
	}
	lvl := ColLevel{
		Labels:   labels,
		Name:     name,
		DataType: options.String,
	}
	lvl.Refresh()
	return lvl
}

// Len returns the number of labels in the column level.
func (lvl ColLevel) Len() int {
	return len(lvl.Labels)
}

// Refresh updates all the label mappings value within a column level.
func (lvl *ColLevel) Refresh() {
	lvl.updateLabelMap()
}

// updateLabelMap updates a single level's map of {label values: [label positions]}.
// A level's label map is agnostic of the actual values in those positions.
func (lvl *ColLevel) updateLabelMap() {
	labelMap := make(LabelMap, lvl.Len())
	for i := 0; i < lvl.Len(); i++ {
		key := fmt.Sprint(lvl.Labels[i])
		labelMap[key] = append(labelMap[key], i)
	}
	lvl.LabelMap = labelMap
}

// ResetDefault converts a column level in place to an ordered default range []int{0, 1, 2,...n}. It is analogous to Reindex but for columns.
func (lvl *ColLevel) ResetDefault() {
	*lvl = NewDefaultColLevel(lvl.Len(), "")
	return
}

// Copy copies a Column Level
func (lvl ColLevel) Copy() ColLevel {
	if reflect.DeepEqual(lvl, ColLevel{}) {
		return ColLevel{}
	}
	lvlCopy := ColLevel{}
	lvlCopy = lvl
	lvlCopy.Labels = make([]string, lvl.Len())
	for i := 0; i < lvl.Len(); i++ {
		lvlCopy.Labels[i] = lvl.Labels[i]
	}
	lvlCopy.LabelMap = make(LabelMap)
	for k, v := range lvl.LabelMap {
		lvlCopy.LabelMap[k] = v
	}
	return lvlCopy
}

// Copy returns a deep copy of each column level.
func (cols Columns) Copy() Columns {
	if reflect.DeepEqual(cols, Columns{Levels: make([]ColLevel, 0), NameMap: make(LabelMap)}) {
		return Columns{Levels: make([]ColLevel, 0), NameMap: make(LabelMap)}
	}
	colsCopy := Columns{NameMap: LabelMap{}}
	for k, v := range cols.NameMap {
		colsCopy.NameMap[k] = v
	}
	for i := 0; i < len(cols.Levels); i++ {
		colsCopy.Levels = append(colsCopy.Levels, cols.Levels[i].Copy())
	}
	return colsCopy
}

// Subset subsets a Columns with all the column levels located at the specified integer positions and modifies the Columns in place.
func (cols *Columns) Subset(colPositions []int) {
	for j := 0; j < cols.NumLevels(); j++ {
		cols.Levels[j].Subset(colPositions)
	}
	cols.updateNameMap()
	return
}

// Subset subsets the label values in a column level at specified integer positions and modifies the ColLevel in place.
func (lvl *ColLevel) Subset(positions []int) {
	var labels []string
	for _, pos := range positions {
		labels = append(labels, lvl.Labels[pos])
	}
	lvl.Labels = labels

	lvl.Refresh()
	return
}
