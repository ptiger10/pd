package series

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
)

func TestNew_emptySeries(t *testing.T) {
	got := newEmptySeries()
	want := &Series{values: values.MustCreateValuesFromInterface(nil).Values, index: index.New(), datatype: options.None}
	if !Equal(got, want) {
		t.Errorf("New(nil) returned %#v, want %#v", got, want)
	}
}

func TestNew_nilWithConfig_emptySeries(t *testing.T) {
	got, err := New(nil, Config{Index: "foo"})
	if err != nil {
		t.Errorf("New(): %v", err)
	}
	want := newEmptySeries()
	if !Equal(got, want) {
		t.Errorf("New(nil) returned %#v, want %#v", got, want)
	}
}

func TestNew(t *testing.T) {
	testDate := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	type args struct {
		data interface{}
	}
	type want struct {
		values interface{}
		dtype  options.DataType
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"all null", args{""}, want{"", options.String}},
		{"float32", args{float32(1)}, want{1.0, options.Float64}},
		{"float64", args{float64(1)}, want{1.0, options.Float64}},
		{"int", args{int(1)}, want{1, options.Int64}},
		{"int8", args{int8(1)}, want{1, options.Int64}},
		{"int16", args{int16(1)}, want{1, options.Int64}},
		{"int32", args{int32(1)}, want{1, options.Int64}},
		{"int64", args{int64(1)}, want{1, options.Int64}},
		{"string", args{"foo"}, want{"foo", options.String}},
		{"bool", args{true}, want{true, options.Bool}},
		{"datetime", args{testDate}, want{testDate, options.DateTime}},

		{"float32_slice", args{[]float32{1}}, want{1.0, options.Float64}},
		{"float64_slice", args{[]float64{1}}, want{1.0, options.Float64}},
		{"int_slice", args{[]int{1}}, want{1, options.Int64}},
		{"int8_slice", args{[]int8{1}}, want{1, options.Int64}},
		{"int16_slice", args{[]int16{1}}, want{1, options.Int64}},
		{"int32_slice", args{[]int32{1}}, want{1, options.Int64}},
		{"int64_slice", args{[]int64{1}}, want{1, options.Int64}},
		{"string_slice", args{[]string{"foo"}}, want{"foo", options.String}},
		{"bool_slice", args{[]bool{true}}, want{true, options.Bool}},
		{"datetime_slice", args{[]time.Time{testDate}}, want{testDate, options.DateTime}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.data)
			if err != nil {
				t.Errorf("New() error = %v, wantErr nil", err)
			}
			container := values.MustCreateValuesFromInterface(tt.want.values)
			wantValues := container.Values
			wantIdx := index.NewDefault(1)
			want := &Series{values: wantValues, index: wantIdx, datatype: tt.want.dtype}
			if !Equal(got, want) {
				t.Errorf("New() = %v, want %v", got, want)
			}
		})
	}
}

func TestNew_conversion(t *testing.T) {
	got, err := New("3.5", Config{DataType: options.Float64})
	if err != nil {
		t.Errorf("New(): %v", err)
	}
	values, _ := values.InterfaceFactory(3.5)
	index := index.NewDefault(1)
	want := &Series{values: values.Values, index: index, datatype: options.Float64}
	if !Equal(got, want) {
		t.Errorf("New(nil) returned %v, want %v", got, want)
	}
}

