package index

import (
	"fmt"
)

func (lvl *Level) UpdateLabelMap() {
	labelMap := make(LabelMap)
	for i, val := range lvl.Labels.All() {
		key := fmt.Sprint(val)
		labelMap[key] = append(labelMap[key], i)
	}
	lvl.LabelMap = labelMap
}

func (idx *Index) UpdateNameMap() {
	nameMap := make(LabelMap)
	for i, lvl := range idx.Levels {
		nameMap[lvl.Name] = append(nameMap[lvl.Name], i)
	}
}
