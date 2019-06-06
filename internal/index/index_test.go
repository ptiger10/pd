package index

import (
	"reflect"
	"testing"
)

func Test_Copy(t *testing.T) {
	idx := New(MustCreateNewLevel([]int{1, 2, 3}, ""))
	copyIdx := idx.Copy()
	for i := 0; i < len(idx.Levels); i++ {
		if reflect.ValueOf(idx.Levels[i].Labels).Pointer() == reflect.ValueOf(copyIdx.Levels[i].Labels).Pointer() {
			t.Errorf("index.Copy() returned original labels at level %v, want fresh copy", i)
		}
		if reflect.ValueOf(idx.Levels[i].LabelMap).Pointer() == reflect.ValueOf(copyIdx.Levels[i].LabelMap).Pointer() {
			t.Errorf("index.Copy() returned original map at level %v, want fresh copy", i)
		}
	}
}

func Test_Drop_oneLevel(t *testing.T) {
	idx := New(MustCreateNewLevel([]int{1, 2, 3}, ""))
	idx.Drop(0)
	want := New(MustCreateNewLevel([]int{1, 2, 3}, ""))
	if !reflect.DeepEqual(idx, want) {
		t.Errorf("idx.Drop() for one level returned %v, want %v", idx, want)
	}
}

func Test_Drop_twoLevels(t *testing.T) {
	idx := New(MustCreateNewLevel([]int{1, 2, 3}, ""), MustCreateNewLevel([]int{4, 5, 6}, ""))
	idx.Drop(1)
	want := New(MustCreateNewLevel([]int{1, 2, 3}, ""))
	if !reflect.DeepEqual(idx, want) {
		t.Errorf("idx.Drop() for two levels returned %v, want %v", idx, want)
	}
}
