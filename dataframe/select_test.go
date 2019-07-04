package dataframe

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/ptiger10/pd/series"
)

func TestSubsetRows(t *testing.T) {
	df := MustNew([]interface{}{[]string{"foo", "bar", "baz"}}, Config{Index: []int{0, 1, 2}})
	tests := []struct {
		name    string
		input   *DataFrame
		args    []int
		want    *DataFrame
		wantErr bool
	}{
		{name: "subset one row", input: df, args: []int{0}, want: MustNew([]interface{}{"foo"}, Config{Index: []int{0}}), wantErr: false},
		{"subset two rows", df, []int{0, 1}, MustNew([]interface{}{[]string{"foo", "bar"}}, Config{Index: []int{0, 1}}), false},
		{"two rows reverse", df, []int{1, 0}, MustNew([]interface{}{[]string{"bar", "foo"}}, Config{Index: []int{1, 0}}), false},
		{"fail: empty args", df, []int{}, newEmptyDataFrame(), true},
		{"fail: invalid row", df, []int{3}, newEmptyDataFrame(), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := tt.input.Copy()
			dfArchive := tt.input.Copy()
			err := df.InPlace.SubsetRows(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("InPlace.SubsetRows() error = %v, want %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(df, tt.want) {
					t.Errorf("InPlace.SubsetRows() got %v, want %v", df, tt.want)
				}
			}

			dfCopy, err := dfArchive.SubsetRows(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("DataFrame.SubsetRows() error = %v, want %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(dfCopy, tt.want) {
					t.Errorf("DataFrame.SubsetRows() got %v, want %v", dfCopy, tt.want)
				}
				if Equal(dfArchive, dfCopy) {
					t.Errorf("DataFrame.SubsetRows() retained access to original, want copy")
				}
			}
		})
	}
}

func TestSubsetColumns(t *testing.T) {
	df := MustNew([]interface{}{"foo", "bar", "baz"}, Config{Col: []string{"0", "1", "2"}})
	tests := []struct {
		name    string
		input   *DataFrame
		args    []int
		want    *DataFrame
		wantErr bool
	}{
		{name: "subset one column", input: df, args: []int{0}, want: MustNew([]interface{}{"foo"}, Config{Col: []string{"0"}}), wantErr: false},
		{"subset two columns", df, []int{0, 1}, MustNew([]interface{}{"foo", "bar"}, Config{Col: []string{"0", "1"}}), false},
		{"two columns reverse", df, []int{1, 0}, MustNew([]interface{}{"bar", "foo"}, Config{Col: []string{"1", "0"}}), false},
		{"fail: empty args", df, []int{}, newEmptyDataFrame(), true},
		{"fail: invalid column", df, []int{10}, newEmptyDataFrame(), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			df := tt.input.Copy()
			dfArchive := tt.input.Copy()
			err := df.InPlace.SubsetColumns(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("InPlace.SubsetColumns() error = %v, want %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(df, tt.want) {
					t.Errorf("InPlace.SubsetColumns() got %v, want %v", df, tt.want)
				}
			}

			dfCopy, err := dfArchive.SubsetColumns(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("DataFrame.SubsetColumns() error = %v, want %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(tt.name, "fail") {
				if !Equal(dfCopy, tt.want) {
					t.Errorf("DataFrame.SubsetColumns() got %v, want %v", dfCopy, tt.want)
				}
				if Equal(dfArchive, dfCopy) {
					t.Errorf("DataFrame.SubsetColumns() retained access to original, want copy")
				}
			}
		})
	}
}

func TestCol(t *testing.T) {
	df := MustNew([]interface{}{[]string{"foo"}, []string{"bar"}}, Config{Col: []string{"baz", "qux"}})
	type args struct {
		label string
	}
	tests := []struct {
		name  string
		input *DataFrame
		args  args
		want  *series.Series
	}{
		{name: "pass", input: df, args: args{label: "baz"},
			want: series.MustNew([]string{"foo"}, series.Config{Name: "baz"})},
		{"fail: invalid col", df, args{label: "corge"}, series.MustNew(nil)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			if got := tt.input.Col(tt.args.label); !series.Equal(got, tt.want) {
				t.Errorf("DataFrame.Col() = %v, want %v", got, tt.want)
			}
			if strings.Contains(tt.name, "fail") {
				if buf.String() == "" {
					t.Errorf("DataFrame.Col() returned no log message, want log due to fail")
				}
			}
		})
	}
}

