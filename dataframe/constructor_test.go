package dataframe

import (
	"reflect"
	"testing"
)

// func Test_New(t *testing.T) {
// 	c := Config{Cols: []interface{}{"fooCol", "barCol"}, Index: []string{"foo", "bar", "baz"}}
// 	df, err := New([]interface{}{[]int64{1, 2, 3}, []float64{4, 5, 6}}, c)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	want, _ := New([]interface{}{[]int64{1, 2, 3}, []float64{4, 5, 6}}, c)
// 	if !reflect.DeepEqual(df, want) {
// 		t.Error(err)
// 	}
// }

func TestNew(t *testing.T) {
	type args struct {
		data   []interface{}
		config []Config
	}
	tests := []struct {
		name    string
		args    args
		want    *DataFrame
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.data, tt.args.config...)
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
