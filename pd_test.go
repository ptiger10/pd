package pd

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

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
			args{"foo", []Config{{Name: "bar"}}},
			want{series.MustNew("foo", series.Config{Name: "bar"}), false}},
		{"config with df field",
			args{"foo", []Config{{Name: "bar", Col: []string{"baz"}}}},
			want{series.MustNew("foo", series.Config{Name: "bar"}), false}},
		{"fail: multiple configs",
			args{"foo", []Config{{Name: "bar"}, {Name: "baz"}}},
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
			args{[]interface{}{"foo"}, []Config{{Name: "bar"}}},
			want{dataframe.MustNew([]interface{}{"foo"}, dataframe.Config{Name: "bar"}), false}},
		{"fail: multiple configs",
			args{[]interface{}{"foo"}, []Config{{Name: "bar"}, {Name: "baz"}}},
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
		{name: "no interpolation", args: args{filepath: "csv_tests/pass.csv", options: []ReadOptions{{Manual: true}}},
			want: want{
				df: dataframe.MustNew([]interface{}{
					[]string{"", "foo", "bar"},
					[]string{"A", "1", "2"},
				}),
				err: false}},
		{name: "pipe delimiter", args: args{filepath: "csv_tests/pipe.csv",
			options: []ReadOptions{{Delimiter: '|', HeaderRows: 1, IndexCols: 1}}},
			want: want{
				df: dataframe.MustNew([]interface{}{1, 2},
					dataframe.Config{Index: "foo", Col: []string{"A", "B"}}),
				err: false}},
		{"interpolation", args{"csv_tests/pass.csv", []ReadOptions{{IndexCols: 1, HeaderRows: 1}}},
			want{
				dataframe.MustNew([]interface{}{
					[]int64{1, 2},
				}, dataframe.Config{Index: []string{"foo", "bar"}, Col: []string{"A"}}),
				false}},
		{"fail: bad path", args{"foo.csv", nil}, want{dataframe.MustNew(nil), true}},
		{"fail: bad delimiter", args{"csv_tests/pass.csv",
			[]ReadOptions{{Delimiter: '\n'}}}, want{dataframe.MustNew(nil), true}},
		{"fail: empty", args{"csv_tests/empty.csv", nil}, want{dataframe.MustNew(nil), true}},
		{"fail: corrupted file", args{"csv_tests/corrupted.csv", nil}, want{dataframe.MustNew(nil), true}},
		{"fail: too many configs", args{"csv_tests/pass.csv", []ReadOptions{{}, {}}}, want{dataframe.MustNew(nil), true}},
		{"fail within ReadInterface", args{"csv_tests/pass.csv", []ReadOptions{{HeaderRows: 10}}},
			want{dataframe.MustNew(nil), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadCSV(tt.args.filepath, tt.args.options...)
			if (err != nil) != tt.want.err {
				t.Errorf("ReadCSV():  error = %v, want %v", err, tt.want.err)
			}
			if !dataframe.Equal(got, tt.want.df) {
				t.Errorf("ReadCSV() got \n%v, \nwant \n%v", got, tt.want.df)
			}
		})
	}
}

