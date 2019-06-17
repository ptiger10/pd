package index

import (
	"reflect"
	"testing"
)

func Test_New_Empty(t *testing.T) {
	idx := New()
	if idx.Len() != 0 {
		t.Error("Len() of empty index did not return 0")
	}
}

func Test_NewDefault(t *testing.T) {
	got := NewDefault(3)
	lvl, err := NewLevel([]int64{0, 1, 2}, "")
	if err != nil {
		t.Error(err)
	}
	want := New(lvl)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Default constructor returned %v, want %v", got, want)
	}
	gotLen := len(got.Levels)
	wantLen := 1
	if gotLen != wantLen {
		t.Errorf("Returned %d index levels, want %d", gotLen, wantLen)
	}
}

func Test_NewMulti(t *testing.T) {
	lvl1, err := NewLevel([]int64{0, 1, 2}, "")
	if err != nil {
		t.Error(err)
	}
	lvl2, err := NewLevel([]int64{100, 101, 102}, "")
	if err != nil {
		t.Error(err)
	}
	index := New(lvl1, lvl2)
	gotLen := len(index.Levels)
	wantLen := 2
	if gotLen != wantLen {
		t.Errorf("Returned %d index levels, want %d", gotLen, wantLen)
	}

}

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
