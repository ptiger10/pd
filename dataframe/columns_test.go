package dataframe

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestColumns_Describe(t *testing.T) {
	singleDefault := MustNew([]interface{}{"foo", "bar", "baz"})
	multiConfig := MustNew([]interface{}{"foo", "bar", "baz"},
		Config{MultiCol: [][]string{{"1", "2", "3"}, {"qux", "quux", "quuz"}}})
	type args struct {
		atLevel int
		atCol   int
	}
	type want struct {
		at     string
		values [][]string
	}
	tests := []struct {
		name  string
		input *DataFrame
		args  args
		want  want
	}{
		{name: "single default", input: singleDefault, args: args{atLevel: 0, atCol: 2},
			want: want{at: "2", values: [][]string{{"0", "1", "2"}}}},
		{"multi from config", multiConfig, args{1, 2}, want{"quuz", [][]string{{"1", "2", "3"}, {"qux", "quux", "quuz"}}}},
		{"soft fail: at invalid level", singleDefault, args{10, 0}, want{"", [][]string{{"0", "1", "2"}}}},
		{"soft fail: at invalid col", singleDefault, args{0, 10}, want{"", [][]string{{"0", "1", "2"}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			cols := Columns{
				df: tt.input,
			}

			gotAt := cols.At(tt.args.atLevel, tt.args.atCol)
			if gotAt != tt.want.at {
				t.Errorf("Columns.At(): got %v, want %v", gotAt, tt.want.at)
			}
			gotValues := cols.Values()
			if !reflect.DeepEqual(gotValues, tt.want.values) {
				t.Errorf("Columns.Values(): got %v, want %v", gotValues, tt.want.values)
			}

			if strings.Contains(tt.name, "fail:") {
				if buf.String() == "" {
					t.Errorf("Columns operation returned no log message, want log due to fail")
				}
			}
		})
	}
}

func TestColumns_RenameLevel(t *testing.T) {
	df := MustNew([]interface{}{"foo"}, Config{MultiCol: [][]string{{"bar"}, {"baz"}}, MultiColNames: []string{"qux", "quuz"}})
	type args struct {
		level int
		name  string
	}
	tests := []struct {
		name    string
		input   *DataFrame
		args    args
		want    *DataFrame
		wantErr bool
	}{
		{name: "pass 0", input: df, args: args{0, "corge"},
			want:    MustNew([]interface{}{"foo"}, Config{MultiCol: [][]string{{"bar"}, {"baz"}}, MultiColNames: []string{"corge", "quuz"}}),
			wantErr: false},
		{"fail", df, args{10, "corge"},
			df, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cols := Columns{
				df: tt.input.Copy(),
			}
			if err := cols.RenameLevel(tt.args.level, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Columns.RenameLevel() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !Equal(cols.df, tt.want) {
				t.Errorf("Columns.RenameLevel(): got %v, want %v", cols.df, tt.want)
			}
		})
	}
}

func TestColumns_SwapLevels(t *testing.T) {
	df := MustNew([]interface{}{[]string{"foo", "bar"}},
		Config{MultiCol: [][]string{{"baz"}, {"qux"}}, MultiColNames: []string{"quux", "quuz"}})
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name    string
		input   *DataFrame
		args    args
		want    *DataFrame
		wantErr bool
	}{
		{name: "pass", input: df, args: args{0, 1},
			want: MustNew([]interface{}{[]string{"foo", "bar"}},
				Config{MultiCol: [][]string{{"qux"}, {"baz"}}, MultiColNames: []string{"quuz", "quux"}}),
			wantErr: false},
		{"fail: i", df, args{10, 1},
			df, true},
		{"fail: j", df, args{0, 10},
			df, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cols := Columns{
				df: tt.input.Copy(),
			}
			err := cols.SwapLevels(tt.args.i, tt.args.j)
			if (err != nil) != tt.wantErr {
				t.Errorf("Column.SwapLevels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(cols.df, tt.want) {
				t.Errorf("Column.SwapLevels() = %v, want %v", cols.df, tt.want)
			}
		})
	}
}

func TestColumns_InsertLevel(t *testing.T) {
	df := MustNew([]interface{}{"foo"}, Config{
		MultiCol: [][]string{{"bar"}, {"baz"}}, MultiColNames: []string{"quux", "quuz"}})
	type args struct {
		pos    int
		labels []string
		name   string
	}
	tests := []struct {
		name    string
		input   *DataFrame
		args    args
		want    *DataFrame
		wantErr bool
	}{
		{name: "0", input: df, args: args{pos: 0, labels: []string{"qux"}, name: "corge"},
			want: MustNew([]interface{}{"foo"},
				Config{MultiCol: [][]string{{"qux"}, {"bar"}, {"baz"}}, MultiColNames: []string{"corge", "quux", "quuz"}}),
			wantErr: false},
		{"1", df, args{1, []string{"qux"}, "corge"},
			MustNew([]interface{}{"foo"},
				Config{MultiCol: [][]string{{"bar"}, {"qux"}, {"baz"}}, MultiColNames: []string{"quux", "corge", "quuz"}}),
			false},
		{"2", df, args{2, []string{"qux"}, "corge"},
			MustNew([]interface{}{"foo"},
				Config{MultiCol: [][]string{{"bar"}, {"baz"}, {"qux"}}, MultiColNames: []string{"quux", "quuz", "corge"}}),
			false},
		{"fail: invalid position", df, args{10, []string{"bar"}, "corge"},
			df, true},
		{"fail: excessive col labels", df, args{1, []string{"bar", "waldo", "fred"}, "corge"},
			df, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cols := Columns{
				df: tt.input.Copy(),
			}
			err := cols.InsertLevel(tt.args.pos, tt.args.labels, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Column.InsertLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(cols.df, tt.want) {
				t.Errorf("Column.InsertLevel() = %v, want %v", cols.df, tt.want)
			}
		})
	}
}

