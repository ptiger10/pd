package dataframe

import (
	"fmt"

	"github.com/ptiger10/pd/internal/index"
	"github.com/ptiger10/pd/internal/values"
)

// stackCol converts a column into a column level and replaces existing column levels
func (df *DataFrame) stackCol(col int) *DataFrame {
	df = df.Copy()
	// preserve original values prior to index modification
	idx := df.index.Copy()
	name := df.cols.Name(col)
	datatype := df.cols.Levels[col].DataType

	// modify index
	df.InPlace.replaceIndex([]int{col})
	g := df.GroupByIndex()
	fmt.Println(g.Groups())
	cols := index.NewColumns(index.NewColLevel(g.Groups(), name))
	vals := make([]values.Container, cols.Len())

	// lookup new values
	for n, group := range g.Groups() {
		var d []interface{}
		var counter int
		for i := 0; i < df.Len(); i++ {
			if counter >= len(g.groups[group].Positions) {
				nulls := values.MakeNullRange(df.Len() - counter)
				d = append(d, nulls...)
			} else if g.groups[group].Positions[counter] == i {
				d = append(d, df.Row(i).Values[col])
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
		container.Values, _ = values.Convert(container.Values, datatype)
		container.DataType = datatype
		vals[n] = container
	}
	return newFromComponents(vals, idx, cols, df.Name())

}
