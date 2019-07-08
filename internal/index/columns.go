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
	col := Columns{
		Levels: levels,
	}
	col.updateNameMap()
	return col
}

// NewColumnsFromConfig returns new Columns with default length n using a config struct.
func NewColumnsFromConfig(config Config, n int) (Columns, error) {
	var columns Columns

	// both nil: return default index
	if len(config.Col) == 0 && len(config.MultiCol) == 0 {
		cols := NewDefaultColLevel(n, config.ColName)
		return NewColumns(cols), nil

	}
	// both not nil: return error
	if len(config.Col) != 0 && len(config.MultiCol) != 0 {
		return Columns{}, fmt.Errorf("columnFactory(): supplying both config.Col and config.MultiCol is ambiguous; supply one or the other")
	}
	// single-level Columns
	if len(config.Col) != 0 {
		newLevel := NewColLevel(config.Col, config.ColName)
		columns = NewColumns(newLevel)
	}

	// multi-level Columns
	if len(config.MultiCol) != 0 {
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
func (col Columns) Len() int {
	if col.NumLevels() == 0 {
		return 0
	}
	return col.Levels[0].Len()
}

// NumLevels returns the number of column levels.
func (col Columns) NumLevels() int {
	return len(col.Levels)
}

// MultiNames returns a slice of the names of the column at every level for every column position.
func (col Columns) MultiNames() [][]string {
	names := make([][]string, col.Len())
	for k := 0; k < col.Len(); k++ {
		nameSlice := make([]string, col.NumLevels())
		for j := 0; j < col.NumLevels(); j++ {
			nameSlice[j] = col.Levels[j].Labels[k]
		}
		names[k] = nameSlice
	}
	return names
}

// Names returns the name of every column by concatenating the labels across every level.
func (col Columns) Names() []string {
	names := make([]string, col.Len())
	for i, mn := range col.MultiNames() {
		name := strings.Join(mn, values.GetMultiColNameSeparator())
		names[i] = name
	}
	return names
}

// Name returns the name of the column at col
func (col Columns) Name(column int) string {
	if col.NumLevels() == 0 {
		return ""
	}
	return col.Names()[column]
}

// MultiName returns the names of the column levels at col
func (col Columns) MultiName(column int) []string {
	if col.NumLevels() == 0 {
		return []string{}
	}
	return col.MultiNames()[column]
}

// MaxNameWidth returns the number of characters in the column name with the most characters.
func (col Columns) MaxNameWidth() int {
	var max int
	for k := range col.NameMap {
		if length := len(fmt.Sprint(k)); length > max {
			max = length
		}
	}
	return max
}

// UpdateNameMap updates the holistic index map of {index level names: [index level positions]}
func (col *Columns) updateNameMap() {
	nameMap := make(LabelMap)
	for i, lvl := range col.Levels {
		nameMap[lvl.Name] = append(nameMap[lvl.Name], i)
	}
	col.NameMap = nameMap
}

// Refresh updates the global name map and the label mappings at every level.
// Should be called after Series selection or index modification
func (col *Columns) Refresh() {
	col.updateNameMap()
	for i := 0; i < len(col.Levels); i++ {
		col.Levels[i].Refresh()
	}
}

// [START Columns modification methods]

// returns an error if any level position does not exist
func (col Columns) ensureLevelPositions(positions []int) error {
	for _, pos := range positions {
		len := col.NumLevels()
		if pos >= len {
			return fmt.Errorf("invalid position: %d (max %v)", pos, len-1)
		}
	}
	return nil
}

// InsertLevel inserts a level into the cols and modifies the DataFrame in place.
func (col *Columns) InsertLevel(pos int, labels []string, name string) error {
	if pos > col.NumLevels() {
		return fmt.Errorf("invalid column level: %d (max: %v)", pos, col.NumLevels()-1)
	}
	lvl := NewColLevel(labels, name)
	if len(labels) != col.Len() {
		return fmt.Errorf("col.InsertLevel(): len(labels) must equal number of columns (%d != %d)",
			len(labels), col.Len())
	}
	col.Levels = append(col.Levels[:pos], append([]ColLevel{lvl}, col.Levels[pos:]...)...)
	col.Refresh()
	return nil
}

// SubsetLevels modifies the DataFrame in place with only the specified cols levels.
func (col *Columns) SubsetLevels(levelPositions []int) error {
	err := col.ensureLevelPositions(levelPositions)
	if err != nil {
		return fmt.Errorf("col.SubsetLevels(): %v", err)
	}
	if len(levelPositions) == 0 {
		return fmt.Errorf("col.SubsetLevels(): no levels provided")
	}

	levels := make([]ColLevel, 0)
	for _, position := range levelPositions {
		levels = append(levels, col.Levels[position])
	}
	col.Levels = levels
	col.Refresh()
	return nil
}

// DropLevel drops the specified cols level and modifies the DataFrame in place.
// If there is only one col level remaining, replaces with a new default col level.
func (col *Columns) DropLevel(level int) error {
	if err := col.ensureLevelPositions([]int{level}); err != nil {
		return fmt.Errorf("Columns.DropLevel(): %v", err)
	}
	if col.NumLevels() == 1 {
		col.Levels = append(col.Levels, NewDefaultColLevel(col.Len(), ""))
	}
	col.Levels = append(col.Levels[:level], col.Levels[level+1:]...)
	col.Refresh()
	return nil
}

// [END Columns modification methods]

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

// Duplicate extends the level by itself n times and modifies the level in place.
func (lvl *ColLevel) Duplicate(n int) {
	archive := make([]string, len(lvl.Labels))
	copy(archive, lvl.Labels)
	for i := 0; i < n; i++ {
		lvl.Labels = append(lvl.Labels, archive...)
		lvl.Refresh()
	}
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
func (col Columns) Copy() Columns {
	if reflect.DeepEqual(col, Columns{Levels: make([]ColLevel, 0), NameMap: make(LabelMap)}) {
		return Columns{Levels: make([]ColLevel, 0), NameMap: make(LabelMap)}
	}
	colCopy := Columns{NameMap: LabelMap{}}
	for k, v := range col.NameMap {
		colCopy.NameMap[k] = v
	}
	for i := 0; i < len(col.Levels); i++ {
		colCopy.Levels = append(colCopy.Levels, col.Levels[i].Copy())
	}
	return colCopy
}

// Subset subsets a Columns with all the column levels located at the specified integer positions and modifies the Columns in place.
func (col *Columns) Subset(colPositions []int) {
	for j := 0; j < col.NumLevels(); j++ {
		col.Levels[j].Subset(colPositions)
	}
	col.updateNameMap()
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
