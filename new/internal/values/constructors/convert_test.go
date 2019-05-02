package constructors

import (
	"math"
	"testing"
	"time"

	"github.com/ptiger10/pd/new/internal/values"
	"github.com/ptiger10/pd/new/kinds"
)

func TestConvert_string(t *testing.T) {
	var v values.Values
	var err error
	var tests = []struct {
		want  string
		input values.Values
	}{
		{"NaN", SliceFloat([]float64{math.NaN()})},
		{"100", SliceInt([]int{100})},
		{"100", SliceString([]string{"100"})},
		{"true", SliceBool([]bool{true})},
		{
			"2019-05-01 00:00:00 +0000 UTC",
			SliceDateTime([]time.Time{time.Date(2019, 5, 1, 0, 0, 0, 0, time.UTC)}),
		},
		{"100", SliceInterface([]interface{}{100})},
	}

	for _, test := range tests {
		v, err = convert(test.input, kinds.String)
		if err != nil {
			t.Errorf("Unable to convert to string: %v", err)
		}
		vals := v.All()
		got := vals[0].(string)
		if got != test.want {
			t.Errorf("String conversion returned %v, want %v", got, test.want)
		}
	}

}
