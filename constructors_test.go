package pd

import (
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
			args{"foo", []Config{Config{Name: "bar"}}},
			want{series.MustNew("foo", series.Config{Name: "bar"}), false}},
		{"config with df field",
			args{"foo", []Config{Config{Name: "bar", Cols: []interface{}{"baz"}}}},
			want{series.MustNew("foo", series.Config{Name: "bar"}), false}},
		{"fail: multiple configs",
			args{"foo", []Config{Config{Name: "bar"}, Config{Name: "baz"}}},
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
			args{[]interface{}{"foo"}, []Config{Config{Name: "bar"}}},
			want{dataframe.MustNew([]interface{}{"foo"}, dataframe.Config{Name: "bar"}), false}},
		{"fail: multiple configs",
			args{[]interface{}{"foo"}, []Config{Config{Name: "bar"}, Config{Name: "baz"}}},
			want{dataframe.MustNew(nil), true}},
		{"fail: unsupported value",
			args{[]interface{}{complex64(1)}, nil},
			want{dataframe.MustNew(nil), true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DataFrame(tt.args.data, tt.args.config...)
			if (err != nil) != tt.want.err {
				t.Errorf("Series():  error = %v, want %v", err, tt.want.err)
			}
			if !dataframe.Equal(got, tt.want.df) {
				t.Errorf("Series() got %v, want %v", got, tt.want.df)
			}
		})
	}
}
