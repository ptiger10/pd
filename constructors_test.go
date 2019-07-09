package pd

import (
	"fmt"
	"testing"

	"github.com/d4l3k/messagediff"
	"github.com/ptiger10/pd/dataframe"
	"github.com/ptiger10/pd/series"
)

func TestSeries(t *testing.T) {
	type args struct {
		data   interface{}
		config []Config
	}
	type want struct {
		series *series.Series
		err    bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"no config", args{"foo", nil}, want{series.MustNew("foo"), false}},
		{"config",
			args{"foo", []Config{Config{Name: "bar"}}},
			want{series.MustNew("foo", series.Config{Name: "bar"}), false}},
		{"config with df field",
			args{"foo", []Config{Config{Name: "bar", Col: []string{"baz"}}}},
			want{series.MustNew("foo", series.Config{Name: "bar"}), false}},
		{"fail: multiple configs",
			args{"foo", []Config{Config{Name: "bar"}, Config{Name: "baz"}}},
			want{series.MustNew(nil), true}},
		{"fail: unsupported value",
			args{complex64(1), nil},
			want{series.MustNew(nil), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Series(tt.args.data, tt.args.config...)
			if (err != nil) != tt.want.err {
				t.Errorf("Series():  error = %v, want %v", err, tt.want.err)
			}
			if !series.Equal(got, tt.want.series) {
				t.Errorf("Series() got %v, want %v", got, tt.want.series)
			}
		})
	}
}

func TestDataFrame(t *testing.T) {
	type args struct {
		data   []interface{}
		config []Config
	}
	type want struct {
		df  *dataframe.DataFrame
		err bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"no config", args{[]interface{}{"foo"}, nil}, want{dataframe.MustNew([]interface{}{"foo"}), false}},
		{"config",
			args{[]interface{}{"foo"}, []Config{Config{Name: "bar"}}},
			want{dataframe.MustNew([]interface{}{"foo"}, dataframe.Config{Name: "bar"}), false}},
		{"fail: multiple configs",
			args{[]interface{}{"foo"}, []Config{Config{Name: "bar"}, Config{Name: "baz"}}},
			want{dataframe.MustNew(nil), true}},
		{"fail: unsupported value",
			args{[]interface{}{complex64(1)}, nil},
			want{dataframe.MustNew(nil), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DataFrame(tt.args.data, tt.args.config...)
			if (err != nil) != tt.want.err {
				t.Errorf("DataFrame():  error = %v, want %v", err, tt.want.err)
			}
			if !dataframe.Equal(got, tt.want.df) {
				t.Errorf("DataFrame() got %v, want %v", got, tt.want.df)
			}
		})
	}
}

func TestReadCSV(t *testing.T) {
	type args struct {
		filepath string
		options  []ReadOptions
	}
	type want struct {
		df  *dataframe.DataFrame
		err bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{name: "no interpolation", args: args{filepath: "csv_test.csv", options: []ReadOptions{ReadOptions{Manual: true}}},
			want: want{
				df: dataframe.MustNew([]interface{}{
					[]string{"", "foo", "bar"},
					[]string{"A", "1", "2"},
				}),
				err: false}},
		{"fail: bad path", args{"foo.csv", nil}, want{dataframe.MustNew(nil), true}},
		{"interpolation", args{"csv_test.csv", []ReadOptions{ReadOptions{IndexCols: 1, HeaderRows: 1}}},
			want{
				df: dataframe.MustNew([]interface{}{
					[]int64{1, 2},
				}, dataframe.Config{Index: []string{"foo", "bar"}, Col: []string{"A"}}),
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadCSV(tt.args.filepath, tt.args.options...)
			if (err != nil) != tt.want.err {
				t.Errorf("ReadCSV():  error = %v, want %v", err, tt.want.err)
			}
			if !dataframe.Equal(got, tt.want.df) {
				t.Errorf("ReadCSV() got \n%v, \nwant \n%v", got, tt.want.df)
				diff, _ := messagediff.PrettyDiff(got, tt.want.df)
				fmt.Println(diff)
			}
		})
	}
}

