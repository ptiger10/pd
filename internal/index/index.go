// Package index is an internal collection of data structures and methods for constructing and modifying n-level indexes and column headers.
package index

import (
	"fmt"
	"log"
	"reflect"

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
	DataType  options.DataType
	Labels    values.Values
	LabelMap  LabelMap
	Name      string
	IsDefault bool
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
	Col             []string
	ColName         string
	MultiCol        [][]string
	MultiColNames   []string
	Manual          bool
}

// [START constructors]

// New receives one or more Levels and returns a new Index.
// Expects that Levels already have .LabelMap and .Longest set.
func New(levels ...Level) Index {
	if levels == nil {
		return Index{Levels: []Level{}, NameMap: LabelMap{}}
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
	// both nil: return default index
	if config.Index == nil && len(config.MultiIndex) == 0 {
		lvl := NewDefaultLevel(n, config.IndexName)
		return New(lvl), nil
	}
	// both not nil: return error
	if config.Index != nil && len(config.MultiIndex) != 0 {
		return Index{}, fmt.Errorf("internal/index.NewFromConfig(): supplying both config.Index and config.MultiIndex is ambiguous; supply one or the other")
	}
	// multi index
	if len(config.MultiIndex) != 0 {
		// name misalignment
		if config.MultiIndexNames != nil && len(config.MultiIndexNames) != len(config.MultiIndex) {
			return Index{}, fmt.Errorf(
				"internal/index.NewFromConfig(): if MultiIndexNames is not nil, it must must have same length as MultiIndex: %d != %d",
				len(config.MultiIndexNames), len(config.MultiIndex))
		}
		multi, err := CreateMultiIndex(config.MultiIndex, config.MultiIndexNames, config.Manual)
		if err != nil {
			return Index{}, fmt.Errorf("internal/index.NewFromConfig(): %v", err)
		}
		return multi, nil
	}
	// default: single index
	var newLevel Level
	var err error
	if !config.Manual {
		newLevel, err = InterpolatedNewLevel(config.Index, config.IndexName)
	} else {
		newLevel, err = NewLevel(config.Index, config.IndexName)
	}
	if err != nil {
		return Index{}, fmt.Errorf("internal/index.NewFromConfig(): %v", err)
	}
	return New(newLevel), nil
}

// CreateMultiIndex creates a multiindex from a slice of interface values representing different levels
func CreateMultiIndex(levels []interface{}, names []string, manualMode bool) (Index, error) {
	var newLevels []Level
	for i := 0; i < len(levels); i++ {
		var levelName string
		if i < len(names) {
			levelName = names[i]
		}
		var newLevel Level
		var err error
		if !manualMode {
			newLevel, err = InterpolatedNewLevel(levels[i], levelName)
		} else {
			newLevel, err = NewLevel(levels[i], levelName)
		}
		if err != nil {
			return Index{}, fmt.Errorf("CreateMultiIndex(): %v", err)
		}
		newLevels = append(newLevels, newLevel)
	}
	return New(newLevels...), nil
}

// Copy returns a deep copy of each index level
func (idx Index) Copy() Index {
	if reflect.DeepEqual(idx, Index{NameMap: LabelMap{}, Levels: []Level{}}) {
		return Index{NameMap: LabelMap{}, Levels: []Level{}}
	}
	idxCopy := Index{NameMap: LabelMap{}}
	for k, v := range idx.NameMap {
		idxCopy.NameMap[k] = v
	}
	for i := 0; i < len(idx.Levels); i++ {
		idxCopy.Levels = append(idxCopy.Levels, idx.Levels[i].Copy())
	}
	return idxCopy
}

// NewDefaultLevel creates an index level with range labels (0, 1, 2, ...n) and optional name.
func NewDefaultLevel(n int, name string) Level {
	vals := values.MakeIntRange(0, n)
	container := values.MustCreateValuesFromInterface(vals)
	lvl := Level{Labels: container.Values, DataType: container.DataType, Name: name, IsDefault: true}
	lvl.Refresh()
	return lvl
}

// InterpolatedNewLevel creates an Index Level and interpolates []interface{}
func InterpolatedNewLevel(data interface{}, name string) (Level, error) {
	if data == nil {
		return Level{}, nil
	}
	container, err := values.InterfaceFactory(data)
	if err != nil {
		return Level{}, fmt.Errorf("NewLevel(): %v", err)
	}
	if vals, isInterfaceSlice := data.([]interface{}); isInterfaceSlice {
		interpolateAs := values.Interpolate(vals)
		if interpolateAs != options.Interface {
			// ducks error because interpolation is controlled
			container.Values, _ = values.Convert(container.Values, interpolateAs)
		}
		container.DataType = interpolateAs
	}

	lvl := Level{Labels: container.Values, DataType: container.DataType, Name: name}
	lvl.Refresh()
	return lvl, nil
}

// NewLevel creates an Index Level from a Scalar or Slice interface{} but returns an error if interface{} is not supported by factory.
func NewLevel(data interface{}, name string) (Level, error) {
	if data == nil {
		return Level{}, nil
	}
	container, err := values.InterfaceFactory(data)
	if err != nil {
		return Level{}, fmt.Errorf("NewLevel(): %v", err)
	}
	lvl := Level{Labels: container.Values, DataType: container.DataType, Name: name}
	lvl.Refresh()
	return lvl, nil
}

// MustNewLevel returns a new level from an interface, but panics on error
func MustNewLevel(data interface{}, name string) Level {
	lvl, err := NewLevel(data, name)
	if err != nil {
		if options.GetLogWarnings() {
			log.Printf("MustNewLevel returned an error: %v", err)
		}
		lvl, _ := NewLevel(nil, "")
		return lvl
	}
	return lvl
}

// Copy copies an Index Level
func (lvl Level) Copy() Level {
	if lvl.Labels == nil {
		return Level{}
	}
	lvlCopy := Level{}
	lvlCopy = lvl
	lvlCopy.Labels = lvlCopy.Labels.Copy()
	lvlCopy.LabelMap = make(LabelMap)
	for k, v := range lvl.LabelMap {
		lvlCopy.LabelMap[k] = v
	}
	return lvlCopy
}

// [END Constructors]

// [START Element]

// Elements refer to all the elements at the same position across all levels of an index.
type Elements struct {
	Labels    []interface{}
	DataTypes []options.DataType
}

// Elements returns all the index elements at an integer position.
func (idx Index) Elements(position int) Elements {
	var labels []interface{}
	var datatypes []options.DataType
	for _, lvl := range idx.Levels {
		label := lvl.Labels.Element(position).Value
		labels = append(labels, label)
		datatypes = append(datatypes, lvl.DataType)
	}
	return Elements{labels, datatypes}
}

// [END Element]

// [START Index]

// Len returns the number of labels in every level of the index.
func (idx Index) Len() int {
	if idx.NumLevels() == 0 {
		return 0
	}
	return idx.Levels[0].Len()
}

// NumLevels returns the number of levels in the index.
func (idx Index) NumLevels() int {
	return len(idx.Levels)
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
	maxWidths := make([]int, idx.NumLevels())
	for j := 0; j < idx.NumLevels(); j++ {
		maxWidths[j] = idx.Levels[j].maxWidth()
	}
	return maxWidths
}

// DataTypes returns a slice of the DataTypes at each level of the index
func (idx Index) DataTypes() []options.DataType {
	idxDataTypes := make([]options.DataType, idx.NumLevels())
	for j := 0; j < idx.NumLevels(); j++ {
		idxDataTypes[j] = idx.Levels[j].DataType
	}
	return idxDataTypes
}

// Aligned ensures that all index levels have the same length.
func (idx Index) Aligned() error {
	if idx.NumLevels() == 0 {
		return nil
	}
	lvl0 := idx.Levels[0].Len()
	for i := 1; i < idx.NumLevels(); i++ {
		if cmpLvl := idx.Levels[i].Len(); lvl0 != cmpLvl {
			return fmt.Errorf("index.Aligned(): index level %v must have same number of labels as level 0, %d != %d",
				i, cmpLvl, lvl0)
		}
	}
	return nil
}

// SwapLevels swaps two levels in the index and modifies the index in place.
func (idx *Index) SwapLevels(i, j int) error {
	if err := idx.ensureLevelPositions([]int{i}); err != nil {
		return fmt.Errorf("invalid i: %v", err)
	}
	if err := idx.ensureLevelPositions([]int{j}); err != nil {
		return fmt.Errorf("invalid j: %v", err)
	}
	idx.Levels[i], idx.Levels[j] = idx.Levels[j], idx.Levels[i]
	idx.Refresh()
	return nil
}

// InsertLevel inserts a level into the index and modifies the index in place.
func (idx *Index) InsertLevel(pos int, values interface{}, name string) error {
	if pos > idx.NumLevels() {
		return fmt.Errorf("invalid index level: %d (max: %v)", pos, idx.NumLevels()-1)
	}
	lvl, err := NewLevel(values, name)
	if err != nil {
		return fmt.Errorf("index.InsertLevel(): %v", err)
	}
	if lvl.Len() != idx.Len() {
		return fmt.Errorf("values to insert must have same length as existing index: %d != %d", lvl.Len(), idx.Len())
	}
	idx.Levels = append(idx.Levels[:pos], append([]Level{lvl}, idx.Levels[pos:]...)...)
	idx.Refresh()
	return nil
}

// Subset subsets the index with the labels located at the specified integer positions and modifies the index in place.
func (idx *Index) Subset(rowPositions []int) error {
	err := idx.ensureRowPositions(rowPositions)
	if err != nil {
		return fmt.Errorf("index.SubsetLevels(): %v", err)
	}
	if len(rowPositions) == 0 {
		return fmt.Errorf("index.SubsetLevels(): no levels provided")
	}

	for j := 0; j < idx.NumLevels(); j++ {
		idx.Levels[j].Labels = idx.Levels[j].Labels.Subset(rowPositions)
	}
	idx.Refresh()
	return nil
}

// SubsetLevels subsets the index with only the levels located at specified integer positions and modifies the index in place.
func (idx *Index) SubsetLevels(levelPositions []int) error {
	err := idx.ensureLevelPositions(levelPositions)
	if err != nil {
		return fmt.Errorf("index.SubsetLevels(): %v", err)
	}
	if len(levelPositions) == 0 {
		return fmt.Errorf("index.SubsetLevels(): no levels provided")
	}

	newIdx := Index{Levels: make([]Level, len(levelPositions))}
	for i, pos := range levelPositions {
		newIdx.Levels[i] = idx.Levels[pos]
	}
	newIdx.updateNameMap()
	*idx = newIdx
	return nil
}

func (idx Index) ensureRowPositions(rowPositions []int) error {
	len := idx.Len()
	for _, pos := range rowPositions {
		if pos >= len {
			return fmt.Errorf("invalid position: %d (max %v)", pos, len-1)
		}
	}
	return nil
}

func (idx Index) ensureLevelPositions(levelPositions []int) error {
	for _, pos := range levelPositions {
		if pos >= idx.NumLevels() {
			return fmt.Errorf("invalid index level: %d (max: %v)", pos, idx.NumLevels()-1)
		}
	}
	return nil
}

// Set sets the value at the specified index row and level to val and modifies the Index in place.
func (idx *Index) Set(row int, level int, val interface{}) error {
	if err := idx.ensureRowPositions([]int{row}); err != nil {
		return fmt.Errorf("index.Set(): %v", err)
	}
	if err := idx.ensureLevelPositions([]int{level}); err != nil {
		return fmt.Errorf("index.Set(): %v", err)
	}
	if _, err := values.InterfaceFactory(val); err != nil {
		return fmt.Errorf("index.Set(): %v", err)
	}

	idx.Levels[level].Labels.Set(row, val)
	idx.Levels[level].Refresh()
	return nil
}

// dropLevel drops an index level.
func (idx *Index) dropLevel(level int) {
	if idx.NumLevels() == 1 {
		idx.Levels = append(idx.Levels, NewDefaultLevel(idx.Len(), ""))
	}
	idx.Levels = append(idx.Levels[:level], idx.Levels[level+1:]...)
	return
}

// DropLevel drops an index level and modifies the Index in place. If there is only one level, replaces with default index.
func (idx *Index) DropLevel(level int) error {
	if err := idx.ensureLevelPositions([]int{level}); err != nil {
		return fmt.Errorf("index.DropLevel(): %v", err)
	}
	idx.dropLevel(level)
	idx.Refresh()
	return nil
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
// Should be called after index modification
func (idx *Index) Refresh() {
	idx.updateNameMap()
	for i := 0; i < idx.NumLevels(); i++ {
		idx.Levels[i].Refresh()
	}
}

// [END Index]

// [START index level]

// Len returns the number of labels in the level
func (lvl Level) Len() int {
	if reflect.DeepEqual(lvl, Level{}) {
		return 0
	}
	return lvl.Labels.Len()
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
	lvl.updateLabelMap()
}

// [END index level]

// [START inplace conversion]

// Convert an index level from one kind to another in place, then refreshes the LabelMap
func (lvl *Level) Convert(kind options.DataType) error {
	switch kind {
	case options.None:
		return fmt.Errorf("unable to convert index level: must supply a valid Kind")
	case options.Float64:
		lvl.ToFloat64()
	case options.Int64:
		lvl.ToInt64()
	case options.String:
		lvl.ToString()
	case options.Bool:
		lvl.ToBool()
	case options.DateTime:
		lvl.ToDateTime()
	case options.Interface:
		lvl.ToInterface()
	default:
		return fmt.Errorf("unable to convert level: kind not supported: %v", kind)
	}
	return nil
}

// ToFloat64 converts an index level to Float
func (lvl *Level) ToFloat64() {
	lvl.Labels = lvl.Labels.ToFloat64()
	lvl.DataType = options.Float64
	lvl.Refresh()
	return
}

// ToInt64 converts an index level to Int
func (lvl *Level) ToInt64() {
	lvl.Labels = lvl.Labels.ToInt64()
	lvl.DataType = options.Int64
	lvl.Refresh()
}

// ToString converts an index level to String
func (lvl *Level) ToString() {
	lvl.Labels = lvl.Labels.ToString()
	lvl.DataType = options.String
	lvl.Refresh()
}

// ToBool converts an index level to Bool
func (lvl *Level) ToBool() {
	lvl.Labels = lvl.Labels.ToBool()
	lvl.DataType = options.Bool
	lvl.Refresh()
}

// ToDateTime converts an index level to DateTime
func (lvl *Level) ToDateTime() {
	lvl.Labels = lvl.Labels.ToDateTime()
	lvl.DataType = options.DateTime
	lvl.Refresh()
}

// ToInterface converts an index level to Interface
func (lvl *Level) ToInterface() {
	lvl.Labels = lvl.Labels.ToInterface()
	lvl.DataType = options.Interface
	lvl.Refresh()
}

// [END conversion]