func TestNew_Fail(t *testing.T) {
	type args struct {
		data   interface{}
		config Config
	}
	tests := []struct {
		name string
		args args
	}{
		{"unsupported value", args{complex64(1), Config{}}},
		{"unsupported single index", args{"foo", Config{Index: complex64(1)}}},
		{"unsupported multiIndex", args{"foo", Config{MultiIndex: []interface{}{complex64(1)}}}},
		{"unsupported conversion", args{"3.5", Config{DataType: options.Unsupported}}},
		{"index-multiIndex ambiguity", args{"foo", Config{Index: "foo", MultiIndex: []interface{}{"bar"}}}},
		{"values-index alignmentV1", args{"foo", Config{Index: []string{"foo", "bar"}}}},
		{"values-index alignmentV2", args{[]string{"foo"}, Config{Index: []string{"foo", "bar"}}}},
		{"values-index alignmentV3", args{[]string{"foo", "bar"}, Config{Index: "foo"}}},
		{"values-index alignmentV4", args{[]string{"foo", "bar"}, Config{Index: []string{"foo"}}}},
		{"values-multiIndex alignmentV1", args{"foo", Config{MultiIndex: []interface{}{[]string{"foo", "bar"}}}}},
		{"values-multiIndex alignment2", args{[]string{"foo"}, Config{MultiIndex: []interface{}{[]string{"foo", "bar"}}}}},
		{"values-multiIndex alignmentV3", args{[]string{"foo", "bar"}, Config{MultiIndex: []interface{}{"foo"}}}},
		{"values-multiIndex alignmentV4", args{[]string{"foo", "bar"}, Config{MultiIndex: []interface{}{"foo"}}}},
		{"values-multiIndex alignmentV5", args{[]string{"foo", "bar"}, Config{MultiIndex: []interface{}{"foo", "bar"}}}},
		{"multiIndex alignment", args{[]string{"foo", "bar"}, Config{
			MultiIndex: []interface{}{[]string{"foo", "bar"}, []string{"baz"}}}}},
		{"multiIndex names", args{[]string{"foo", "bar"}, Config{
			MultiIndex:      []interface{}{[]string{"foo", "bar"}, []string{"baz", "qux"}},
			MultiIndexNames: []string{"1"},
		}}},
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
	_, err := New("foo", Config{}, Config{})
	if err == nil {
		t.Error("New() error = nil, want error due to multiple configs")
	}
}

func TestMustNew(t *testing.T) {
	v, _ := values.InterfaceFactory(1.0)
	tests := []struct {
		name string
		args interface{}
		want *Series
	}{
		{name: "pass", args: 1.0,
			want: &Series{values: v.Values, index: index.NewDefault(1), datatype: options.Float64}},
		{name: "fail", args: complex64(1),
			want: newEmptySeries()},
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
func TestMustNew_fail(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)
	MustNew(complex64(1))
	if buf.String() == "" {
		t.Errorf("MustNew() returned no log message, want log due to fail")
	}
}

func Test_Copy(t *testing.T) {
	tests := []struct {
		name  string
		input *Series
		want  *Series
	}{
		{name: "pass", input: MustNew("foo"), want: MustNew("foo")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Copy()
			if !Equal(got, tt.want) {
				t.Errorf("s.Copy() returned %v, want %v", got, tt.want)
			}
			if reflect.ValueOf(tt.input).Pointer() == reflect.ValueOf(tt.want).Pointer() {
				t.Errorf("s.Copy() retained reference to original, want copy")
			}
			if reflect.ValueOf(tt.input.values).Pointer() == reflect.ValueOf(tt.want.values).Pointer() {
				t.Errorf("s.Copy() retained reference to original values, want copy")
			}
			if reflect.ValueOf(tt.input.index.Levels).Pointer() == reflect.ValueOf(tt.want.index.Levels).Pointer() {
				t.Errorf("s.Copy() retained reference to original index, want copy")
			}
		})
	}
}

func TestFromInternalComponents(t *testing.T) {
	vals := values.MustCreateValuesFromInterface("foo")
	index := index.NewDefault(1)
	got := FromInternalComponents(vals, index, "bar")
	want := MustNew("foo", Config{Name: "bar"})
	if !Equal(got, want) {
		t.Errorf("FromInternalComponents() returned %v, want %v", got, want)
	}

}

func TestToInternalComponents(t *testing.T) {
	s := MustNew("foo")
	vals, idx := s.ToInternalComponents()
	wantVals := values.MustCreateValuesFromInterface("foo")
	wantIdx := index.NewDefault(1)
	if !reflect.DeepEqual(vals, wantVals) {
		t.Errorf("ToInternalComponents() returned %v, want %v", vals, wantVals)
	}
	if !reflect.DeepEqual(idx, wantIdx) {
		t.Errorf("ToInternalComponents() returned %v, want %v", idx, wantIdx)
	}

}
