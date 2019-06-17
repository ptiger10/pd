package series

// import (
// 	"testing"
// )

// func Test_Group_Sum(t *testing.T) {
// 	s, _ := NewWithConfig([]int{1, 2, 3, 4}, Config{Index: []int{1, 1, 2, 2}})
// 	g := s.GroupByIndex()
// 	got := g.Sum()
// 	want, _ := NewWithConfig([]float64{3, 7}, Config{Index: []int{1, 2}})
// 	if !Equal(got, want) {
// 		t.Errorf("s.GroupByIndex.Sum() returned %v, want %v", got.index, want.index)
// 	}
// }

// func Test_Group_Sum_extended(t *testing.T) {
// 	s, _ := NewWithConfig([]int{1, 2, 3, 4, 5, 6, 7, 8}, Config{
// 		MultiIndex: []interface{}{[]string{"foo", "foo", "bar", "bar", "foo", "foo", "bar", "bar"}, []int{1, 2, 1, 2, 1, 2, 1, 2}},
// 	})
// 	g := s.GroupByIndex()
// 	got := g.Sum()
// 	want, _ := NewWithConfig([]float64{10, 12, 6, 8}, Config{
// 		MultiIndex: []interface{}{[]string{"bar", "bar", "foo", "foo"}, []int{1, 2, 1, 2}},
// 	})
// 	if !Equal(got, want) {
// 		t.Errorf("s.GroupByIndex.Sum() returned %v, want %v", got, want)
// 	}
// }
