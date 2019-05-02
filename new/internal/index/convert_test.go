package index_test

import (
	"fmt"
	"testing"

	constructIdx "github.com/ptiger10/pd/new/internal/index/constructors"
	"github.com/ptiger10/pd/new/kinds"
)

func TestConvertIndex_float(t *testing.T) {
	var err error
	lvl := constructIdx.SliceInt([]int{1, 2, 3}, "")
	fmt.Println(lvl)
	lvl, err = lvl.Convert(kinds.Float)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(lvl)
}
