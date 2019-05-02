package constructors

import (
	"reflect"
	"testing"

	constructVal "github.com/ptiger10/pd/new/internal/values/constructors"
	"github.com/ptiger10/pd/new/kinds"
)

func Test_New(t *testing.T) {
	labels := constructVal.SliceInt([]int64{0, 1, 2})
	lvl := level(labels, kinds.Int, "")
	index := New(lvl)
	gotLen := len(index.Levels)
	wantLen := 1
	if gotLen != wantLen {
		t.Errorf("Returned %d index levels, want %d", gotLen, wantLen)
	}
}

func Test_NewMulti(t *testing.T) {
	lvl1 := level(
		constructVal.SliceInt([]int64{0, 1, 2}),
		kinds.Int, "")
	lvl2 := level(
		constructVal.SliceInt([]int64{100, 101, 102}),
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
	want := New(level(
		constructVal.SliceInt([]int64{0, 1, 2}),
		kinds.Int, "",
	))
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Default constructor returned %v, want %v", got, want)
	}
}
