package values

import (
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/ptiger10/pd/options"
)

func TestInterfaceFactory(t *testing.T) {
	type args struct {
		data interface{}
	}
	type want struct {
		factory Factory
		err     bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"nil", args{nil},
			want{Factory{Values: &interfaceValues{}, DataType: options.None}, false}},
		{"scalar", args{"foo"},
			want{Factory{Values: &stringValues{stringValue{"foo", false}}, DataType: options.String}, false}},
		{"slice", args{[]string{"foo"}},
			want{Factory{Values: &stringValues{stringValue{"foo", false}}, DataType: options.String}, false}},
		{"fail: unsupported", args{[1]string{"foo"}},
			want{Factory{Values: &interfaceValues{}, DataType: options.None}, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InterfaceFactory(tt.args.data)
			if (err != nil) != tt.want.err {
				t.Errorf("InterfaceFactory() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want.factory) {
				t.Errorf("InterfaceFactory() = %v, want %v", got, tt.want.factory)
			}
		})
	}
}

func TestScalarFactory(t *testing.T) {
	type args struct {
		data interface{}
	}
	type want struct {
		factory Factory
		err     bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"float32",
			args{float32(1.0)},
			want{Factory{Values: &float64Values{float64Value{1.0, false}}, DataType: options.Float64}, false}},
		{"float64",
			args{float64(1.0)},
			want{Factory{Values: &float64Values{float64Value{1.0, false}}, DataType: options.Float64}, false}},
		{"int",
			args{int(1)},
			want{Factory{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"int8",
			args{int8(1)},
			want{Factory{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"int16",
			args{int16(1)},
			want{Factory{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"int32",
			args{int32(1)},
			want{Factory{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"int64",
			args{int64(1)},
			want{Factory{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"uint",
			args{uint(1)},
			want{Factory{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"uint8",
			args{uint8(1)},
			want{Factory{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"uint16",
			args{uint16(1)},
			want{Factory{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"uint32",
			args{uint32(1)},
			want{Factory{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"uint64",
			args{uint64(1)},
			want{Factory{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"string",
			args{"foo"},
			want{Factory{Values: &stringValues{stringValue{"foo", false}}, DataType: options.String}, false}},
		{"string null",
			args{""},
			want{Factory{Values: &stringValues{stringValue{"NaN", true}}, DataType: options.String}, false}},
		{"bool",
			args{true},
			want{Factory{Values: &boolValues{boolValue{true, false}}, DataType: options.Bool}, false}},
		{"datetime",
			args{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)},
			want{Factory{Values: &dateTimeValues{dateTimeValue{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC), false}}, DataType: options.DateTime}, false}},
		{"unsupported",
			args{complex64(1)},
			want{Factory{Values: nil, DataType: options.None}, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ScalarFactory(tt.args.data)
			if (err != nil) != tt.want.err {
				t.Errorf("ScalarFactory() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want.factory) {
				t.Errorf("ScalarFactory() = %v, want %v", got, tt.want.factory)
			}
		})
	}
}

func TestScalarConstructor_NullFloat(t *testing.T) {
	vals, err := ScalarFactory(math.NaN())
	if err != nil {
		t.Errorf("Unable to construct values from null float: %v", err)
	}
	val := vals.Values.Element(0).Value.(float64)
	if !math.IsNaN(val) {
		t.Errorf("Returned %v, want NaN", val)
	}
}

func TestSliceFactory(t *testing.T) {
	type args struct {
		data interface{}
	}
	type want struct {
		factory Factory
		err     bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"float32",
			args{[]float32{1.0}},
			want{Factory{Values: &float64Values{float64Value{1.0, false}}, DataType: options.Float64}, false}},
		{"float64",
			args{[]float64{1.0}},
			want{Factory{Values: &float64Values{float64Value{1.0, false}}, DataType: options.Float64}, false}},
		{"int",
			args{[]int{1}},
			want{Factory{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"int8",
			args{[]int8{1}},
			want{Factory{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"int16",
			args{[]int16{1}},
			want{Factory{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"int32",
			args{[]int32{1}},
			want{Factory{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"int64",
			args{[]int64{1}},
			want{Factory{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"uint",
			args{[]uint{1}},
			want{Factory{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"uint8",
			args{[]uint8{1}},
			want{Factory{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"uint16",
			args{[]uint16{1}},
			want{Factory{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"uint32",
			args{[]uint32{1}},
			want{Factory{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"uint64",
			args{[]uint64{1}},
			want{Factory{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"string",
			args{[]string{"foo"}},
			want{Factory{Values: &stringValues{stringValue{"foo", false}}, DataType: options.String}, false}},
		{"string null",
			args{[]string{""}},
			want{Factory{Values: &stringValues{stringValue{"NaN", true}}, DataType: options.String}, false}},
		{"bool",
			args{[]bool{true}},
			want{Factory{Values: &boolValues{boolValue{true, false}}, DataType: options.Bool}, false}},
		{"datetime",
			args{[]time.Time{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)}},
			want{Factory{Values: &dateTimeValues{dateTimeValue{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC), false}}, DataType: options.DateTime}, false}},
		{"unsupported",
			args{[]complex64{complex64(1)}},
			want{Factory{Values: nil, DataType: options.None}, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SliceFactory(tt.args.data)
			if (err != nil) != tt.want.err {
				t.Errorf("SliceFactory() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want.factory) {
				t.Errorf("SliceFactory() = %v, want %v", got, tt.want.factory)
			}
		})
	}
}

func TestSliceConstructor_NullFloat(t *testing.T) {
	vals, err := SliceFactory([]float64{math.NaN()})
	if err != nil {
		t.Errorf("Unable to construct values from null float: %v", err)
	}
	val := vals.Values.Element(0).Value.(float64)
	if !math.IsNaN(val) {
		t.Errorf("Returned %v, want NaN", val)
	}
}

func TestSliceConstructor_NullFloatInterface(t *testing.T) {
	vals, err := SliceFactory([]interface{}{math.NaN()})
	if err != nil {
		t.Errorf("Unable to construct values from null float: %v", err)
	}
	val := vals.Values.Element(0).Value.(float64)
	if !math.IsNaN(val) {
		t.Errorf("Returned %v, want NaN", val)
	}
}

func TestSliceConstructor_Unsupported(t *testing.T) {
	data := []complex64{1, 2, 3}
	_, err := SliceFactory(data)
	if err == nil {
		t.Errorf("Returned nil error, expected error due to unsupported type %T", data)
	}
}

func TestMakeIntRange(t *testing.T) {
	got := MakeIntRange(0, 3)
	want := []int{0, 1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MakeIntRange(): got %v, want %v", got, want)
	}
}

func TestMakeIntRangeInclusive(t *testing.T) {
	type args struct {
		start int
		end   int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"ascending", args{0, 3}, []int{0, 1, 2, 3}},
		{"descending", args{3, 0}, []int{3, 2, 1, 0}},
		{"equal", args{1, 1}, []int{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MakeIntRangeInclusive(tt.args.start, tt.args.end)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeIntRange(): got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeInterfaceRange(t *testing.T) {
	got := MakeInterfaceRange(0, 3)
	want := []interface{}{0, 1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MakeIntRange(): got %v, want %v", got, want)
	}
}
