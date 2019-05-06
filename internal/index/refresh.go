package index

import (
	"fmt"

	"github.com/ptiger10/pd/options"
)

// UpdateLabelMap updates a single level's map of {label values: [label positions]}
func (lvl *Level) UpdateLabelMap() {
	labels := lvl.Labels.All()
	labelMap := make(LabelMap, len(labels))
	for i, val := range labels {
		key := fmt.Sprint(val)
		labelMap[key] = append(labelMap[key], i)
	}
	lvl.LabelMap = labelMap
}

// UpdateNameMap updates the holistic index map of {index level names: [index level positions]}
func (idx *Index) UpdateNameMap() {
	nameMap := make(LabelMap)
	for i, lvl := range idx.Levels {
		nameMap[lvl.Name] = append(nameMap[lvl.Name], i)
	}
	idx.NameMap = nameMap
}

// UpdateLongest finds the max length of either the level name or the longest string in the LabelMap,
// for use in printing a Series or DataFrame
func (lvl *Level) UpdateLongest() {
	var max int
	for k := range lvl.LabelMap {
		if len(k) > max {
			max = len(k)
		}
	}
	if len(lvl.Name) > max {
		max = len(lvl.Name)
	}
	if max > options.DisplayIndexMaxWidth {
		max = options.DisplayIndexMaxWidth
	}
	lvl.Longest = max
}

// Refresh updates all the label mappings and metadata within the index.
// Should be called after Series selection or index modification
func (idx *Index) Refresh() {
	idx.UpdateNameMap()
	for i := 0; i < len(idx.Levels); i++ {
		idx.Levels[i].Refresh()
	}
}

// Refresh updates all the label mappings and metadata within a level.
func (lvl *Level) Refresh() {
	lvl.UpdateLabelMap()
	lvl.UpdateLongest()
}
