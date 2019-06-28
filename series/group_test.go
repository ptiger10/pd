package series

import (
	"testing"
)

func Test_Group_Sum(t *testing.T) {
	s, _ := New([]int{1, 2, 3, 4}, Config{Index: []int{1, 1, 2, 2}})
	g := s.GroupByIndex()
	got := g.Sum()
	want, _ := New([]float64{3, 7}, Config{Index: []int{1, 2}})
	if !Equal(got, want) {
		t.Errorf("s.GroupByIndex.Sum() returned %v, want %v", got.index, want.index)
	}
}

func Test_Group(t *testing.T) {
	type args struct {
		label string
	}
	tests := []struct {
		name string
		args args
		want *Series
	}{
		{name: "pass", args: args{"1"}, want: MustNew([]int{1, 2}, Config{Index: []int{1, 1}})},
		{name: "fail", args: args{"100"}, want: newEmptySeries()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MustNew([]int{1, 2, 3, 4}, Config{Index: []int{1, 1, 2, 2}})
			g := s.GroupByIndex()
			got := g.Group(tt.args.label)
			if !Equal(got, tt.want) {
				t.Errorf("Grouping.Group() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Group_Sum_extended(t *testing.T) {
	s, _ := New([]int{1, 2, 3, 4, 5, 6, 7, 8}, Config{
		MultiIndex: []interface{}{[]string{"foo", "foo", "bar", "bar", "foo", "foo", "bar", "bar"}, []int{1, 2, 1, 2, 1, 2, 1, 2}},
	})
	g := s.GroupByIndex()
	got := g.Sum()
	want, _ := New([]float64{10, 12, 6, 8}, Config{
		MultiIndex: []interface{}{[]string{"bar", "bar", "foo", "foo"}, []int{1, 2, 1, 2}},
	})
	if !Equal(got, want) {
		t.Errorf("s.GroupByIndex.Sum() returned %v, want %v", got, want)
	}
}
