package index

import (
	"fmt"

	"github.com/ptiger10/pd/internal/values"
)

// Columns is a collection of column levels, plus name mappings.
type Columns struct {
	NameMap LabelMap
	Levels  []ColLevel
}

// NewColumns returns a new Columns collection from a slice of column levels.
func NewColumns(levels ...ColLevel) Columns {
	if levels == nil {
		emptyLevel := NewColLevel(nil, "")
		levels = append(levels, emptyLevel)
	}
	cols := Columns{
		Levels: levels,
	}
	cols.updateNameMap()
	return cols
}

// NewDefaultColumns returns a new Columns collection with default range labels (0, 1, 2, ... n).
func NewDefaultColumns(n int) Columns {
	return NewColumns(NewDefaultColLevel(n))
}

// Len returns the number of labels in every level of the column.
func (cols Columns) Len() int {
	if len(cols.Levels) == 0 {
		return 0
	} else {
		return cols.Levels[0].Len()
	}
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
	if cols.Len() == 0 {
		return
	}
	cols.updateNameMap()
	for i := 0; i < len(cols.Levels); i++ {
		cols.Levels[i].Refresh()
	}
}

// A ColLevel is a single collection of column labels within a Columns collection, plus label mappings and metadata.
// It is identical to an index Level except for the Labels, which are a simple []interface{} that do not satisfy the values.Values interface.
type ColLevel struct {
	Name     string
	Labels   []interface{}
	LabelMap LabelMap
}

// NewDefaultColLevel returns []string values {"0", "1", "2", ... n} for use in default DataFrame columns.
func NewDefaultColLevel(n int) ColLevel {
	colsInt := values.MakeInterfaceRange(0, n)
	return NewColLevel(colsInt, "")
}

// NewColLevel returns a Columns level with updated label map.
func NewColLevel(labels []interface{}, name string) ColLevel {
	lvl := ColLevel{
		Labels: labels,
		Name:   name,
	}
	lvl.Refresh()
	return lvl
}

// Len returns the number of labels in the column level.
func (lvl ColLevel) Len() int {
	return len(lvl.Labels)
}

// Refresh updates all the label mappings value within a column level.
func (lvl ColLevel) Refresh() {
	if lvl.Labels == nil {
		return
	}
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