func TestInterface(t *testing.T) {
	data := [][]interface{}{{"foo", "bar"}, {"baz", "qux"}}
	noConfig := dataframe.MustNew([]interface{}{
		[]string{"foo", "baz"},
		[]string{"bar", "qux"},
	})
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
				df:  noConfig,
				err: false}},
		{"drop 1 row", args{data, []ReadOptions{{DropRows: 1}}},
			want{
				dataframe.MustNew([]interface{}{
					[]string{"baz"},
					[]string{"qux"},
				}),
				false}},
		{"1 header row", args{data, []ReadOptions{{HeaderRows: 1}}},
			want{
				dataframe.MustNew([]interface{}{
					[]string{"baz"},
					[]string{"qux"},
				}, dataframe.Config{Col: []string{"foo", "bar"}}),
				false}},
		{"1 index column", args{data, []ReadOptions{{IndexCols: 1}}},
			want{
				dataframe.MustNew([]interface{}{
					[]string{"bar", "qux"},
				}, dataframe.Config{Index: []string{"foo", "baz"}}),
				false}},
		{"1 header row, 1 index column", args{data, []ReadOptions{{IndexCols: 1, HeaderRows: 1}}},
			want{
				dataframe.MustNew([]interface{}{
					[]string{"qux"},
				}, dataframe.Config{Index: []string{"baz"}, Col: []string{"bar"}}),
				false}},
		{"1 header row, 1 index column, datatype conversion", args{data,
			[]ReadOptions{{
				IndexCols:  1,
				HeaderRows: 1,
				DataTypes:  map[string]string{"bar": "bool"},
			}}},
			want{dataframe.MustNew([]interface{}{
				[]bool{true},
			}, dataframe.Config{Index: []string{"baz"}, Col: []string{"bar"}}),
				false}},
		{"1 header row, 1 index column, rename column", args{data,
			[]ReadOptions{{
				IndexCols:  1,
				HeaderRows: 1,
				Rename:     map[string]string{"bar": "corge"}}}},
			want{dataframe.MustNew([]interface{}{
				[]string{"qux"},
			}, dataframe.Config{Index: []string{"baz"}, Col: []string{"corge"}}),
				false}},
		{"1 header row, 1 index column, convert index type", args{data,
			[]ReadOptions{{
				IndexCols:      1,
				HeaderRows:     1,
				IndexDataTypes: map[int]string{0: "bool"}}}},
			want{
				dataframe.MustNew([]interface{}{
					[]string{"qux"},
				}, dataframe.Config{Index: []bool{true}, Col: []string{"bar"}}),
				false}},
		{"fail: too many headers", args{data,
			[]ReadOptions{{
				HeaderRows: 10,
			}}}, want{dataframe.MustNew(nil), true}},
		{"fail: too many index columns", args{data,
			[]ReadOptions{{
				IndexCols: 10,
			}}}, want{dataframe.MustNew(nil), true}},
		{"fail: drop too many rows", args{data,
			[]ReadOptions{{
				DropRows: 10,
			}}}, want{dataframe.MustNew(nil), true}},
		{"fail: excessive ReadOptions", args{data,
			[]ReadOptions{{}, {}}}, want{dataframe.MustNew(nil), true}},
		{"fail: no rows", args{[][]interface{}{}, nil}, want{dataframe.MustNew(nil), true}},
		{"fail: no columns", args{[][]interface{}{{}}, nil}, want{dataframe.MustNew(nil), true}},
		{"soft fail: missing column in value conversion", args{data,
			[]ReadOptions{{
				DataTypes: map[string]string{"NOSUCHCOLUMN": "string"},
			}}}, want{noConfig, false}},
		{"soft fail: unsupported value conversion",
			args{data, []ReadOptions{{
				HeaderRows: 1,
				DataTypes:  map[string]string{"NOSUCHCOLUMN": "unsupported"}}}},
			want{dataframe.MustNew([]interface{}{
				[]string{"baz"},
				[]string{"qux"},
			}, dataframe.Config{Col: []string{"foo", "bar"}}),
				false}},
		{"soft fail: unsupported index conversion", args{data,
			[]ReadOptions{{
				IndexDataTypes: map[int]string{0: "unsupported"},
			}}}, want{noConfig, false}},

		// TODO: []interface{}{complex64} should fail
		// {"fail: unsupported data", args{[][]interface{}{{complex64(1)}}, nil}, want{dataframe.MustNew(nil), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			got, err := ReadInterface(tt.args.data, tt.args.options...)
			if (err != nil) != tt.want.err {
				t.Errorf("ReadInterface():  error = %v, want %v", err, tt.want.err)
			}
			if !dataframe.Equal(got, tt.want.df) {
				t.Errorf("ReadInterface() got \n%v, \nwant \n%v", got, tt.want.df)
			}

			if strings.Contains(tt.name, "soft fail") {
				if buf.String() == "" {
					t.Errorf("pd.ReadInterface() returned no log message, want log due to fail")
				}
			}
		})
	}
}
