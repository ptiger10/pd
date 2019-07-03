package dataframe

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/ptiger10/pd/internal/index"
)

func TestDataFrame_GroupByIndex(t *testing.T) {
	lvl1 := index.MustNewLevel(1, "")
	lvl2 := index.MustNewLevel(2, "")
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

			df := tt.input.Copy()
			got := df.GroupByIndex(tt.args.levelPositions...).groups
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DataFrame.GroupByIndex() = %#v, want %#v", got, tt.want)
			}

			if strings.Contains(tt.name, "fail") {
				if buf.String() == "" {
					t.Errorf("DataFrame.GroupByIndex() returned no log message, want log due to fail")
				}
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
