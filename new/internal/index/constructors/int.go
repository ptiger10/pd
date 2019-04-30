package constructors

import (
	"reflect"
	"strconv"

	"github.com/ptiger10/pd/new/datatypes"
	"github.com/ptiger10/pd/new/internal/index"
)

// SliceInt converts []int (of any int type) -> IndexLevel with IntLabels
func SliceInt(data interface{}) index.Level {
	var labels index.IntLabels
	labelMap := make(index.LabelMap)
	d := reflect.ValueOf(data)
	for i := 0; i < d.Len(); i++ {
		label := d.Index(i).Int()
		labels = append(labels, label)
		strLabel := strconv.Itoa(int(label))
		labelMap[strLabel] = append(labelMap[strLabel], i)
	}
	ret := index.Level{
		Kind:     datatypes.Int,
		Labels:   labels,
		LabelMap: labelMap,
	}
	return ret

}
