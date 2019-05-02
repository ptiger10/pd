package constructors

import (
	"github.com/ptiger10/pd/new/internal/index"
)

// LevelFromSlice creates an Index Level from an interface{} that reflects Slice
func LevelFromSlice(miniIdx index.MiniIndex) index.Level {
	switch miniIdx.Data.(type) {
	case []int, []int8, []int16, []int32, []int64:
		return SliceInt(miniIdx.Data, miniIdx.Name)
	default:
		return index.Level{}
	}
}
