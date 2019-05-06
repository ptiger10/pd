package index

import (
	"reflect"
	"testing"
)

func Test_Default(t *testing.T) {
	got := Default(3)
	lvl, err := NewLevelFromSlice([]int64{0, 1, 2}, "")
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
	lvl1, err := NewLevelFromSlice([]int64{0, 1, 2}, "")
	if err != nil {
		t.Error(err)
	}
	lvl2, err := NewLevelFromSlice([]int64{100, 101, 102}, "")
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
