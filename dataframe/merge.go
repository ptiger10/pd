package dataframe

// // Join extends the columns, rows, or columns and rows of a dataframe by appending s2 and modifies the DataFrame in place.
// // If extending rows, the values within a Values container are converted to []interface if the container datatypes are not the same.
// //
// // Allowable append values: "rows", "cols", "both"
// //
// // Allowable method values: "left", "right", "inner", "outer"
// func (ip InPlace) Join(append string, method string, df2 *DataFrame) error {
// 	if ip.df.vals == nil {
// 		ip.df.replace(df2)
// 		return nil
// 	}
// 	switch append {
// 	case "rows":
// 		switch method {
// 		case "left":

// 		}
// 	}
// 	return nil
// }

// assumes equivalent index levels and column positions
func (ip InPlace) appendDataFrameRow(df2 *DataFrame) {
	// Handling empty DataFrame
	if Equal(ip.df, newEmptyDataFrame()) {
		ip.df.replace(df2)
		return
	}

	// Append
	// Index Levels
	for j := 0; j < ip.df.IndexLevels(); j++ {
		ip.df.index.Levels[j].Labels.Append(df2.index.Levels[j].Labels)
		ip.df.index.Levels[j].NeedsRefresh = true
	}
	// Values
	for m := 0; m < ip.df.NumCols(); m++ {
		ip.df.vals[m].Values.Append(df2.vals[m].Values)
	}
	return
}


func (ip InPlace) appendDataFrameColumn(df2 *DataFrame) error {
	// Handling empty DataFrame
	if Equal(ip.df, newEmptyDataFrame()) {
		ip.df.replace(df2)
		return nil
	}

	// Append
	for m := 0; m < df2.NumCols(); m++ {
		// drop errors for now, because input is controlled
		// err := ip.AppendCol(
		ip.AppendCol(
			df2.hydrateSeries(m),
			df2.cols.MultiName(m)...,
		)
		// if err != nil {
		// 	return fmt.Errorf("appendDataFrameColumn(): %v", err)
		// }
	}
	return nil
}

// // Join extends the columns, rows, or columns and rows of a dataframe by appending s2 and returns a new DataFrame.
// // If extending rows, the values within a Values container are converted to []interface if the container datatypes are not the same.
// func (df *DataFrame) Join(append string, method string, df2 *DataFrame) *DataFrame {
// 	df = df.Copy()
// 	df.InPlace.Join(append, method, df2)
// 	return df
// }
