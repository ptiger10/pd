package dataframe

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/options"
	"github.com/ptiger10/pd/series"
)

func TestNew(t *testing.T) {
	type args struct {
		data   []interface{}
		config []Config
	}
	tests := []struct {
		name string
		args args
		want *DataFrame
	}{
		{name: "no config",
			args: args{data: []interface{}{"foo", "bar"}, config: nil},
			want: &DataFrame{
				name:  "",
				s:     []*series.Series{series.MustNew("foo", series.Config{Name: "0"}), series.MustNew("bar", series.Config{Name: "1"})},
				cols:  index.NewDefaultColumns(2),
				index: index.NewDefault(1),
			},
		},
		{"config with name",
			args{[]interface{}{"foo", "bar"}, []Config{Config{Name: "foobar"}}},
			&DataFrame{
				name: "foobar",
				s: []*series.Series{
					series.MustNew("foo", series.Config{Name: "0"}),
					series.MustNew("bar", series.Config{Name: "1"})},
				cols:  index.NewDefaultColumns(2),
				index: index.NewDefault(1),
			},
		},
		{"config with index",
			args{[]interface{}{"foo", "bar"}, []Config{Config{Index: "baz"}}},
			&DataFrame{
				name: "",
				s: []*series.Series{
					series.MustNew("foo", series.Config{Name: "0", Index: "baz"}),
					series.MustNew("bar", series.Config{Name: "1", Index: "baz"})},
				cols:  index.NewDefaultColumns(2),
				index: index.New(index.MustNewLevel("baz", "")),
			},
		},
		{"config with named index",
			args{[]interface{}{"foo", "bar"}, []Config{Config{Index: "baz", IndexName: "corge"}}},
			&DataFrame{
				name: "",
				s: []*series.Series{
					series.MustNew("foo", series.Config{Name: "0", Index: "baz", IndexName: "corge"}),
					series.MustNew("bar", series.Config{Name: "1", Index: "baz", IndexName: "corge"})},
				cols:  index.NewDefaultColumns(2),
				index: index.New(index.MustNewLevel("baz", "corge")),
			},
		},
		{"config with multiIndex",
			args{[]interface{}{"foo", "bar"}, []Config{Config{MultiIndex: []interface{}{"baz", "qux"}}}},
			&DataFrame{
				name: "",
				s: []*series.Series{
					series.MustNew("foo", series.Config{Name: "0", MultiIndex: []interface{}{"baz", "qux"}}),
					series.MustNew("bar", series.Config{Name: "1", MultiIndex: []interface{}{"baz", "qux"}})},
				cols:  index.NewDefaultColumns(2),
				index: index.New(index.MustNewLevel("baz", ""), index.MustNewLevel("qux", "")),
			},
		},
		{"config with named multiIndex",
			args{[]interface{}{"foo", "bar"}, []Config{Config{MultiIndex: []interface{}{"baz", "qux"}, MultiIndexNames: []string{"waldo", "fred"}}}},
			&DataFrame{
				name: "",
				s: []*series.Series{
					series.MustNew("foo", series.Config{Name: "0", MultiIndex: []interface{}{"baz", "qux"}, MultiIndexNames: []string{"waldo", "fred"}}),
					series.MustNew("bar", series.Config{Name: "1", MultiIndex: []interface{}{"baz", "qux"}, MultiIndexNames: []string{"waldo", "fred"}})},
				cols:  index.NewDefaultColumns(2),
				index: index.New(index.MustNewLevel("baz", "waldo"), index.MustNewLevel("qux", "fred")),
			},
		},
		{"config with columns",
			args{[]interface{}{"foo", "bar"}, []Config{Config{Cols: []interface{}{"baz", "qux"}}}},
			&DataFrame{
				name: "",
				s: []*series.Series{
					series.MustNew("foo", series.Config{Name: "baz"}),
					series.MustNew("bar", series.Config{Name: "qux"})},
				cols:  index.NewColumns(index.NewColLevel([]interface{}{"baz", "qux"}, "")),
				index: index.NewDefault(1),
			},
		},
		{"config with named columns",
			args{[]interface{}{"foo", "bar"}, []Config{Config{Cols: []interface{}{"baz", "qux"}, ColsName: "corge"}}},
			&DataFrame{
				name: "",
				s: []*series.Series{
					series.MustNew("foo", series.Config{Name: "baz"}),
					series.MustNew("bar", series.Config{Name: "qux"})},
				cols:  index.NewColumns(index.NewColLevel([]interface{}{"baz", "qux"}, "corge")),
				index: index.NewDefault(1),
			},
		},
		{"config with multicolumns",
			args{[]interface{}{"foo", "bar"}, []Config{Config{MultiCol: [][]interface{}{{"baz", "qux"}, {"quux", "quuz"}}, MultiColNames: []string{"fred", "waldo"}}}},
			&DataFrame{
				name: "",
				s: []*series.Series{
					series.MustNew("foo", series.Config{Name: "baz | quux"}),
					series.MustNew("bar", series.Config{Name: "qux | quuz"})},
				cols:  index.NewColumns(index.NewColLevel([]interface{}{"baz", "qux"}, "fred"), index.NewColLevel([]interface{}{"quux", "quuz"}, "waldo")),
				index: index.NewDefault(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.data, tt.args.config...)
			if err != nil {
				t.Errorf("New() error = %v, want nil", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %#v, want %#v", got.index, tt.want.index)
			}
		})
	}
}

func TestNew_Fail(t *testing.T) {
	type args struct {
		data   []interface{}
		config Config
	}
	tests := []struct {
		name string
		args args
	}{
		{"unsupported value", args{[]interface{}{complex64(1)}, Config{}}},
		{"unsupported single index", args{[]interface{}{"foo"}, Config{Index: complex64(1)}}},
		{"unsupported multiIndex", args{[]interface{}{"foo"}, Config{MultiIndex: []interface{}{complex64(1)}}}},
		{"unsupported conversion", args{[]interface{}{"3.5"}, Config{DataType: options.Unsupported}}},
		{"mixed value slice lengths", args{[]interface{}{[]string{"foo"}, []string{"bar", "baz"}}, Config{}}},
		{"index-multiIndex ambiguity", args{[]interface{}{"foo"}, Config{Index: "foo", MultiIndex: []interface{}{"bar"}}}},
		{"values-index alignmentV1", args{[]interface{}{"foo"}, Config{Index: []string{"foo", "bar"}}}},
		{"values-index alignmentV2", args{[]interface{}{[]string{"foo"}}, Config{Index: []string{"foo", "bar"}}}},
		{"values-index alignmentV3", args{[]interface{}{[]string{"foo", "bar"}}, Config{Index: "foo"}}},
		{"values-index alignmentV4", args{[]interface{}{[]string{"foo", "bar"}}, Config{Index: []string{"foo"}}}},
		{"values-multiIndex alignmentV1", args{[]interface{}{"foo"}, Config{MultiIndex: []interface{}{[]string{"foo", "bar"}}}}},
		{"values-multiIndex alignment2", args{[]interface{}{[]string{"foo"}}, Config{MultiIndex: []interface{}{[]string{"foo", "bar"}}}}},
		{"values-multiIndex alignmentV3", args{[]interface{}{[]string{"foo", "bar"}}, Config{MultiIndex: []interface{}{"foo"}}}},
		{"values-multiIndex alignmentV4", args{[]interface{}{[]string{"foo", "bar"}}, Config{MultiIndex: []interface{}{"foo"}}}},
		{"values-multiIndex alignmentV5", args{[]interface{}{[]string{"foo", "bar"}}, Config{MultiIndex: []interface{}{"foo", "bar"}}}},
		{"multiIndex alignment", args{[]interface{}{[]string{"foo", "bar"}}, Config{
			MultiIndex: []interface{}{[]string{"foo", "bar"}, []string{"baz"}}}}},
		{"multiIndex names", args{[]interface{}{[]string{"foo", "bar"}}, Config{
			MultiIndex:      []interface{}{[]string{"foo", "bar"}, []string{"baz", "qux"}},
			MultiIndexNames: []string{"1"},
		}}},
		// {"column slice", args{[]interface{}{"foo"}, Config{Cols: []interface{}{[]string{"foo", "bar"}}}}},
		{"values-cols alignmentV1", args{[]interface{}{"foo"}, Config{Cols: []interface{}{"foo", "bar"}}}},
		// {"values-cols alignmentV2", args{[]interface{}{"foo", "bar"}, Config{Cols: []interface{}{"baz"}}}},
		{"cols-multicols ambiguity", args{[]interface{}{"foo"}, Config{Cols: []interface{}{"baz"}, MultiCol: [][]interface{}{{"qux"}}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.data, tt.args.config)
			if err == nil {
				t.Error("New() error = nil, want error")
				return
			}
		})
	}
}

func TestNew_Fail_multipleConfigs(t *testing.T) {
	_, err := New([]interface{}{"foo"}, Config{}, Config{})
	if err == nil {
		t.Error("New() error = nil, want error due to multiple configs")
	}
}

func TestMustNew(t *testing.T) {
	tests := []struct {
		name string
		args []interface{}
		want *DataFrame
	}{
		{name: "pass", args: []interface{}{"foo"},
			want: &DataFrame{
				name: "", s: []*series.Series{series.MustNew("foo", series.Config{Name: "0"})}, cols: index.NewDefaultColumns(1), index: index.NewDefault(1)}},
		{name: "fail", args: []interface{}{complex64(1)},
			want: newEmptyDataFrame()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			got := MustNew(tt.args)
			if !Equal(got, tt.want) {
				t.Errorf("MustNew() = %v, want %v", got, tt.want)
			}
			if strings.Contains(tt.name, "fail") {
				if buf.String() == "" {
					t.Errorf("series.MustNew() returned no log message, want log due to fail")
				}
			}
		})
	}
}
