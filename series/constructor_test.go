package series

import (
	"bytes"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
)

func TestNew_nil(t *testing.T) {
	got, err := New(nil)
	if err != nil {
		t.Errorf("New() error = %v, wantErr nil", err)
	}
	values, _ := values.InterfaceFactory(nil)
	index := index.New()
	want := &Series{values: values.Values, index: index}
	if !Equal(got, want) {
		t.Errorf("New(nil) returned %#v, want %#v", got, want)
	}
}

func TestNew_allNull(t *testing.T) {
	got, err := New([]string{""})
	if err != nil {
		t.Errorf("New() error = %v, wantErr nil", err)
	}
	values, _ := values.InterfaceFactory("")
	index := index.NewDefault(1)
	want := &Series{values: values.Values, index: index, datatype: options.String}
	if !Equal(got, want) {
		t.Errorf("New(nil) returned %v, want %v", got, want)
	}
}

func TestNew_unsupported(t *testing.T) {
	_, err := New(complex64(1))
	if err == nil {
		t.Error("New() error = nil, want error due to unsupported input type")
	}
}

func TestNew_scalarTypeConversion(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    options.DataType
		wantErr bool
	}{
		{"float32", args{float32(1)}, options.Float64, false},
		{"float64", args{float64(1)}, options.Float64, false},
		{"int", args{int(1)}, options.Int64, false},
		{"int8", args{int8(1)}, options.Int64, false},
		{"int16", args{int16(1)}, options.Int64, false},
		{"int32", args{int32(1)}, options.Int64, false},
		{"int64", args{int64(1)}, options.Int64, false},
		{"string", args{"foo"}, options.String, false},
		{"bool", args{true}, options.Bool, false},
		{"datetime", args{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)}, options.DateTime, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.datatype != tt.want {
				t.Errorf("New().datatype = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew_sliceTypeConversion(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    options.DataType
		wantErr bool
	}{
		{"float32", args{[]float32{1}}, options.Float64, false},
		{"float64", args{[]float64{1}}, options.Float64, false},
		{"int", args{[]int{1}}, options.Int64, false},
		{"int8", args{int8(1)}, options.Int64, false},
		{"int16", args{int16(1)}, options.Int64, false},
		{"int32", args{int32(1)}, options.Int64, false},
		{"int64", args{int64(1)}, options.Int64, false},
		{"string", args{"foo"}, options.String, false},
		{"bool", args{true}, options.Bool, false},
		{"datetime", args{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)}, options.DateTime, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.datatype != tt.want {
				t.Errorf("New().datatype = %v, want %v", got.datatype, tt.want)
			}
		})
	}
}

func TestNew_float(t *testing.T) {
	got, err := New(1.0)
	if err != nil {
		t.Errorf("New() error = %v, wantErr nil", err)
	}
	v, err := values.InterfaceFactory(1.0)
	if err != nil {
		t.Error(err)
	}
	idx := index.NewDefault(1)
	want := &Series{values: v.Values, index: idx, datatype: options.Float64}
	if !Equal(got, want) {
		t.Errorf("New() = %v, want %v", got, want)
	}
}

func TestNew_int(t *testing.T) {
	got, err := New(1)
	if err != nil {
		t.Errorf("New() error = %v, wantErr nil", err)
	}
	v, err := values.InterfaceFactory(1)
	if err != nil {
		t.Error(err)
	}
	idx := index.NewDefault(1)
	want := &Series{values: v.Values, index: idx, datatype: options.Int64}
	if !Equal(got, want) {
		t.Errorf("New() = %#v, want %#v", got, want)
	}
}

func TestNew_string(t *testing.T) {
	got, err := New("foo")
	if err != nil {
		t.Errorf("New() error = %v, wantErr nil", err)
	}
	v, err := values.InterfaceFactory("foo")
	if err != nil {
		t.Error(err)
	}
	idx := index.NewDefault(1)
	want := &Series{values: v.Values, index: idx, datatype: options.String}
	if !Equal(got, want) {
		t.Errorf("New() = %#v, want %#v", got, want)
	}
}

func TestNew_bool(t *testing.T) {
	got, err := New(true)
	if err != nil {
		t.Errorf("New() error = %v, wantErr nil", err)
	}
	v, err := values.InterfaceFactory(true)
	if err != nil {
		t.Error(err)
	}
	idx := index.NewDefault(1)
	want := &Series{values: v.Values, index: idx, datatype: options.Bool}
	if !Equal(got, want) {
		t.Errorf("New() = %#v, want %#v", got, want)
	}
}

func TestNew_datetime(t *testing.T) {
	got, err := New(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Errorf("New() error = %v, wantErr nil", err)
	}
	v, err := values.InterfaceFactory(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Error(err)
	}
	idx := index.NewDefault(1)
	want := &Series{values: v.Values, index: idx, datatype: options.DateTime}
	if !Equal(got, want) {
		t.Errorf("New() = %#v, want %#v", got, want)
	}
}

func TestNew_multiIndex(t *testing.T) {
	got, _ := New("foo", Config{
		MultiIndex: []interface{}{"bar", "baz"},
	})
	v, err := values.InterfaceFactory("foo")
	if err != nil {
		t.Error(err)
	}
	idx1, _ := index.NewLevel("bar", "")
	idx2, _ := index.NewLevel("baz", "")
	idx := index.New(idx1, idx2)
	want := &Series{values: v.Values, index: idx, datatype: options.String}
	if !Equal(got, want) {
		t.Errorf("New() = %v, want %v", got, want)
	}
}

func TestNew_multiIndex_named(t *testing.T) {
	got, _ := New("foo", Config{
		MultiIndex:      []interface{}{"bar", "baz"},
		MultiIndexNames: []string{"qux", "quux"},
	})
	v, err := values.InterfaceFactory("foo")
	if err != nil {
		t.Error(err)
	}
	idx1, _ := index.NewLevel("bar", "qux")
	idx2, _ := index.NewLevel("baz", "quux")
	idx := index.New(idx1, idx2)
	want := &Series{values: v.Values, index: idx, datatype: options.String}
	if !Equal(got, want) {
		t.Errorf("New() = %v, want %v", got, want)
	}
}

// Alignment tests
func TestNew_configFailure(t *testing.T) {
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

func TestNew_multiple_configs(t *testing.T) {
	_, err := New(nil, Config{}, Config{})
	if err == nil {
		t.Error("New() error = nil, want error due to multiple configs")
	}
}

func TestNew_nil_withConfig(t *testing.T) {
	got, err := New(nil, Config{})
	if err != nil {
		t.Errorf("New(): %v", err)
	}
	values, _ := values.InterfaceFactory(nil)
	index := index.New()
	want := &Series{values: values.Values, index: index}
	if !Equal(got, want) {
		t.Errorf("New(nil) returned %#v, want %#v", got, want)
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

func TestMustNew_success(t *testing.T) {
	got := MustNew(1.0)
	v, err := values.InterfaceFactory(1.0)
	if err != nil {
		t.Error(err)
	}
	idx := index.NewDefault(1)
	want := &Series{values: v.Values, index: idx, datatype: options.Float64}
	if !Equal(got, want) {
		t.Errorf("MustNew() = %v, want %v", got, want)
	}
}
func TestMustNew_fail(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	MustNew(complex64(1))
	if buf.String() == "" {
		t.Errorf("MustNew() returned no log message, want log due to fail")
	}
}

func TestMustNew_failWithConfig(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	MustNew(complex64(1), Config{})
	if buf.String() == "" {
		t.Errorf("mustNew() returned no log message, want log due to fail")
	}
}
