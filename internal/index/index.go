package index

import (
	"fmt"
	"log"
	"sort"

	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
)

// An Index is a collection of levels, plus label mappings
type Index struct {
	Levels  []Level
	NameMap LabelMap
}

// A Level is a single collection of labels within an index, plus label mappings and metadata
type Level struct {
	DataType options.DataType
	Labels   values.Values
	LabelMap LabelMap
	Name     string
}

// A LabelMap records the position of labels, in the form {label name: [label position(s)]}
type LabelMap map[string][]int

// A Config customizes the construction of an Index or Columns object.
type Config struct {
	Name            string
	DataType        options.DataType
	Index           interface{}
	IndexName       string
	MultiIndex      []interface{}
	MultiIndexNames []string
	Cols            []interface{}
	ColsName        string
	MultiCol        [][]interface{}
	MultiColNames   []string
}

// Element is a single Element in an index level.
type Element struct {
	Label    interface{}
	DataType options.DataType
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

// NewFromConfig returns a new Index with default length n using a config struct.
func NewFromConfig(config Config, n int) (Index, error) {
	var index Index
	// both nil: return default index
	if config.Index == nil && config.MultiIndex == nil {
		lvl := NewDefaultLevel(n, config.IndexName)
		return New(lvl), nil
	}
	// both not nil: return error
	if config.Index != nil && config.MultiIndex != nil {
		return Index{}, fmt.Errorf("internal/index.NewFromConfig(): supplying both config.Index and config.MultiIndex is ambiguous; supply one or the other")
	}
	// single index
	if config.Index != nil {
		newLevel, err := NewLevel(config.Index, config.IndexName)
		if err != nil {
			return Index{}, fmt.Errorf("internal/index.NewFromConfig(): %v", err)
		}
		return New(newLevel), nil
	}
	// multi index
	if config.MultiIndex != nil {
		// name misalignment
		if config.MultiIndexNames != nil && len(config.MultiIndexNames) != len(config.MultiIndex) {
			return Index{}, fmt.Errorf(
				"internal/index.NewFromConfig(): if MultiIndexNames is not nil, it must must have same length as MultiIndex: %d != %d",
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
				return Index{}, fmt.Errorf("internal/index.NewFromConfig(): %v", err)
			}
			newLevels = append(newLevels, newLevel)
		}
		return New(newLevels...), nil
	}
	return index, nil
}

// NewDefaultLevel creates an index level with range labels (0, 1, 2, ...n) and optional name.
func NewDefaultLevel(n int, name string) Level {
	v := values.MakeIntRange(0, n)
	level := MustNewLevel(v, name)
	return level
}

// NewLevel creates an Index Level from a Scalar or Slice interface{} but returns an error if interface{} is not supported by factory.
func NewLevel(data interface{}, name string) (Level, error) {
	factory, err := values.InterfaceFactory(data)
	if err != nil {
		return Level{}, fmt.Errorf("NewLevel(): %v", err)
	}
	lvl := Level{Labels: factory.Values, DataType: factory.DataType, Name: name}
	lvl.Refresh()
	return lvl, nil
}

// MustNewLevel returns a new level from an interface, but panics on error
func MustNewLevel(data interface{}, name string) Level {
	lvl, err := NewLevel(data, name)
	if err != nil {
		log.Fatalf("MustNewLevel returned an error: %v", err)
	}
	return lvl
}

// [START Index]

// Len returns the number of labels in every level of the index.
func (idx Index) Len() int {
	if len(idx.Levels) == 0 {
		return 0
	}
	return idx.Levels[0].Len()
}

// NumLevels returns the number of levels in the index.
func (idx Index) NumLevels() int {
	return len(idx.Levels)
}

// Aligned ensures that all index levels have the same length.
func (idx Index) Aligned() error {
	lvl0 := idx.Levels[0].Len()
	for i := 1; i < idx.NumLevels(); i++ {
		if cmpLvl := idx.Levels[i].Len(); lvl0 != cmpLvl {
			return fmt.Errorf("index.Aligned(): index level %v must have same number of labels as level 0, %d != %d",
				i, cmpLvl, lvl0)
		}
	}
	return nil
}

// Unnamed returns true if all index levels are unnamed
func (idx Index) Unnamed() bool {
	for _, lvl := range idx.Levels {
		if lvl.Name != "" {
			return false
		}
	}
	return true
}

// MaxWidths returns the max number of characters in each level of an index.
func (idx Index) MaxWidths() []int {
	var maxWidths []int
	for _, lvl := range idx.Levels {
		maxWidths = append(maxWidths, lvl.maxWidth())
	}
	return maxWidths
}

// Subset returns a new index with all the labels located at the specified integer positions
func (idx Index) Subset(rowPositions []int) (Index, error) {
	var err error
	idx = idx.Copy()
	for i, level := range idx.Levels {
		idx.Levels[i].Labels, err = level.Labels.In(rowPositions)
		if err != nil {
			return Index{}, fmt.Errorf("selecting rows in index: %v", err)
		}
	}
	idx.Refresh()
	return idx, nil
}

// SubsetLevels returns a copy of the index with only those levels located at specified integer positions
func (idx Index) SubsetLevels(levelPositions []int) (Index, error) {
	idx = idx.Copy()
	var lvls []Level
	for _, pos := range levelPositions {
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

// [END Index]

// [START index level]

// Copy copies an Index Level
func (lvl Level) Copy() Level {
	lvlCopy := Level{}
	lvlCopy = lvl
	lvlCopy.Labels = lvlCopy.Labels.Copy()
	lvlCopy.LabelMap = make(LabelMap)
	for k, v := range lvl.LabelMap {
		lvlCopy.LabelMap[k] = v
	}
	return lvlCopy
}

// Len returns the number of labels in the level
func (lvl Level) Len() int {
	return lvl.Labels.Len()
}

// Element returns an Element at an integer position within an index level.
func (lvl Level) Element(position int) Element {
	return Element{
		Label:    lvl.Labels.Element(position).Value,
		DataType: lvl.DataType,
	}
}

// maxWidth finds the max length of either the level name or the longest string in the LabelMap,
// for use in printing a Series or DataFrame
func (lvl *Level) maxWidth() int {
	var max int
	for k := range lvl.LabelMap {
		if length := len(fmt.Sprint(k)); length > max {
			max = length
		}
	}
	if len(lvl.Name) > max {
		max = len(lvl.Name)
	}
	return max
}

// updateLabelMap updates a single index level's map of {label values: [label positions]}.
// A level's label map is unaware of the actual values in those positions.
func (lvl *Level) updateLabelMap() {
	labelMap := make(LabelMap, lvl.Len())
	for i := 0; i < lvl.Len(); i++ {
		key := fmt.Sprint(lvl.Labels.Element(i).Value)
		labelMap[key] = append(labelMap[key], i)
	}
	lvl.LabelMap = labelMap
}

// Refresh updates all the label mappings value within a level.
func (lvl *Level) Refresh() {
	if lvl.Labels == nil {
		return
	}
	lvl.updateLabelMap()
}

// [END index level]

// [START conversion]

// Convert an index level from one kind to another, then refresh the LabelMap
func (lvl Level) Convert(kind options.DataType) (Level, error) {
	var convertedLvl Level
	switch kind {
	case options.None:
		return Level{}, fmt.Errorf("unable to convert index level: must supply a valid Kind")
	case options.Float64:
		convertedLvl = lvl.ToFloat64()
	case options.Int64:
		convertedLvl = lvl.ToInt64()
	case options.String:
		convertedLvl = lvl.ToString()
	case options.Bool:
		convertedLvl = lvl.ToBool()
	case options.DateTime:
		convertedLvl = lvl.ToDateTime()
	case options.Interface:
		convertedLvl = lvl.ToInterface()
	default:
		return Level{}, fmt.Errorf("unable to convert level: kind not supported: %v", kind)
	}
	return convertedLvl, nil
}

// ToFloat64 converts an index level to Float
func (lvl Level) ToFloat64() Level {
	lvl.Labels = lvl.Labels.ToFloat64()
	lvl.DataType = options.Float64
	lvl.Refresh()
	return lvl
}

// ToInt64 converts an index level to Int
func (lvl Level) ToInt64() Level {
	lvl.Labels = lvl.Labels.ToInt64()
	lvl.DataType = options.Int64
	lvl.Refresh()
	return lvl
}

// ToString converts an index level to String
func (lvl Level) ToString() Level {
	lvl.Labels = lvl.Labels.ToString()
	lvl.DataType = options.String
	lvl.Refresh()
	return lvl
}

// ToBool converts an index level to Bool
func (lvl Level) ToBool() Level {
	lvl.Labels = lvl.Labels.ToBool()
	lvl.DataType = options.Bool
	lvl.Refresh()
	return lvl
}

// ToDateTime converts an index level to DateTime
func (lvl Level) ToDateTime() Level {
	lvl.Labels = lvl.Labels.ToDateTime()
	lvl.DataType = options.DateTime
	lvl.Refresh()
	return lvl
}

// ToInterface converts an index level to Interface
func (lvl Level) ToInterface() Level {
	lvl.Labels = lvl.Labels.ToInterface()
	lvl.DataType = options.Interface
	lvl.Refresh()
	return lvl
}

// [END conversion]
