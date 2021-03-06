package dataframe

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/ptiger10/pd/options"
)

func TestGroup_Copy(t *testing.T) {
	s := MustNew([]interface{}{[]int{1, 2, 3, 4}}, Config{Index: []int{1, 1, 2, 2}})
	got := s.GroupByIndex(0).copy().groups
	want := map[string]*group{
		"1": {Positions: []int{0, 1}, FirstPosition: 0},
		"2": {Positions: []int{2, 3}, FirstPosition: 2},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("group.copy() got %v, want %v", got, want)
	}
}

func TestDataFrame_GroupByIndex(t *testing.T) {
	multi := MustNew([]interface{}{[]string{"foo", "bar", "baz"}}, Config{MultiIndex: []interface{}{[]int{1, 1, 2}, []int{2, 2, 1}}})
	type args struct {
		levelPositions []int
	}
	tests := []struct {
		name  string
		input *DataFrame
		args  args
		want  map[string]*group
	}{
		{name: "single no args",
			input: MustNew([]interface{}{[]string{"foo", "bar", "baz"}}, Config{Index: []int{1, 1, 2}}),
			args:  args{[]int{}},
			want: map[string]*group{
				"1": {Positions: []int{0, 1}, FirstPosition: 0},
				"2": {Positions: []int{2}, FirstPosition: 2},
			}},
		{"multi no args",
			multi,
			args{[]int{}},
			map[string]*group{
				"1 | 2": {Positions: []int{0, 1}, FirstPosition: 0},
				"2 | 1": {Positions: []int{2}, FirstPosition: 2},
			}},
		{"multi, select one level",
			multi,
			args{[]int{0}},
			map[string]*group{
				"1": {Positions: []int{0, 1}, FirstPosition: 0},
				"2": {Positions: []int{2}, FirstPosition: 2},
			}},
		{"multi in reverse order, select two levels",
			MustNew([]interface{}{[]string{"foo", "bar", "baz"}}, Config{MultiIndex: []interface{}{[]int{2, 2, 1}, []int{1, 1, 2}}}),
			args{[]int{0, 1}},
			map[string]*group{
				"2 | 1": {Positions: []int{0, 1}, FirstPosition: 0},
				"1 | 2": {Positions: []int{2}, FirstPosition: 2},
			}},
		{"multi - select two levels in reverse",
			multi,
			args{[]int{1, 0}},
			map[string]*group{
				"2 | 1": {Positions: []int{0, 1}, FirstPosition: 0},
				"1 | 2": {Positions: []int{2}, FirstPosition: 2},
			}},
		{"fail: invalid level",
			multi,
			args{[]int{10}},
			newEmptyGrouping().groups},
		{"fail: partial invalid level",
			multi,
			args{[]int{0, 10}},
			newEmptyGrouping().groups},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			df := tt.input.Copy()
			got := df.GroupByIndex(tt.args.levelPositions...).groups
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DataFrame.GroupByIndex() = %v, want %v", got, tt.want)
			}

			if strings.Contains(tt.name, "fail") {
				if buf.String() == "" {
					t.Errorf("DataFrame.GroupByIndex() returned no log message, want log due to fail")
				}
			}
		})
	}
}

func TestDataFrame_GroupBy(t *testing.T) {
	single := MustNew([]interface{}{[]string{"foo", "bar", "baz"}, []int{1, 1, 2}})
	multi := MustNew([]interface{}{[]string{"foo", "bar", "baz"}, []int{1, 1, 2}, []int{2, 2, 1}})
	type args struct {
		cols []int
	}
	tests := []struct {
		name  string
		input *DataFrame
		args  args
		want  map[string]*group
	}{
		{name: "single",
			input: single,
			args:  args{[]int{1}},
			want: map[string]*group{
				"1": {Positions: []int{0, 1}, FirstPosition: 0},
				"2": {Positions: []int{2}, FirstPosition: 2},
			}},
		{"multi",
			multi,
			args{[]int{1, 2}},
			map[string]*group{
				"1 | 2": {Positions: []int{0, 1}, FirstPosition: 0},
				"2 | 1": {Positions: []int{2}, FirstPosition: 2},
			}},
		{"fail: invalid level",
			single,
			args{[]int{10}},
			newEmptyGrouping().groups},
		{"fail: partial invalid level",
			single,
			args{[]int{0, 10}},
			newEmptyGrouping().groups},
		{"fail: no args provided",
			single,
			args{[]int{}},
			newEmptyGrouping().groups},
		{"fail: no columns left ungrouped",
			single,
			args{[]int{0, 1}},
			newEmptyGrouping().groups},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			df := tt.input.Copy()
			got := df.GroupBy(tt.args.cols...).groups
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DataFrame.GroupBy() = %#v, want %#v", got, tt.want)
			}

			if strings.Contains(tt.name, "fail") {
				if buf.String() == "" {
					t.Errorf("DataFrame.GroupByIndex() returned no log message, want log due to fail")
				}
			}
		})
	}
}

