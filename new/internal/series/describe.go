package series

import "fmt"

// Len returns the length of the Series (including null values)
//
// Applies to: All
func (s Series) Len() int {
	return s.values.Len()
}

// Describe the key details of the Series
//
// Applies to: All
func (s Series) Describe() {
	fmt.Println(s.values.Describe())
}
