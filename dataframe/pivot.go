package dataframe

import (
	"fmt"
	"strings"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/series"
)

func (df *DataFrame) stackVals(col int, g Grouping) []values.Container {
	vals := make([]values.Container, len(g.Groups()))

	// lookup new values
	for n, group := range g.Groups() {
		var d []interface{}
		var counter int
		for i := 0; i < df.Len(); i++ {
			if counter >= len(g.groups[group].Positions) {
				nulls := values.MakeNullRange(df.Len() - counter)
				d = append(d, nulls...)
			} else if g.groups[group].Positions[counter] == i {
				// TODO replace with df.At()
				d = append(d, df.Row(i).Values[0])
				counter++
			} else {
				d = append(d, "")
			}
		}
		container, err := values.InterfaceFactory(d)
		if err != nil {
			fmt.Printf("stackCol(): internal error: %v", err)
		}
		// ducks error because values is assumed to be supported
		container.Values, _ = values.Convert(container.Values, df.dataType())
		container.DataType = df.dataType()
		vals[n] = container
	}
	return vals
}

// stackCol converts a column into a column level and replaces existing column levels
func (df *DataFrame) stackCol(col int) *DataFrame {
	df = df.Copy()
	// preserve original values prior to index modification
	idx := df.index.Copy()
	name := df.cols.Name(col)

	// modify index
	df.InPlace.replaceIndex([]int{col})
	g := df.GroupByIndex()
	cols := index.NewColumns(index.NewColLevel(g.Groups(), name))

	vals := df.stackVals(col, g)
	return newFromComponents(vals, idx, cols, df.Name())
}

// assumes one column and that the level being dropped is not the last one
func (df *DataFrame) transposeIndex(level int) *DataFrame {
	archive := df.Copy()
	archive.index.DropLevel(level)
	archive.index.Subset([]int{0})

	df = df.Copy()
	df.index.SubsetLevels([]int{level})
	df = df.Transpose()
	df.index = archive.index
	return df
}

// stackIndex converts an index level into a column level and replaces existing column levels
func (df *DataFrame) stackIndex(level int) *DataFrame {
	df = df.Copy()
	// modify index
	g := df.GroupByIndex(level)
	cols := index.NewColumns(index.NewColLevel(g.Groups(), df.index.Levels[level].Name))

	vals := df.stackVals(level, g)

	// Remove index to create snapshot of a new index (if level is only level, create default range)
	df, _ = df.ResetIndex(level)

	idx := df.index.Copy()
	return newFromComponents(vals, idx, cols, df.Name())
}

// Pivot transforms data into the desired form and calls aggFunc on the reshaped data.
func (df *DataFrame) Pivot(index int, values int, columns int, aggFunc string) *DataFrame {
	df = df.Copy()
	g := df.GroupBy(index, columns)
	switch aggFunc {
	case "sum":
		df = g.Sum()
	case "mean":
		df = g.Mean()
	case "median":
		df = g.Median()
	case "min":
		df = g.Min()
	case "max":
		df = g.Max()
	case "std":
		df = g.Std()
	}
	df = df.transposeIndex(1)
	return df
}

// Transpose transforms all rows to columns.
func (df *DataFrame) Transpose() *DataFrame {
	ret := newEmptyDataFrame()
	for m := 0; m < df.NumCols(); m++ {
		row := transposeSeries(df.hydrateSeries(m))
		ret.InPlace.appendDataFrameRow(row)
	}
	return ret
}

func transposeSeries(s *series.Series) *DataFrame {
	// Columns
	lvls := make([]index.ColLevel, s.NumLevels())
	cols := index.NewColumns(lvls...)
	container, idx := s.ToInternalComponents()
	for j := 0; j < s.NumLevels(); j++ {
		cols.Levels[j].IsDefault = idx.Levels[j].IsDefault
		cols.Levels[j].DataType = idx.Levels[j].DataType
		cols.Levels[j].Name = idx.Levels[j].Name
		for m := 0; m < s.Len(); m++ {
			elem := idx.Levels[j].Labels.Element(m)
			if !elem.Null {
				cols.Levels[j].Labels = append(cols.Levels[j].Labels, fmt.Sprint(elem.Value))
			} else {
				cols.Levels[j].Labels = append(cols.Levels[j].Labels, "")
			}
		}
	}
	cols.Refresh()

	// Index
	names := strings.Split(s.Name(), values.GetMultiColNameSeparator())
	idxLvls := make([]index.Level, len(names))
	retIdx := index.New(idxLvls...)
	for j := 0; j < len(names); j++ {
		name := names[j]
		idxContainer := parseStringIntoValuesContainer(name)
		retIdx.Levels[j].Labels = idxContainer.Values
		retIdx.Levels[j].DataType = idxContainer.DataType
	}
	retIdx.Refresh()

	// Values
	vals := make([]values.Container, s.Len())
	for m := 0; m < s.Len(); m++ {
		vals[m].Values = container.Values.Subset([]int{m})
		vals[m].DataType = container.DataType
	}

	return newFromComponents(vals, retIdx, cols, "")
}
