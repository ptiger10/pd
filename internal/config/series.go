package config

import "github.com/ptiger10/pd/kinds"

// A MiniIndex is an untyped representation of one index level.
// It is used for unpacking client-supplied index data and Constructor metadata.
type MiniIndex struct {
	Data interface{}
	Kind kinds.Kind
	Name string
}

// A ConstructorConfig is an internal type used for configuring the Series.New() function with optional parameters
type ConstructorConfig struct {
	Indices []MiniIndex
	Kind    kinds.Kind
	Name    string
}

// A SelectionConfig is an internal type used for configuring the Series.Select() method with optional parameters
type SelectionConfig struct {
	LevelPositions []int
	LevelNames     []string
	RowPositions   []int
	RowLabels      []string
}
