package opt

import (
	"github.com/ptiger10/pd/datatypes"
	"github.com/ptiger10/pd/internal/config"
)

// [START Constructor options]

// A ConstructorOption is an optional parameter in the Series constructor.
type ConstructorOption func(*config.ConstructorConfig)

// DataType will convert either values or an index level to the specified kind
func DataType(datatype datatypes.DataType) ConstructorOption {
	return func(c *config.ConstructorConfig) {
		c.DataType = datatype
	}
}

// Name will name either values or an index level
func Name(n string) ConstructorOption {
	return func(c *config.ConstructorConfig) {
		c.Name = n
	}
}

// [END Constructor options]

// [START Selection options]

// A SelectionOption is an optional parameter in a Series selector.
type SelectionOption func(*config.SelectionConfig)

// ByLevels selects one or more index levels by their integer positions
func ByLevels(positions []int) SelectionOption {
	return func(c *config.SelectionConfig) {
		c.LevelPositions = positions
	}
}

// ByLevelNames selects one or more index levels by their names
func ByLevelNames(names []string) SelectionOption {
	return func(c *config.SelectionConfig) {
		c.LevelNames = names
	}
}

// ByRows selects one or more rows by their integer positions
func ByRows(positions []int) SelectionOption {
	return func(c *config.SelectionConfig) {
		c.RowPositions = positions
	}
}

// ByLabels selects one or more rows by their stringified index labels
func ByLabels(labels []string) SelectionOption {
	return func(c *config.SelectionConfig) {
		c.RowLabels = labels
	}
}

// [END Selection options]