func Test_Groups(t *testing.T) {
	type want struct {
		groups []string
		sorted []string
	}
	tests := []struct {
		name  string
		input *DataFrame
		want  want
	}{
		// {name: "one group", input: MustNew([]interface{}{[]int{1, 2}}, Config{Index: []int{1, 1}}),
		// 	want: want{[]string{"1"}, []string{"1"}}},
		// {"two groups", MustNew([]interface{}{[]int{1, 2}}, Config{Index: []int{1, 2}}),
		// 	want{[]string{"1", "2"}, []string{"1", "2"}}},
		{"two groups - reverse order", MustNew([]interface{}{[]int{1, 2}}, Config{Index: []int{2, 1}}),
			want{[]string{"2", "1"}, []string{"1", "2"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.input.GroupByIndex()
			got := g.Groups()
			if !reflect.DeepEqual(got, tt.want.groups) {
				t.Errorf("Grouping.Groups() = %v, want %v", got, tt.want.groups)
			}
			gotSorted := g.SortedGroups()
			if !reflect.DeepEqual(gotSorted, tt.want.sorted) {
				t.Errorf("Grouping.SortedGroups() = %v, want %v", gotSorted, tt.want.sorted)
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
		want *DataFrame
	}{
		{name: "pass", args: args{"1"}, want: MustNew([]interface{}{[]int{1, 2}}, Config{Index: []int{1, 1}})},
		{name: "fail", args: args{"100"}, want: newEmptyDataFrame()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			df := MustNew([]interface{}{[]int{1, 2, 3, 4}}, Config{Index: []int{1, 1, 2, 2}})
			g := df.GroupByIndex()
			got := g.Group(tt.args.label)
			if !Equal(got, tt.want) {
				t.Errorf("Grouping.Group() = %v, want %v", got, tt.want)
			}
			if strings.Contains(tt.name, "fail") {
				if buf.String() == "" {
					t.Errorf("Grouping.Group() returned no log message, want log due to fail")
				}
			}

		})
	}
}

func TestGrouping_Nth(t *testing.T) {
	s := MustNew([]interface{}{[]string{"foo", "bar", "baz"}}, Config{MultiIndex: []interface{}{[]int{1, 1, 2}, []int{2, 2, 1}}})
	g := s.GroupByIndex()
	gotFirst := g.First()
	wantFirst := MustNew([]interface{}{[]string{"foo", "baz"}}, Config{MultiIndex: []interface{}{[]int{1, 2}, []int{2, 1}}})
	if !Equal(gotFirst, wantFirst) {
		t.Errorf("Grouping.First() = %v, want %v", gotFirst, wantFirst)
	}
	gotLast := g.Last()
	wantLast := MustNew([]interface{}{[]string{"bar", "baz"}}, Config{MultiIndex: []interface{}{[]int{1, 2}, []int{2, 1}}})
	if !Equal(gotLast, wantLast) {
		t.Errorf("Grouping.Last() = %v, want %v", gotLast, wantLast)
	}
}

func TestGrouping_Math(t *testing.T) {
	df := MustNew([]interface{}{[]int{1, 2, 3, 4}},
		Config{Col: []string{"A"}, Index: []int{1, 1, 2, 2}})
	tests := []struct {
		name  string
		input *DataFrame
		fn    func(Grouping) *DataFrame
		want  *DataFrame
	}{
		{name: "fail: empty", input: newEmptyDataFrame(), fn: Grouping.Sum,
			want: newEmptyDataFrame()},
		{"sum", df, Grouping.Sum,
			MustNew([]interface{}{[]float64{3, 7}}, Config{Index: []int{1, 2}, Col: []string{"A"}})},
		{"mean", df, Grouping.Mean,
			MustNew([]interface{}{[]float64{1.5, 3.5}}, Config{Index: []int{1, 2}, Col: []string{"A"}})},
		{"min", df, Grouping.Min,
			MustNew([]interface{}{[]float64{1, 3}}, Config{Index: []int{1, 2}, Col: []string{"A"}})},
		{"max", df, Grouping.Max,
			MustNew([]interface{}{[]float64{2, 4}}, Config{Index: []int{1, 2}, Col: []string{"A"}})},
		{"median", df, Grouping.Median,
			MustNew([]interface{}{[]float64{1.5, 3.5}}, Config{Index: []int{1, 2}, Col: []string{"A"}})},
		{"standard deviation", df, Grouping.Std,
			MustNew([]interface{}{[]float64{0.5, 0.5}}, Config{Index: []int{1, 2}, Col: []string{"A"}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.input.GroupByIndex()
			// Test Asynchronously
			got := tt.fn(g)
			if !Equal(got, tt.want) {
				t.Errorf("df.GroupByIndex math operation returned %v, want %v", got, tt.want)
			}
			// Test Synchronously
			options.SetAsync(false)
			gotSync := tt.fn(g)
			if !Equal(gotSync, tt.want) {
				t.Errorf("df.GroupByIndex synchronous math operation returned %v, want %v", gotSync, tt.want)
			}
			options.RestoreDefaults()
		})
	}
}
