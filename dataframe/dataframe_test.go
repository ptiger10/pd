package dataframe

import (
	"fmt"
	"testing"
)

func TestDataTypes(t *testing.T) {
	df, _ := New([]interface{}{"foo"})
	fmt.Println(df.DataTypes())
}

func TestDataType(t *testing.T) {
	df, _ := New([]interface{}{"foo"})
	fmt.Println(df.DataType())
}
