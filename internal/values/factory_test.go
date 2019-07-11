package values

import (
	"bytes"
	"log"
	"math"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/ptiger10/pd/options"
)

func TestInterfaceFactory(t *testing.T) {
	type args struct {
		data interface{}
	}
	type want struct {
		container Container
		err       bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"nil", args{nil},
			want{Container{Values: &interfaceValues{}, DataType: options.None}, false}},
		{"scalar", args{"foo"},
			want{Container{Values: &stringValues{stringValue{"foo", false}}, DataType: options.String}, false}},
		{"slice", args{[]string{"foo"}},
			want{Container{Values: &stringValues{stringValue{"foo", false}}, DataType: options.String}, false}},
		{"fail: unsupported", args{[1]string{"foo"}},
			want{Container{Values: &interfaceValues{}, DataType: options.None}, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InterfaceFactory(tt.args.data)
			if (err != nil) != tt.want.err {
				t.Errorf("InterfaceFactory() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want.container) {
				t.Errorf("InterfaceFactory() = %v, want %v", got, tt.want.container)
			}
		})
	}
}

func TestMustCreateValuesFromInterface(t *testing.T) {
	type args struct {
		data interface{}
	}
	type want struct {
		container Container
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"nil", args{nil},
			want{Container{Values: &interfaceValues{}, DataType: options.None}}},
		{"scalar", args{"foo"},
			want{Container{Values: &stringValues{stringValue{"foo", false}}, DataType: options.String}}},
		{"slice", args{[]string{"foo"}},
			want{Container{Values: &stringValues{stringValue{"foo", false}}, DataType: options.String}}},
		{"fail: unsupported", args{[1]string{"foo"}},
			want{Container{Values: &interfaceValues{}, DataType: options.None}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(os.Stderr)

			got := MustCreateValuesFromInterface(tt.args.data)
			if !reflect.DeepEqual(got, tt.want.container) {
				t.Errorf("MustCreateValuesFromInterface() = %v, want %v", got, tt.want.container)
			}

			if strings.Contains(tt.name, "fail") {
				if buf.String() == "" {
					t.Errorf("Grouping.Group() returned no log message, want log due to fail")
				}
			}
		})
	}
}

