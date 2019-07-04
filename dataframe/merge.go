package dataframe

import (
	"reflect"

	"github.com/ptiger10/pd/internal/index"
)

// Join extends the columns, rows, or columns and rows of a dataframe by appending s2 and modifies the DataFrame in place.
// If extending rows, the values within a Values container are converted to []interface if the container datatypes are not the same.
//
// Allowable append values: "rows", "cols", "both"
//
// Allowable method values: "left", "right", "inner", "outer"
func (ip InPlace) Join(append string, method string, df2 *DataFrame) error {
	if ip.df.vals == nil {
		ip.df.replace(df2)
		return nil
	}
	switch append {
	case "rows":
		switch method {
		case "left":

		}
	}
	return nil
}

func matchIndexLevel(df *DataFrame, idx index.Elements) int {
	for i := 0; i < df.Len(); i++ {
		if reflect.DeepEqual(idx, df.index.Elements(i)) {
			return i
		}
	}
	return -1
}

// Join extends the columns, rows, or columns and rows of a dataframe by appending s2 and returns a new DataFrame.
// If extending rows, the values within a Values container are converted to []interface if the container datatypes are not the same.
func (df *DataFrame) Join(append string, method string, df2 *DataFrame) *DataFrame {
	df = df.Copy()
	df.InPlace.Join(append, method, df2)
	return df
}
