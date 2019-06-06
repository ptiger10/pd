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
	err := idx.Drop(0)
	if err != nil {
		t.Errorf("idx.Drop(): %v", err)
	}
	want := New(MustCreateNewLevel([]int{1, 2, 3}, ""))
	if !reflect.DeepEqual(idx, want) {
		t.Errorf("idx.Drop() for one level returned %v, want %v", idx, want)
	}
}

func Test_Drop_multilevel(t *testing.T) {
	idx := New(MustCreateNewLevel([]int{1, 2, 3}, ""), MustCreateNewLevel([]int{4, 5, 6}, ""))
	idx.Drop(1)
	want := New(MustCreateNewLevel([]int{1, 2, 3}, ""))
	if !reflect.DeepEqual(idx, want) {
		t.Errorf("idx.Drop() for multilevel returned %v, want %v", idx, want)
	}
}

func Test_Droplevels(t *testing.T) {
	idx := New(MustCreateNewLevel([]int{1, 2, 3}, ""), MustCreateNewLevel([]int{4, 5, 6}, ""), MustCreateNewLevel([]int{7, 8, 9}, ""))
	err := idx.dropLevels([]int{2, 0})
	if err != nil {
		t.Errorf("idx.Droplevels(): %v", err)
	}
	want := New(MustCreateNewLevel([]int{4, 5, 6}, ""))
	if !reflect.DeepEqual(idx, want) {
		t.Errorf("idx.Drop() for multilevel returned %v, want %v", idx, want)
	}
}