func TestInterfaceSliceFactory(t *testing.T) {
	type args struct {
		data     []interface{}
		manual   bool
		dataType options.DataType
	}
	type want struct {
		vals []Container
		err  bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{name: "normal", args: args{data: []interface{}{"foo"}, manual: false, dataType: options.None},
			want: want{vals: []Container{{DataType: options.String, Values: &stringValues{stringValue{"foo", false}}}},
				err: false}},
		{name: "normal manual", args: args{data: []interface{}{"foo"}, manual: true, dataType: options.None},
			want: want{vals: []Container{{DataType: options.String, Values: &stringValues{stringValue{"foo", false}}}},
				err: false}},
		{"interpolate []interface", args{data: []interface{}{[]interface{}{"foo"}}, manual: false, dataType: options.None},
			want{vals: []Container{{DataType: options.String, Values: &stringValues{stringValue{"foo", false}}}},
				err: false}},
		{"no interpolation", args{data: []interface{}{[]interface{}{"foo"}}, manual: true, dataType: options.None},
			want{vals: []Container{{DataType: options.Interface, Values: &interfaceValues{interfaceValue{"foo", false}}}},
				err: false}},
		{"with conversion", args{data: []interface{}{[]interface{}{"foo"}}, manual: true, dataType: options.Bool},
			want{vals: []Container{{DataType: options.Bool, Values: &boolValues{boolValue{true, false}}}},
				err: false}},
		{"fail: unsupported value", args{data: []interface{}{complex64(1)}, manual: false, dataType: options.None},
			want{vals: nil, err: true}},
		{"fail: unsupported datatype", args{data: []interface{}{"foo"}, manual: false, dataType: options.Unsupported},
			want{vals: nil, err: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InterfaceSliceFactory(tt.args.data, tt.args.manual, tt.args.dataType)
			if (err != nil) != tt.want.err {
				t.Errorf("InterfaceSliceFactory() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want.vals) {
				t.Errorf("InterfaceSliceFactory() = %v, want %v", got, tt.want.vals)
			}
		})
	}
}

func TestScalarFactory(t *testing.T) {
	type args struct {
		data interface{}
	}
	type want struct {
		container Container
		err       bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"float32",
			args{float32(1.0)},
			want{Container{Values: &float64Values{float64Value{1.0, false}}, DataType: options.Float64}, false}},
		{"float64",
			args{float64(1.0)},
			want{Container{Values: &float64Values{float64Value{1.0, false}}, DataType: options.Float64}, false}},
		{"int",
			args{int(1)},
			want{Container{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"int8",
			args{int8(1)},
			want{Container{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"int16",
			args{int16(1)},
			want{Container{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"int32",
			args{int32(1)},
			want{Container{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"int64",
			args{int64(1)},
			want{Container{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"uint",
			args{uint(1)},
			want{Container{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"uint8",
			args{uint8(1)},
			want{Container{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"uint16",
			args{uint16(1)},
			want{Container{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"uint32",
			args{uint32(1)},
			want{Container{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"uint64",
			args{uint64(1)},
			want{Container{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"string",
			args{"foo"},
			want{Container{Values: &stringValues{stringValue{"foo", false}}, DataType: options.String}, false}},
		{"string null",
			args{""},
			want{Container{Values: &stringValues{stringValue{"NaN", true}}, DataType: options.String}, false}},
		{"bool",
			args{true},
			want{Container{Values: &boolValues{boolValue{true, false}}, DataType: options.Bool}, false}},
		{"datetime",
			args{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)},
			want{Container{Values: &dateTimeValues{dateTimeValue{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC), false}}, DataType: options.DateTime}, false}},
		{"unsupported",
			args{complex64(1)},
			want{Container{Values: nil, DataType: options.None}, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ScalarFactory(tt.args.data)
			if (err != nil) != tt.want.err {
				t.Errorf("ScalarFactory() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want.container) {
				t.Errorf("ScalarFactory() = %v, want %v", got, tt.want.container)
			}
		})
	}
}

func TestScalarConstructor_NullFloat(t *testing.T) {
	vals, err := ScalarFactory(math.NaN())
	if err != nil {
		t.Errorf("Unable to construct values from null float: %v", err)
	}
	val := vals.Values.Value(0).(float64)
	if !math.IsNaN(val) {
		t.Errorf("Returned %v, want NaN", val)
	}
}

func TestSliceFactory(t *testing.T) {
	type args struct {
		data interface{}
	}
	type want struct {
		container Container
		err       bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{"float32",
			args{[]float32{1.0}},
			want{Container{Values: &float64Values{float64Value{1.0, false}}, DataType: options.Float64}, false}},
		{"float64",
			args{[]float64{1.0}},
			want{Container{Values: &float64Values{float64Value{1.0, false}}, DataType: options.Float64}, false}},
		{"int",
			args{[]int{1}},
			want{Container{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"int8",
			args{[]int8{1}},
			want{Container{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"int16",
			args{[]int16{1}},
			want{Container{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"int32",
			args{[]int32{1}},
			want{Container{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"int64",
			args{[]int64{1}},
			want{Container{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"uint",
			args{[]uint{1}},
			want{Container{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"uint8",
			args{[]uint8{1}},
			want{Container{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"uint16",
			args{[]uint16{1}},
			want{Container{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"uint32",
			args{[]uint32{1}},
			want{Container{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"uint64",
			args{[]uint64{1}},
			want{Container{Values: &int64Values{int64Value{1, false}}, DataType: options.Int64}, false}},
		{"string",
			args{[]string{"foo"}},
			want{Container{Values: &stringValues{stringValue{"foo", false}}, DataType: options.String}, false}},
		{"string null",
			args{[]string{""}},
			want{Container{Values: &stringValues{stringValue{"NaN", true}}, DataType: options.String}, false}},
		{"bool",
			args{[]bool{true}},
			want{Container{Values: &boolValues{boolValue{true, false}}, DataType: options.Bool}, false}},
		{"datetime",
			args{[]time.Time{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)}},
			want{Container{Values: &dateTimeValues{dateTimeValue{time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC), false}}, DataType: options.DateTime}, false}},
		{"unsupported",
			args{[]complex64{complex64(1)}},
			want{Container{Values: nil, DataType: options.None}, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SliceFactory(tt.args.data)
			if (err != nil) != tt.want.err {
				t.Errorf("SliceFactory() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want.container) {
				t.Errorf("SliceFactory() = %v, want %v", got, tt.want.container)
			}
		})
	}
}

func TestMapSplitter(t *testing.T) {
	type args struct {
		data []interface{}
	}
	type want struct {
		isSplit          bool
		extractedData    []interface{}
		extractedColumns []string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{name: "nil", args: args{data: nil},
			want: want{isSplit: false, extractedData: nil, extractedColumns: nil}},
		{"non-splitting map", args{[]interface{}{map[string]bool{"foo": true}}},
			want{isSplit: false, extractedData: nil, extractedColumns: nil}},
		{"pass", args{[]interface{}{map[string]interface{}{"foo": true}}},
			want{isSplit: true, extractedData: []interface{}{true}, extractedColumns: []string{"foo"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isSplit, extractedData, extractedColumns := MapSplitter(tt.args.data)
			if isSplit != tt.want.isSplit {
				t.Errorf("MapSplitter().isSplit = %v, want %v", isSplit, tt.want.isSplit)
			}
			if !reflect.DeepEqual(extractedData, tt.want.extractedData) {
				t.Errorf("MapSplitter().extractedData = %v, want %v", extractedData, tt.want.extractedData)
			}
			if !reflect.DeepEqual(extractedColumns, tt.want.extractedColumns) {
				t.Errorf("MapSplitter().extractedColumns = %v, want %v", extractedColumns, tt.want.extractedColumns)
			}
		})
	}
}

func TestValues_InterpolateString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{name: "int64", args: args{s: "1"},
			want: 1},
		{"float64", args{s: "1.0"},
			1.0},
		{"bool", args{s: "true"},
			true},
		{"datetime", args{s: "Jan 1, 2019"},
			time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"default to string", args{s: "anything else"},
			"anything else"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InterpolateString(tt.args.s)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InterpolateString = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValues_Interpolate(t *testing.T) {
	dt := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	long := make([]interface{}, 100)
	for i := 0; i < 100; i++ {
		long[i] = "foo"
	}
	type args struct {
		data interface{}
	}
	tests := []struct {
		name string
		args args
		want options.DataType
	}{
		{name: "float64", args: args{data: []interface{}{"foo", 1.2, 2.2, float32(3), float64(4)}}, want: options.Float64},
		{name: "int64", args: args{data: []interface{}{"foo", 1, 2, int8(3), uint(4)}}, want: options.Int64},
		{name: "mixed numbers -> float64", args: args{data: []interface{}{"foo", 1.0, float64(2), int8(3), uint(4)}}, want: options.Float64},
		{name: "string", args: args{data: []interface{}{"foo", "bar", "baz", "qux", 4}}, want: options.String},
		{name: "bool", args: args{data: []interface{}{true, false, true, false, "foo"}}, want: options.Bool},
		{name: "dateTime", args: args{data: []interface{}{dt, dt, dt, dt, "foo"}}, want: options.DateTime},
		{name: "none -> interface", args: args{data: []interface{}{1.5, 1, "foo", true, dt}}, want: options.Interface},
		{name: "long string", args: args{data: long}, want: options.String},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Interpolate(tt.args.data)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Interpolate = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceConstructor_NullFloat(t *testing.T) {
	vals, err := SliceFactory([]float64{math.NaN()})
	if err != nil {
		t.Errorf("Unable to construct values from null float: %v", err)
	}
	val := vals.Values.Value(0).(float64)
	if !math.IsNaN(val) {
		t.Errorf("Returned %v, want NaN", val)
	}
}

func TestSliceConstructor_NullFloatInterface(t *testing.T) {
	vals, err := SliceFactory([]interface{}{math.NaN()})
	if err != nil {
		t.Errorf("Unable to construct values from null float: %v", err)
	}
	val := vals.Values.Value(0).(float64)
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

func TestMakeNullRange(t *testing.T) {
	got := MakeNullRange(3)
	want := []interface{}{"", "", ""}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MakeNullRange(): got %v, want %v", got, want)
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

func TestMakeStringRange(t *testing.T) {
	got := MakeStringRange(0, 3)
	want := []string{"0", "1", "2"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MakeIntRange(): got %v, want %v", got, want)
	}
}

func TestContainer_Copy(t *testing.T) {
	container := MustCreateValuesFromInterface([]string{"foo", "bar"})
	got := container.Copy()
	want := Container{&stringValues{stringValue{"foo", false}, stringValue{"bar", false}}, options.String}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Container.Copy(): got %v, want %v", got, want)
	}
}

func TestValues_Transpose(t *testing.T) {
	type args struct {
		data [][]interface{}
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{name: "pass", args: args{[][]interface{}{{1, 2, 3}, {4, 5, 6}}},
			want: []interface{}{[]interface{}{1, 4}, []interface{}{2, 5}, []interface{}{3, 6}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TransposeValues(tt.args.data)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransposeValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
