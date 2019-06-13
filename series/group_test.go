package series

import (
	"testing"
)

func Test_Group_Sum(t *testing.T) {
	s, _ := New([]int{1, 2, 3, 4}, Idx([]int{1, 1, 2, 2}))
	g := s.GroupByIndex()
	got := g.Sum()
	want, _ := New([]float64{3, 7}, Idx([]int{1, 2}))
	if !seriesEquals(got, want) {
		t.Errorf("s.GroupByIndex.Sum() returned %v, want %v", got.index, want.index)
	}
}

func Test_Group_Sum_extended(t *testing.T) {
	s, _ := New([]int{1, 2, 3, 4, 5, 6, 7, 8}, Idx([]string{"foo", "foo", "bar", "bar", "foo", "foo", "bar", "bar"}), Idx([]int{1, 2, 1, 2, 1, 2, 1, 2}))
	g := s.GroupByIndex()
	got := g.Sum()
	want, _ := New([]float64{10, 12, 6, 8}, Idx([]string{"bar", "bar", "foo", "foo"}), Idx([]int{1, 2, 1, 2}))
	if !seriesEquals(got, want) {
		t.Errorf("s.GroupByIndex.Sum() returned %v, want %v", got, want)
	}
}
