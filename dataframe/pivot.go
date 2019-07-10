package dataframe

import (
	"fmt"
	"strings"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
	"github.com/ptiger10/pd/options"
	"github.com/ptiger10/pd/series"
)

// values:
// make a [][]interface{} valsMatrix for rows x cols
// # rows: unique non-stacked labels
// # cols = unique stacked labels * number of columns
// isolate first value of the stacked label within each non-stacked label
// transpose to []interface and feed into interface factory to create []Values.Container
func (df *DataFrame) stack(level int) (newIdxPositions []int, valsMatrix [][]interface{}, newColLvl []string) {
	var unstackedIndexLevels []int
	for j := 0; j < df.IndexLevels(); j++ {
		if j != level {
			unstackedIndexLevels = append(unstackedIndexLevels, j)
		}
	}
	g := df.GroupByIndex(unstackedIndexLevels...)

	labelsToStack := df.Index.unique(level)
	numRows := g.Len()
	numCols := len(labelsToStack) * df.NumCols()
	valsMatrix = make([][]interface{}, numRows)
	for i := 0; i < numRows; i++ {
		valsMatrix[i] = make([]interface{}, numCols)
	}

	// only extend the labels for the columns-to-be-stacked once
	extendColLevel := true
	for i, group := range g.Groups() {
		newIdxPositions = append(newIdxPositions, g.groups[group].Positions[0])
		rows, _ := df.SubsetRows(g.groups[group].Positions)
		for labelOffset, label := range labelsToStack {
			// log warnings disabled because frequently a label will not exist in an index
			archive := options.GetLogWarnings()
			options.SetLogWarnings(false)
			row := rows.SelectLabels([]string{label}, level)
			options.SetLogWarnings(archive)
			// log warnings restored
			for m := 0; m < df.NumCols(); m++ {
				if len(row) > 0 {
					valsMatrix[i][m+labelOffset*df.NumCols()] = rows.vals[m].Values.Element(row[0]).Value
				}
				if extendColLevel {
					newColLvl = append(newColLvl, label)
				}
			}
		}
		extendColLevel = false
	}
	return newIdxPositions, valsMatrix, newColLvl
}

func (df *DataFrame) stackIndex(level int) *DataFrame {
	newIdxPositions, valsMatrix, newColLevel := df.stack(level)
	transposedVals := values.TransposeValues(valsMatrix)
	var containers []values.Container
	for i := 0; i < len(transposedVals); i++ {
		container := values.MustCreateValuesFromInterface(transposedVals[i])
		containers = append(containers, container)
	}

	idx := df.index.Copy()
	idx.Subset(newIdxPositions)
	idx.DropLevel(level)

	cols := df.cols.Copy()
	for j := 0; j < df.ColLevels(); j++ {
		// duplicate each level enough times that it is same length as new column level
		cols.Levels[j].Duplicate((len(newColLevel) / df.NumCols()) - 1)
	}

	// ducks error because input is controlled
	cols.InsertLevel(0, newColLevel, df.index.Levels[level].Name)

	ret := newFromComponents(containers, idx, cols, df.Name())
	if df.dataType() != options.Interface {
		ret.InPlace.Convert(df.dataType().String())
	}
	return ret
}

// Pivot transforms data into the desired form and calls aggFunc on the reshaped data.
func (df *DataFrame) Pivot(index int, values int, columns int, aggFunc string) (*DataFrame, error) {
	df = df.Copy()
	df.InPlace.SubsetColumns([]int{index, columns, values})
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
	default:
		return newEmptyDataFrame(), fmt.Errorf("df.Pivot(): aggFunc (%v) does not exist", aggFunc)
	}
	df = df.stackIndex(1)
	df.Columns.DropLevel(1)
	return df, nil
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
			// TODO: test null value
			// if !elem.Null {
			// 	cols.Levels[j].Labels = append(cols.Levels[j].Labels, fmt.Sprint(elem.Value))
			// } else {
			// 	cols.Levels[j].Labels = append(cols.Levels[j].Labels, "")
			// }
			cols.Levels[j].Labels = append(cols.Levels[j].Labels, fmt.Sprint(elem.Value))
		}
	}
	cols.Refresh()

	// Index
	names := strings.Split(s.Name(), values.GetMultiColNameSeparator())
	idxLvls := make([]index.Level, len(names))
	retIdx := index.New(idxLvls...)
	for j := 0; j < len(names); j++ {
		name := names[j]
		// ducks error because type is known to be supported
		retIdx.Levels[j], _ = index.NewLevel(values.InterpolateString(name), "")
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
