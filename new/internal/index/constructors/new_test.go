package constructors

import (
	"reflect"
	"testing"

	constructVal "github.com/ptiger10/pd/new/internal/values/constructors"
	"github.com/ptiger10/pd/new/kinds"
)

func Test_New(t *testing.T) {
	labels := constructVal.SliceInt([]int{0, 1, 2})
	lvl := Level(labels, kinds.Int, "")
	index := New(lvl)
	gotLen := len(index.Levels)
	wantLen := 1
	if gotLen != wantLen {
		t.Errorf("Returned %d index levels, want %d", gotLen, wantLen)
	}
}

func Test_NewMulti(t *testing.T) {
	lvl1 := Level(
		constructVal.SliceInt([]int{0, 1, 2}),
		kinds.Int, "")
	lvl2 := Level(
		constructVal.SliceInt([]int{100, 101, 102}),
		kinds.Int, "")
	index := New(lvl1, lvl2)
	gotLen := len(index.Levels)
	wantLen := 2
	if gotLen != wantLen {
		t.Errorf("Returned %d index levels, want %d", gotLen, wantLen)
	}

}

func Test_Default(t *testing.T) {
	got := Default(3)
	want := New(Level(
		constructVal.SliceInt([]int{0, 1, 2}),
		kinds.Int, "",
	))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Default constructor returned %v, want %v", got, want)
	}
}
