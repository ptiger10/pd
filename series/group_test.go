package series

import (
	"reflect"
	"testing"

	"github.com/ptiger10/pd/internal/index"
)

func Test_Group_Sum(t *testing.T) {
	s := MustNew([]int{1, 2, 3, 4}, Config{Index: []int{1, 1, 2, 2}})
	tests := []struct {
		name  string
		input *Series
		fn    func(Grouping) *Series
		want  *Series
	}{
		{"fail: empty", newEmptySeries(), Grouping.Sum, newEmptySeries()},
		{"sum", s, Grouping.Sum,
			MustNew([]float64{3, 7}, Config{Index: []int{1, 2}})},
		{"mean", s, Grouping.Mean,
			MustNew([]float64{1.5, 3.5}, Config{Index: []int{1, 2}})},
		{"min", s, Grouping.Min,
			MustNew([]float64{1, 3}, Config{Index: []int{1, 2}})},
		{"max", s, Grouping.Max,
			MustNew([]float64{2, 4}, Config{Index: []int{1, 2}})},
		{"median", s, Grouping.Median,
			MustNew([]float64{1.5, 3.5}, Config{Index: []int{1, 2}})},
		{"standard deviation", s, Grouping.Std,
			MustNew([]float64{0.5, 0.5}, Config{Index: []int{1, 2}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.input.GroupByIndex()
			got := tt.fn(g)
			if !Equal(got, tt.want) {
				t.Errorf("s.GroupByIndex math operation returned %v, want %v", got, tt.want)
			}
		})
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

func TestSeries_GroupByIndex(t *testing.T) {
	lvl1 := index.MustNewLevel(1, "")
	lvl2 := index.MustNewLevel(2, "")
	tests := []struct {
		name  string
		input *Series
		want  map[string]*group
	}{
		{name: "single", input: MustNew([]string{"foo", "bar", "baz"}, Config{Index: []int{1, 1, 2}}),
			want: map[string]*group{
				"1": &group{Positions: []int{0, 1}, Index: index.New(lvl1)},
				"2": &group{Positions: []int{2}, Index: index.New(lvl2)},
			}},

		{name: "multi",
			input: MustNew([]string{"foo", "bar", "baz"}, Config{MultiIndex: []interface{}{[]int{1, 1, 2}, []int{2, 2, 1}}}),
			want: map[string]*group{
				"1 2": &group{Positions: []int{0, 1}, Index: index.New(lvl1, lvl2)},
				"2 1": &group{Positions: []int{2}, Index: index.New(lvl2, lvl1)},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.input.Copy()
			got := s.GroupByIndex().groups
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Series.GroupByIndex() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
