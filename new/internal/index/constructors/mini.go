package constructors

import (
	"fmt"
	"reflect"

	"github.com/ptiger10/pd/new/internal/index"
)

func IndexFromMiniIndex(minis []index.MiniIndex) (index.Index, error) {
	var levels []index.Level
	for _, miniIdx := range minis {
		if reflect.ValueOf(miniIdx.Data).Kind() != reflect.Slice {
			return index.Index{}, fmt.Errorf("Unable to construct index: custom index must be a Slice: unsupported index type: %T", miniIdx.Data)
		}
		level, err := LevelFromSlice(miniIdx)
		if err != nil {
			return index.Index{}, fmt.Errorf("Unable to construct index: %v", err)
		}
		if miniIdx.Kind != nil {
			level.Convert(miniIdx.Kind)
		}
		levels = append(levels, level)
		}
	}
	idx := New(levels)
	return idx, nil

}
