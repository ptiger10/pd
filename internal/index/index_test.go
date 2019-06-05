package index

import (
	"reflect"
	"testing"
)

func Test_Copy(t *testing.T) {
	idx := New(mustCreateNewLevelFromSlice([]int{1, 2, 3}))
	copyIdx := idx.Copy()
	for i := 0; i < len(idx.Levels); i++ {
		if reflect.ValueOf(idx.Levels[i].Labels).Pointer() == reflect.ValueOf(copyIdx.Levels[i].Labels).Pointer() {
			t.Errorf("index.Copy() returned original labels at level %v, want fresh copy", i)
		}
		if reflect.ValueOf(idx.Levels[i].LabelMap).Pointer() == reflect.ValueOf(copyIdx.Levels[i].LabelMap).Pointer() {
			t.Errorf("index.Copy() returned original map at level %v, want fresh copy", i)
		}
	}
	// lvl1a := mustCreateNewLevelFromSlice([]int{1, 2, 3})
	// origIdx := New(lvl1a)
	// newIdx := idx.Copy()
	// newIdx.Levels[0] = mustCreateNewLevelFromSlice([]int{4, 5, 6})
	// newIdx.Levels[0].Labels.Set(0, 10)
	// newIdx.Levels[0].Refresh()
	// if !reflect.DeepEqual(idx, origIdx) {
	// 	t.Errorf("Index.Copy() modified in place, want copy")
	// }
}
