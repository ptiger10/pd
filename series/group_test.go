package series

import (
	"testing"
)

func Test_Group_Sum(t *testing.T) {
	s, _ := NewPointer([]int{1, 2, 3, 4}, Idx([]int{1, 1, 2, 2}))
	g := s.GroupByIndex()
	got := g.Sum()
	want := mustNew([]float64{3, 7}, Idx([]int{1, 2}))
	if !seriesEquals(got, want) {
		t.Errorf("s.GroupByIndex.Sum() returned %v, want %v", got.index, want.index)
	}
}
