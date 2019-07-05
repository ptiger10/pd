package dataframe

import "testing"

func TestDataFrame_Pivot_stack(t *testing.T) {
	type args struct {
		col int
	}
	tests := []struct {
		name  string
		input *DataFrame
		args  args
		want  *DataFrame
	}{
		{name: "one group", input: MustNew([]interface{}{[]int64{1, 1, 1}, []string{"foo", "bar", "baz"}}, Config{Col: []string{"A", "B"}}),
			args: args{col: 0}, want: MustNew([]interface{}{[]string{"foo", "bar", "baz"}}, Config{Col: []string{"1"}, ColName: "A"})},
		{name: "two groups", input: MustNew([]interface{}{[]int64{1, 1, 2}, []string{"foo", "bar", "baz"}}, Config{Col: []string{"A", "B"}}),
			args: args{col: 0},
			want: MustNew([]interface{}{[]string{"foo", "bar", ""}, []string{"", "", "baz"}}, Config{Col: []string{"1", "2"}, ColName: "A"})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.stackCol(tt.args.col)
			if !Equal(got, tt.want) {
				t.Errorf("DataFrame.stack() = %v, want %v", got, tt.want)
			}
		})
	}
}
