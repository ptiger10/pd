package series

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	// "github.com/d4l3k/messagediff"
	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
)

// // Float Tests
// // ------------------------------------------------
// func TestConstructor_Slice_Float(t *testing.T) {
// 	var err error

// 	_, err = New([]float32{-1.5, 0, 1.5})
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	_, err = New([]float64{-1.5, 0, 1.5})
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestConstructor_Name_Slice_Float(t *testing.T) {
// 	var err error

// 	_, err = New([]float32{-1.5, 0, 1.5}, options.Name("float32"))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	_, err = New([]float64{-1.5, 0, 1.5}, options.Name("float64"))
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestConstructor_Index_Slice_Float(t *testing.T) {
// 	var err error

// 	_, err = New([]float32{-1.5, 0, 1.5}, Idx([]float32{1, 2, 3}))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	_, err = New([]float64{-1.5, 0, 1.5}, Idx([]float64{1, 2, 3}))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	_, err = New([]float64{-1.5, 0, 1.5}, Idx([]float64{1}))
// 	if err == nil {
// 		t.Errorf("Returned nil error, want error due to mismatched value/index lengths")
// 	}
// }

// func TestConstructor_Kind_Slice_Float(t *testing.T) {
// 	var err error

// 	_, err = New([]float32{-1.5, 0, 1.5}, options.DataType(options.Float64))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	_, err = New([]float32{-1.5, 0, 1.5}, options.DataType(options.Int64))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	_, err = New([]float32{-1.5, 0, 1.5}, options.DataType(options.String))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	_, err = New([]float32{-1.5, 0, 1.5}, options.DataType(options.Bool))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	_, err = New([]float32{-1.5, 0, 1.5}, options.DataType(options.DateTime))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	_, err = New([]float32{-1.5, 0, 1.5}, options.DataType(options.Interface))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	_, err = New([]float32{-1.5, 0, 1.5}, options.DataType(options.Unsupported))
// 	if err == nil {
// 		t.Error("Returned nil error, want error due to unsupported conversion type")
// 	}

// }

// func TestMini_single(t *testing.T) {
// 	mini := config.MiniIndex{
// 		Data:     []int{1, 2, 3},
// 		DataType: options.Int64,
// 		Name:     "test",
// 	}
// 	got, err := indexFromMiniIndex([]config.MiniIndex{mini}, 3)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	lvl, err := index.NewLevel([]int64{1, 2, 3}, "test")
// 	want := index.New(lvl)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("MiniIndex returned %v, want %v", got, want)
// 	}

// }

// func TestNew_float(t *testing.T) {
// 	got, err := New(1.0)
// 	if err != nil {
// 		t.Errorf("New() error = %v, wantErr nil", err)
// 	}
// 	v, err := values.InterfaceFactory(1.0)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	idx := index.Default(1)
// 	want := &Series{values: v.Values, index: idx, datatype: options.Float64}
// 	if !seriesEquals(*got, *want) {
// 		t.Errorf("New() = %v, want %v", got, want)
// 	}
// }

func TestNew_test(t *testing.T) {
	s, _ := New("foo", Idx("bar"))
	s2, _ := New("foo", Idx("bar"))
	if !reflect.DeepEqual(s, s2) {
		t.Errorf("New() = %#v, want %#v", s, s2)
	}
}

func TestNew_multi(t *testing.T) {
	got, _ := New("foo", Idx("bar"), Idx("baz"))
	v, err := values.InterfaceFactory("foo")
	if err != nil {
		t.Error(err)
	}
	idx1, _ := index.NewLevel("bar", "")
	idx2, _ := index.NewLevel("baz", "")
	idx := index.New(idx1, idx2)
	want := &Series{values: v.Values, index: idx, datatype: options.String}
	if !seriesEquals(got, want) {
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
	idx := index.Default(1)
	want := &Series{values: v.Values, index: idx, datatype: options.Int64}
	if !seriesEquals(got, want) {
		t.Errorf("New() = %#v, want %#v", got, want)
	}
	// diff, _ := messagediff.PrettyDiff(got, want)
	// fmt.Println(diff)
}

func TestNew_scalar(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *Series
		wantErr bool
	}{
		{"unsupported", args{complex64(1)}, nil, true},
		{"float32", args{float32(1)}, mustNew(float64(1)), false},
		{"float64", args{float64(1)}, mustNew(float64(1)), false},
		{"int", args{int(1)}, mustNew(int64(1)), false},
		{"int8", args{int8(1)}, mustNew(int64(1)), false},
		{"int16", args{int16(1)}, mustNew(int64(1)), false},
		{"int32", args{int32(1)}, mustNew(int64(1)), false},
		{"int64", args{int64(1)}, mustNew(int64(1)), false},
		{"string", args{"foo"}, mustNew("foo"), false},
		{"bool", args{true}, mustNew(true), false},
		{"datetime", args{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)}, mustNew(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew_slice(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *Series
		wantErr bool
	}{
		{"unsupported", args{[]complex64{1}}, nil, true},
		{"float32", args{[]float32{1}}, mustNew([]float64{1}), false},
		{"float64", args{[]float64{1}}, mustNew([]float64{1}), false},
		{"int", args{[]int{1}}, mustNew([]int64{1}), false},
		{"int8", args{int8(1)}, mustNew(int64(1)), false},
		{"int16", args{int16(1)}, mustNew(int64(1)), false},
		{"int32", args{int32(1)}, mustNew(int64(1)), false},
		{"int64", args{int64(1)}, mustNew(int64(1)), false},
		{"string", args{"foo"}, mustNew("foo"), false},
		{"bool", args{true}, mustNew(true), false},
		{"datetime", args{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)}, mustNew(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWithConfig(t *testing.T) {
	s, err := NewWithConfig(Config{}, int8(1))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(s)
}

// func TestNewCustom(t *testing.T) {
// 	type args struct {
// 		data   interface{}
// 		config *Config
// 		idx    []IndexLevel
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *Series
// 		wantErr bool
// 	}{
// 		{"no index, no config", args{data: "foo", config: nil}, mustNewCustom("foo", nil), false},
// 		{"no index, config", args{data: "foo", config: nil}, mustNewCustom("foo", nil), false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := NewCustom(tt.args.data, tt.args.config, tt.args.idx...)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("NewCustom() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NewCustom() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