func TestSeries_SelectLabel(t *testing.T) {
	df := MustNew([]interface{}{[]string{"foo", "bar", "baz"}}, Config{Index: []int{1, 1, 1}})
	type args struct {
		label string
	}
	tests := []struct {
		name  string
		input *DataFrame
		args  args
		want  int
	}{
		{name: "pass", input: df, args: args{label: "1"}, want: 0},
		{"fail: empty DataFrame", newEmptyDataFrame(), args{label: "1"}, -1},
		{"fail: label not in DataFrame", df, args{label: "100"}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			if got := tt.input.SelectLabel(tt.args.label); got != tt.want {
				t.Errorf("DataFrame.SelectLabel() = %v, want %v", got, tt.want)
			}
			if strings.Contains(tt.name, "fail") {
				if buf.String() == "" {
					t.Errorf("DataFrame.SelectLabel() returned no log message, want log due to fail")
				}
			}
		})
	}
}

func TestSeries_SelectLabels(t *testing.T) {
	df := MustNew([]interface{}{[]string{"foo", "bar", "baz"}}, Config{MultiIndex: []interface{}{[]int{1, 1, 1}, []string{"qux", "quux", "quuz"}}})
	type args struct {
		labels []string
		level  int
	}
	tests := []struct {
		name  string
		input *DataFrame
		args  args
		want  []int
	}{
		{name: "pass", input: df, args: args{labels: []string{"1"}, level: 0}, want: []int{0, 1, 2}},
		{"pass", df, args{[]string{"qux", "quux"}, 1}, []int{0, 1}},
		{"fail: empty DataFrame", newEmptyDataFrame(), args{[]string{"1"}, 0}, []int{}},
		{"fail: label not in DataFrame", df, args{[]string{"100"}, 0}, []int{}},
		{"fail: invalid level", df, args{[]string{"1"}, 100}, []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			if got := tt.input.SelectLabels(tt.args.labels, tt.args.level); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Series.SelectLabels() = %v, want %v", got, tt.want)
			}
			if strings.Contains(tt.name, "fail") {
				if buf.String() == "" {
					t.Errorf("Series.SelectLabels() returned no log message, want log due to fail")
				}
			}
		})
	}
}

func TestSeries_SelectColumn(t *testing.T) {
	df := MustNew([]interface{}{"foo", "bar", "baz"}, Config{Col: []string{"qux", "quuz", "qux"}})
	type args struct {
		label string
	}
	tests := []struct {
		name  string
		input *DataFrame
		args  args
		want  int
	}{
		{name: "pass", input: df, args: args{label: "qux"}, want: 0},
		{"fail: empty DataFrame", newEmptyDataFrame(), args{label: "1"}, -1},
		{"fail: label not in DataFrame", df, args{label: "100"}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			if got := tt.input.SelectColumn(tt.args.label); got != tt.want {
				t.Errorf("DataFrame.SelectColumn() = %v, want %v", got, tt.want)
			}
			if strings.Contains(tt.name, "fail") {
				if buf.String() == "" {
					t.Errorf("DataFrame.SelectColumn() returned no log message, want log due to fail")
				}
			}
		})
	}
}

func TestSeries_SelectColumns(t *testing.T) {
	df := MustNew([]interface{}{"foo", "bar", "baz"}, Config{MultiCol: [][]string{{"qux", "quux", "qux"}, {"corge", "waldo", "fred"}}})
	type args struct {
		labels []string
		level  int
	}
	tests := []struct {
		name  string
		input *DataFrame
		args  args
		want  []int
	}{
		{name: "pass", input: df, args: args{labels: []string{"qux"}, level: 0}, want: []int{0, 2}},
		{"pass", df, args{[]string{"waldo", "corge"}, 1}, []int{1, 0}},
		{"fail: empty DataFrame", newEmptyDataFrame(), args{[]string{"1"}, 0}, []int{}},
		{"fail: label not in DataFrame", df, args{[]string{"100"}, 0}, []int{}},
		{"fail: invalid level", df, args{[]string{"1"}, 100}, []int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			if got := tt.input.SelectColumns(tt.args.labels, tt.args.level); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Series.SelectColumns() = %v, want %v", got, tt.want)
			}
			if strings.Contains(tt.name, "fail") {
				if buf.String() == "" {
					t.Errorf("Series.SelectColumns() returned no log message, want log due to fail")
				}
			}
		})
	}
}
