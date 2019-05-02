package constructors

import (
	"fmt"
	"reflect"

	"github.com/ptiger10/pd/new/internal/index"
	"github.com/ptiger10/pd/new/kinds"
)

// IndexFromMiniIndex creates a full index from a mini representation of an index level
func IndexFromMiniIndex(minis []index.MiniIndex) (index.Index, error) {
	var levels []index.Level
	for _, miniIdx := range minis {
		if reflect.ValueOf(miniIdx.Data).Kind() != reflect.Slice {
			return index.Index{}, fmt.Errorf("Unable to construct index: custom index must be a Slice: unsupported index type: %T", miniIdx.Data)
		}
		level, err := LevelFromSlice(miniIdx.Data, miniIdx.Name)
		if err != nil {
			return index.Index{}, fmt.Errorf("Unable to construct index: %v", err)
		}
		if miniIdx.Kind != kinds.None {
			level.Convert(miniIdx.Kind)
		}
		levels = append(levels, level)
	}
	idx := New(levels...)
	return idx, nil

}
