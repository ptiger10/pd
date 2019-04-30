package constructors

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ptiger10/pd/new/internal/index"
)

func Test_New(t *testing.T) {
	lvl1 := SliceInt([]int{0, 1, 2})
	index := New([]index.Level{lvl1})
	gotLen := len(index.Levels)
	wantLen := 1
	if gotLen != wantLen {
		t.Errorf("Returned %d index levels, want %d", gotLen, wantLen)
	}
}
func Test_NewMulti(t *testing.T) {
	lvl1 := SliceInt([]int{0, 1, 2})
	lvl2 := SliceInt([]int{100, 101, 102})
	index := New([]index.Level{lvl1, lvl2})
	gotLen := len(index.Levels)
	wantLen := 2
	if gotLen != wantLen {
		t.Errorf("Returned %d index levels, want %d", gotLen, wantLen)
	}
	idx := reflect.ValueOf(index.Levels)
	fmt.Println(idx)
	// for i := 0; i < vals.Len(); i++ {
	for j := 0; j < idx.Len(); j++ {
		fmt.Printf("%v ", idx.Index(j).FieldByName("Labels"))
	}

}