func TestColumns_AppendLevel(t *testing.T) {
	df := MustNew([]interface{}{"foo"}, Config{
		MultiCol: [][]string{{"bar"}, {"baz"}}, MultiColNames: []string{"quux", "quuz"}})
	type args struct {
		labels []string
		name   string
	}
	tests := []struct {
		name    string
		input   *DataFrame
		args    args
		want    *DataFrame
		wantErr bool
	}{
		{name: "pass", input: df, args: args{labels: []string{"qux"}, name: "corge"},
			want: MustNew([]interface{}{"foo"},
				Config{MultiCol: [][]string{{"bar"}, {"baz"}, {"qux"}}, MultiColNames: []string{"quux", "quuz", "corge"}}),
			wantErr: false},
		{"fail: misaligned length", df, args{[]string{"bar", "waldo", "fred"}, ""},
			df, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cols := Columns{
				df: tt.input.Copy(),
			}
			err := cols.AppendLevel(tt.args.labels, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Column.AppendLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(cols.df, tt.want) {
				t.Errorf("Column.AppendLevel() = %v, want %v", cols.df, tt.want)
			}
		})
	}
}

func TestColumns_DropLevel(t *testing.T) {
	df := MustNew([]interface{}{"foo"}, Config{
		MultiCol: [][]string{{"bar"}, {"baz"}}, MultiColNames: []string{"quux", "quuz"}})
	type args struct {
		level int
	}
	tests := []struct {
		name    string
		input   *DataFrame
		args    args
		want    *DataFrame
		wantErr bool
	}{
		{name: "pass", input: df, args: args{0},
			want: MustNew([]interface{}{"foo"}, Config{
				MultiCol: [][]string{{"baz"}}, MultiColNames: []string{"quuz"}}),
			wantErr: false},
		{name: "replace with default", input: MustNew([]interface{}{"foo"}, Config{Col: []string{"quz"}}), args: args{0},
			want:    MustNew([]interface{}{"foo"}),
			wantErr: false},
		{"fail: invalid level", df, args{10},
			df, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cols := Columns{
				df: tt.input.Copy(),
			}
			err := cols.DropLevel(tt.args.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("Column.DropLevel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(cols.df, tt.want) {
				t.Errorf("Column.DropLevel() = %v, want %v", cols.df, tt.want)
			}
		})
	}
}

func TestColumns_SelectName(t *testing.T) {
	df := MustNew([]interface{}{"foo"}, Config{
		MultiCol: [][]string{{"bar"}, {"baz"}}, MultiColNames: []string{"quux", "quuz"}})
	type args struct {
		name string
	}
	tests := []struct {
		name  string
		input *DataFrame
		args  args
		want  int
	}{
		{name: "pass", input: df, args: args{"quux"}, want: 0},
		{name: "soft fail: invalid name", input: df, args: args{"fred"}, want: -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cols := Columns{
				df: tt.input,
			}
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)
			if got := cols.SelectName(tt.args.name); got != tt.want {
				t.Errorf("Column.SelectName() = %v, want %v", got, tt.want)
			}
			if strings.Contains(tt.name, "fail:") {
				if buf.String() == "" {
					t.Errorf("Column.SelectName() returned no log message, want log due to fail")
				}
			}
		})
	}
}

func TestColumns_SelectNames(t *testing.T) {
	df := MustNew([]interface{}{"foo"}, Config{
		MultiCol: [][]string{{"bar"}, {"baz"}, {"qux"}}, MultiColNames: []string{"quux", "quuz", "quux"}})
	type args struct {
		names []string
	}
	tests := []struct {
		name  string
		input *DataFrame
		args  args
		want  []int
	}{
		{name: "pass", input: df, args: args{[]string{"quux"}}, want: []int{0, 2}},
		{name: "soft fail: invalid name", input: df, args: args{[]string{"fred"}}, want: []int{}},
		{name: "soft fail: partial invalid name", input: df, args: args{[]string{"quux", "fred"}}, want: []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cols := Columns{
				df: tt.input.Copy(),
			}
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)
			if got := cols.SelectNames(tt.args.names); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Column.SelectNames(): got %v, want %v", got, tt.want)
			}
			if strings.Contains(tt.name, "fail:") {
				if buf.String() == "" {
					t.Errorf("Column.SelectNames() returned no log message, want log due to fail")
				}
			}
		})
	}
}

func TestColumns_SubsetLevels(t *testing.T) {
	df := MustNew([]interface{}{"foo"}, Config{
		MultiCol: [][]string{{"bar"}, {"baz"}, {"qux"}}})
	type args struct {
		levelPositions []int
	}
	tests := []struct {
		name    string
		input   *DataFrame
		args    args
		want    *DataFrame
		wantErr bool
	}{
		{name: "one level", input: df, args: args{[]int{0}},
			want: MustNew([]interface{}{"foo"}, Config{MultiCol: [][]string{{"bar"}}}), wantErr: false},
		{name: "multiple levels", input: df, args: args{[]int{0, 1}},
			want: MustNew([]interface{}{"foo"}, Config{MultiCol: [][]string{{"bar"}, {"baz"}}}), wantErr: false},
		{"fail: invalid level", df, args{[]int{10}},
			df, true},
		{"fail: no levels", df, args{[]int{}},
			df, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cols := Columns{
				df: tt.input.Copy(),
			}
			err := cols.SubsetLevels(tt.args.levelPositions)
			if (err != nil) != tt.wantErr {
				t.Errorf("Column.Subset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(cols.df, tt.want) {
				t.Errorf("Column.Subset() = %v, want %v", cols.df, tt.want)
			}
		})
	}
}
