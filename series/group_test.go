package series

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/ptiger10/pd/options"

	"github.com/ptiger10/pd/internal/index"
)

func TestGrouping_Math(t *testing.T) {
	s := MustNew([]int{1, 2, 3, 4}, Config{Index: []int{1, 1, 2, 2}})
	tests := []struct {
		name  string
		input *Series
		fn    func(Grouping) *Series
		want  *Series
	}{
		{name: "fail: empty", input: newEmptySeries(), fn: Grouping.Sum,
			want: newEmptySeries()},
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
			// Test Asynchronously
			got := tt.fn(g)
			if !Equal(got, tt.want) {
				t.Errorf("s.GroupByIndex math operation returned %v, want %v", got, tt.want)
			}
			// Test Synchronously
			options.SetAsync(false)
			gotSync := tt.fn(g)
			if !Equal(gotSync, tt.want) {
				t.Errorf("s.GroupByIndex synchronous math operation returned %v, want %v", gotSync, tt.want)
			}
			options.RestoreDefaults()
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

func TestSeries_GroupByIndex(t *testing.T) {
	lvl1 := index.MustNewLevel(1, "")
	lvl2 := index.MustNewLevel(2, "")
	multi := MustNew([]string{"foo", "bar", "baz"}, Config{MultiIndex: []interface{}{[]int{1, 1, 2}, []int{2, 2, 1}}})
	type args struct {
		levelPositions []int
	}
	tests := []struct {
		name  string
		input *Series
		args  args
		want  map[string]*group
	}{
		{name: "single no args",
			input: MustNew([]string{"foo", "bar", "baz"}, Config{Index: []int{1, 1, 2}}),
			args:  args{[]int{}},
			want: map[string]*group{
				"1": &group{Positions: []int{0, 1}, Index: index.New(lvl1)},
				"2": &group{Positions: []int{2}, Index: index.New(lvl2)},
			}},
		{"multi no args",
			multi,
			args{[]int{}},
			map[string]*group{
				"1 2": &group{Positions: []int{0, 1}, Index: index.New(lvl1, lvl2)},
				"2 1": &group{Positions: []int{2}, Index: index.New(lvl2, lvl1)},
			}},
		{"multi one level",
			multi,
			args{[]int{0}},
			map[string]*group{
				"1": &group{Positions: []int{0, 1}, Index: index.New(lvl1)},
				"2": &group{Positions: []int{2}, Index: index.New(lvl2)},
			}},
		{"multi two levels reversed",
			multi,
			args{[]int{1, 0}},
			map[string]*group{
				"2 1": &group{Positions: []int{0, 1}, Index: index.New(lvl2, lvl1)},
				"1 2": &group{Positions: []int{2}, Index: index.New(lvl1, lvl2)},
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

			s := tt.input.Copy()
			got := s.GroupByIndex(tt.args.levelPositions...).groups
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Series.GroupByIndex() = %#v, want %#v", got, tt.want)
			}

			if strings.Contains(tt.name, "fail") {
				if buf.String() == "" {
					t.Errorf("Series.GroupByIndex() returned no log message, want log due to fail")
				}
			}
		})
	}
}

func TestGrouping_Nth(t *testing.T) {
	s := MustNew([]string{"foo", "bar", "baz"}, Config{MultiIndex: []interface{}{[]int{1, 1, 2}, []int{2, 2, 1}}})
	g := s.GroupByIndex()
	gotFirst := g.First()
	wantFirst := MustNew([]string{"foo", "baz"}, Config{MultiIndex: []interface{}{[]int{1, 2}, []int{2, 1}}})
	if !Equal(gotFirst, wantFirst) {
		t.Errorf("Grouping.First() = %#v, want %#v", gotFirst, wantFirst)
	}
	gotLast := g.Last()
	wantLast := MustNew([]string{"bar", "baz"}, Config{MultiIndex: []interface{}{[]int{1, 2}, []int{2, 1}}})
	if !Equal(gotLast, wantLast) {
		t.Errorf("Grouping.Last() = %#v, want %#v", gotLast, wantLast)
	}
}