func TestInterface(t *testing.T) {
	data := [][]interface{}{{"foo", "bar"}, {"baz", "qux"}}
	type args struct {
		data    [][]interface{}
		options []ReadOptions
	}
	type want struct {
		df  *dataframe.DataFrame
		err bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{name: "no config", args: args{data: data, options: nil},
			want: want{
				df: dataframe.MustNew([]interface{}{
					[]string{"foo", "baz"},
					[]string{"bar", "qux"},
				}),
				err: false}},
		{"drop 1 row", args{data, []ReadOptions{ReadOptions{DropRows: 1}}},
			want{
				dataframe.MustNew([]interface{}{
					[]string{"baz"},
					[]string{"qux"},
				}),
				false}},
		{"1 header row", args{data, []ReadOptions{ReadOptions{HeaderRows: 1}}},
			want{
				dataframe.MustNew([]interface{}{
					[]string{"baz"},
					[]string{"qux"},
				}, dataframe.Config{Col: []string{"foo", "bar"}}),
				false}},
		{"1 index column", args{data, []ReadOptions{ReadOptions{IndexCols: 1}}},
			want{
				dataframe.MustNew([]interface{}{
					[]string{"bar", "qux"},
				}, dataframe.Config{Index: []string{"foo", "baz"}}),
				false}},
		{"1 header row, 1 index column", args{data, []ReadOptions{ReadOptions{IndexCols: 1, HeaderRows: 1}}},
			want{
				dataframe.MustNew([]interface{}{
					[]string{"qux"},
				}, dataframe.Config{Index: []string{"baz"}, Col: []string{"bar"}}),
				false}},
		{"1 header row, 1 index column, datatype conversion", args{data,
			[]ReadOptions{ReadOptions{
				IndexCols:  1,
				HeaderRows: 1,
				DataTypes:  map[string]string{"bar": "bool"},
			}}},
			want{
				dataframe.MustNew([]interface{}{
					[]bool{true},
				}, dataframe.Config{Index: []string{"baz"}, Col: []string{"bar"}}),
				false}},
		{"1 header row, 1 index column, rename column", args{data,
			[]ReadOptions{ReadOptions{
				IndexCols:  1,
				HeaderRows: 1,
				Rename:     map[string]string{"bar": "corge"}}}},
			want{
				dataframe.MustNew([]interface{}{
					[]string{"qux"},
				}, dataframe.Config{Index: []string{"baz"}, Col: []string{"corge"}}),
				false}},
		{"1 header row, 1 index column, convert index type", args{data,
			[]ReadOptions{ReadOptions{
				IndexCols:      1,
				HeaderRows:     1,
				IndexDataTypes: map[int]string{0: "bool"}}}},
			want{
				dataframe.MustNew([]interface{}{
					[]string{"qux"},
				}, dataframe.Config{Index: []bool{true}, Col: []string{"bar"}}),
				false}},
		{"fail: too many headers", args{data,
			[]ReadOptions{ReadOptions{
				HeaderRows: 10,
			}}}, want{dataframe.MustNew(nil), true}},
		{"fail: too many index columns", args{data,
			[]ReadOptions{ReadOptions{
				IndexCols: 10,
			}}}, want{dataframe.MustNew(nil), true}},
		{"fail: drop too many rows", args{data,
			[]ReadOptions{ReadOptions{
				DropRows: 10,
			}}}, want{dataframe.MustNew(nil), true}},
		{"fail: excessive ReadOptions", args{data,
			[]ReadOptions{ReadOptions{}, ReadOptions{}}}, want{dataframe.MustNew(nil), true}},
		{"fail: no rows", args{[][]interface{}{}, nil}, want{dataframe.MustNew(nil), true}},
		{"fail: no columns", args{[][]interface{}{[]interface{}{}}, nil}, want{dataframe.MustNew(nil), true}},
		// TODO: []interface{}{complex64} should fail
		// {"fail: unsupported data", args{[][]interface{}{{complex64(1)}}, nil}, want{dataframe.MustNew(nil), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadInterface(tt.args.data, tt.args.options...)
			if (err != nil) != tt.want.err {
				t.Errorf("ReadInterface():  error = %v, want %v", err, tt.want.err)
			}
			if !dataframe.Equal(got, tt.want.df) {
				t.Errorf("ReadInterface() got \n%v, \nwant \n%v", got, tt.want.df)
				diff, _ := messagediff.PrettyDiff(got, tt.want.df)
				fmt.Println(diff)
			}
		})
	}
}
