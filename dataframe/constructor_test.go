package dataframe

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
)

func TestNew_emptyDataFrame(t *testing.T) {
	got := newEmptyDataFrame()
	want := &DataFrame{vals: []values.Container{}, index: index.New(), cols: index.NewColumns()}
	if !Equal(got, want) {
		t.Errorf("New(nil) returned %#v, want %#v", got, want)
	}
	_ = got.Len()
	_ = got.ColLevels()
	_ = got.IndexLevels()
	_ = got.NumCols()
	_ = got.Name()
	_ = got.ensureAlignment()
	_ = got.valsAligned()
	_ = got.ensureColumnLevelPositions([]int{})
	_ = got.ensureColumnPositions([]int{})
	_ = got.ensureIndexLevelPositions([]int{})
	_ = got.ensureRowPositions([]int{})
	_ = got.Index.Len()

}

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
				vals:  []values.Container{values.MustCreateValuesFromInterface("foo"), values.MustCreateValuesFromInterface("bar")},
				cols:  index.NewDefaultColumns(2),
				index: index.NewDefault(1),
			},
		},
		{"config with name",
			args{[]interface{}{"foo", "bar"}, []Config{Config{Name: "foobar"}}},
			&DataFrame{
				name:  "foobar",
				vals:  []values.Container{values.MustCreateValuesFromInterface("foo"), values.MustCreateValuesFromInterface("bar")},
				cols:  index.NewDefaultColumns(2),
				index: index.NewDefault(1),
			},
		},
		{"config with datatype",
			args{[]interface{}{1.0, 2.0}, []Config{Config{DataType: options.Int64}}},
			&DataFrame{
				name:  "",
				vals:  []values.Container{values.MustCreateValuesFromInterface(int64(1)), values.MustCreateValuesFromInterface(int64(2))},
				cols:  index.NewDefaultColumns(2),
				index: index.NewDefault(1),
			},
		},
		{"config with index",
			args{[]interface{}{"foo", "bar"}, []Config{Config{Index: "baz"}}},
			&DataFrame{
				name:  "",
				vals:  []values.Container{values.MustCreateValuesFromInterface("foo"), values.MustCreateValuesFromInterface("bar")},
				cols:  index.NewDefaultColumns(2),
				index: index.New(index.MustNewLevel("baz", "")),
			},
		},
		{"config with named index",
			args{[]interface{}{"foo", "bar"}, []Config{Config{Index: "baz", IndexName: "corge"}}},
			&DataFrame{
				name:  "",
				vals:  []values.Container{values.MustCreateValuesFromInterface("foo"), values.MustCreateValuesFromInterface("bar")},
				cols:  index.NewDefaultColumns(2),
				index: index.New(index.MustNewLevel("baz", "corge")),
			},
		},
		{"config with multiIndex",
			args{[]interface{}{"foo", "bar"}, []Config{Config{MultiIndex: []interface{}{"baz", "qux"}}}},
			&DataFrame{
				name:  "",
				vals:  []values.Container{values.MustCreateValuesFromInterface("foo"), values.MustCreateValuesFromInterface("bar")},
				cols:  index.NewDefaultColumns(2),
				index: index.New(index.MustNewLevel("baz", ""), index.MustNewLevel("qux", "")),
			},
		},
		{"config with named multiIndex",
			args{[]interface{}{"foo", "bar"}, []Config{Config{MultiIndex: []interface{}{"baz", "qux"}, MultiIndexNames: []string{"waldo", "fred"}}}},
			&DataFrame{
				name:  "",
				vals:  []values.Container{values.MustCreateValuesFromInterface("foo"), values.MustCreateValuesFromInterface("bar")},
				cols:  index.NewDefaultColumns(2),
				index: index.New(index.MustNewLevel("baz", "waldo"), index.MustNewLevel("qux", "fred")),
			},
		},
		{"config with columns",
			args{[]interface{}{"foo", "bar"}, []Config{Config{Col: []string{"baz", "qux"}}}},
			&DataFrame{
				name:  "",
				vals:  []values.Container{values.MustCreateValuesFromInterface("foo"), values.MustCreateValuesFromInterface("bar")},
				cols:  index.NewColumns(index.NewColLevel([]string{"baz", "qux"}, "")),
				index: index.NewDefault(1),
			},
		},
		{"config with named columns",
			args{[]interface{}{"foo", "bar"}, []Config{Config{Col: []string{"baz", "qux"}, ColName: "corge"}}},
			&DataFrame{
				name:  "",
				vals:  []values.Container{values.MustCreateValuesFromInterface("foo"), values.MustCreateValuesFromInterface("bar")},
				cols:  index.NewColumns(index.NewColLevel([]string{"baz", "qux"}, "corge")),
				index: index.NewDefault(1),
			},
		},
		{"config with multicolumns",
			args{[]interface{}{"foo", "bar"}, []Config{Config{MultiCol: [][]string{{"baz", "qux"}, {"quux", "quuz"}}, MultiColNames: []string{"fred", "waldo"}}}},
			&DataFrame{
				name:  "",
				vals:  []values.Container{values.MustCreateValuesFromInterface("foo"), values.MustCreateValuesFromInterface("bar")},
				cols:  index.NewColumns(index.NewColLevel([]string{"baz", "qux"}, "fred"), index.NewColLevel([]string{"quux", "quuz"}, "waldo")),
				index: index.NewDefault(1),
			},
		},
		{"map[string]interface without config",
			args{[]interface{}{map[string]interface{}{"foo": []string{"bar", "baz"}}}, nil},
			&DataFrame{
				name:  "",
				vals:  []values.Container{values.MustCreateValuesFromInterface([]string{"bar", "baz"})},
				cols:  index.NewColumns(index.NewColLevel([]string{"foo"}, "")),
				index: index.NewDefault(2),
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
			if !Equal(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO fix to discriminate when test is failing
func TestDataFrame_constructor_internalReferences(t *testing.T) {
	tests := []struct {
		name  string
		input *DataFrame
	}{
		{"New", MustNew([]interface{}{"foo"})},
		{"Copy", MustNew([]interface{}{"foo"}).Copy()},
		{"newFromComponents", newFromComponents(nil, index.Index{}, index.Columns{}, "")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input.Columns.df == nil {
				t.Errorf("Constructor did not initialize Columns correctly")
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
		{"values-cols alignment: too many columns", args{[]interface{}{"foo"}, Config{Col: []string{"foo", "bar"}}}},
		{"values-cols alignment: too few columns", args{[]interface{}{"foo", "bar"}, Config{Col: []string{"baz"}}}},
		{"cols-multicols ambiguity", args{[]interface{}{"foo"}, Config{Col: []string{"baz"}, MultiCol: [][]string{{"qux"}}}}},
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
				name: "", vals: []values.Container{values.MustCreateValuesFromInterface("foo")},
				cols: index.NewDefaultColumns(1), index: index.NewDefault(1)}},
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
