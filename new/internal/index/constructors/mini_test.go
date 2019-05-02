package constructors

import (
	"reflect"
	"testing"

	"github.com/ptiger10/pd/new/internal/index"
	"github.com/ptiger10/pd/new/kinds"
)

func TestMini_single(t *testing.T) {
	mini := index.MiniIndex{
		Data: []int{1, 2, 3},
		Kind: kinds.Int,
		Name: "test",
	}
	got, err := IndexFromMiniIndex([]index.MiniIndex{mini})
	if err != nil {
		t.Error(err)
	}
	want := New(
		SliceInt([]int{1, 2, 3}, "test"),
	)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MiniIndex returned %v, want %v", got, want)
	}

}
