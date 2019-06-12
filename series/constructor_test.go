package series

import (
	"reflect"
	"testing"
	"time"

	"github.com/ptiger10/pd/internal/config"
	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/kinds"
	"github.com/ptiger10/pd/opt"
)

// Float Tests
// ------------------------------------------------
func TestConstructor_Slice_Float(t *testing.T) {
	var err error

	_, err = New([]float32{-1.5, 0, 1.5})
	if err != nil {
		t.Error(err)
	}

	_, err = New([]float64{-1.5, 0, 1.5})
	if err != nil {
		t.Error(err)
	}
}

func TestConstructor_Name_Slice_Float(t *testing.T) {
	var err error

	_, err = New([]float32{-1.5, 0, 1.5}, opt.Name("float32"))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float64{-1.5, 0, 1.5}, opt.Name("float64"))
	if err != nil {
		t.Error(err)
	}
}

func TestConstructor_Index_Slice_Float(t *testing.T) {
	var err error

	_, err = New([]float32{-1.5, 0, 1.5}, Idx([]float32{1, 2, 3}))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float64{-1.5, 0, 1.5}, Idx([]float64{1, 2, 3}))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float64{-1.5, 0, 1.5}, Idx([]float64{1}))
	if err == nil {
		t.Errorf("Returned nil error, want error due to mismatched value/index lengths")
	}
}

func TestConstructor_Kind_Slice_Float(t *testing.T) {
	var err error

	_, err = New([]float32{-1.5, 0, 1.5}, opt.Kind(kinds.Float64))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, opt.Kind(kinds.Int64))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, opt.Kind(kinds.String))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, opt.Kind(kinds.Bool))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, opt.Kind(kinds.DateTime))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, opt.Kind(kinds.Interface))
	if err != nil {
		t.Error(err)
	}
	_, err = New([]float32{-1.5, 0, 1.5}, opt.Kind(kinds.Unsupported))
	if err == nil {
		t.Error("Returned nil error, want error due to unsupported conversion type")
	}

}

func TestMini_single(t *testing.T) {
	mini := config.MiniIndex{
		Data: []int{1, 2, 3},
		Kind: kinds.Int64,
		Name: "test",
	}
	got, err := indexFromMiniIndex([]config.MiniIndex{mini}, 3)
	if err != nil {
		t.Error(err)
	}
	lvl, err := index.NewLevel([]int64{1, 2, 3}, "test")
	want := index.New(lvl)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MiniIndex returned %v, want %v", got, want)
	}

}

func TestNew2_float(t *testing.T) {
	got, err := New2(1.0)
	if err != nil {
		t.Errorf("New() error = %v, wantErr nil", err)
	}
	v, err := values.InterfaceFactory(1.0)
	if err != nil {
		t.Error(err)
	}
	idx := index.Default(1)
	want := &Series{values: v.Values, index: idx, kind: kinds.Float64}
	if !seriesEquals(*got, *want) {
		t.Errorf("New() = %v, want %v", got, want)
	}
}

func TestNew2_int(t *testing.T) {
	got, err := New2(1)
	if err != nil {
		t.Errorf("New() error = %v, wantErr nil", err)
	}
	v, err := values.InterfaceFactory(1)
	if err != nil {
		t.Error(err)
	}
	idx := index.Default(1)
	want := &Series{values: v.Values, index: idx, kind: kinds.Int64}
	if !seriesEquals(*got, *want) {
		t.Errorf("New() = %v, want %v", got, want)
	}
}

func TestNew2_scalar(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *Series
		wantErr bool
	}{
		{"float32", args{float32(1)}, mustNew2(float64(1)), false},
		{"float64", args{float64(1)}, mustNew2(float64(1)), false},
		{"int", args{int(1)}, mustNew2(int64(1)), false},
		{"int8", args{int8(1)}, mustNew2(int64(1)), false},
		{"int16", args{int16(1)}, mustNew2(int64(1)), false},
		{"int32", args{int32(1)}, mustNew2(int64(1)), false},
		{"int64", args{int64(1)}, mustNew2(int64(1)), false},
		{"string", args{"foo"}, mustNew2("foo"), false},
		{"bool", args{true}, mustNew2(true), false},
		{"datetime", args{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)}, mustNew2(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New2(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("New2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew2_Slice(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *Series
		wantErr bool
	}{
		{"float32", args{[]float32{1}}, mustNew2([]float64{1}), false},
		{"float64", args{[]float64{1}}, mustNew2([]float64{1}), false},
		{"int", args{[]int{1}}, mustNew2([]int64{1}), false},
		{"int8", args{int8(1)}, mustNew2(int64(1)), false},
		{"int16", args{int16(1)}, mustNew2(int64(1)), false},
		{"int32", args{int32(1)}, mustNew2(int64(1)), false},
		{"int64", args{int64(1)}, mustNew2(int64(1)), false},
		{"string", args{"foo"}, mustNew2("foo"), false},
		{"bool", args{true}, mustNew2(true), false},
		{"datetime", args{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)}, mustNew2(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New2(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("New2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New2() = %v, want %v", got, tt.want)
			}
		})
	}
}
