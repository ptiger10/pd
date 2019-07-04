package dataframe

import (
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
	df := MustNew([]interface{}{[]string{"foo"}, []string{"bar"}},
		Config{Col: []string{"baz", "qux"}})
	got := df.Col("qux")
	want := series.MustNew([]string{"bar"}, series.Config{Name: "qux"})
	if !series.Equal(got, want) {
		t.Errorf("Col(): got %v, want %v", got, want)
	}
}
