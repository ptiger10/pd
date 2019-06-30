package values

var displayValuesWhitespaceBuffer = 4
var displayElementWhitespaceBuffer = 1
var displayIndexWhitespaceBuffer = 1

// GetDisplayValuesWhitespaceBuffer returns DisplayValuesWhitespaceBuffer.
// DisplayValuesWhitespaceBuffer is an option when printing a Series or DataFrame.
// It is the number of spaces between the last level of index labels
// and the first collection of values. In a Series, there is only one collection of values.
// In a DataFrame, the first collection of values is the first Series.
//
// Default buffer: 4 spaces
func GetDisplayValuesWhitespaceBuffer() int {
	return displayValuesWhitespaceBuffer
}

// GetDisplayElementWhitespaceBuffer returns DisplayElementWhitespaceBuffer.
// DisplayElementWhitespaceBuffer is an option when printing an Element.
// It is the number of spaces between the last level of index labels and the first value.
//
// // Default buffer: 1 space
func GetDisplayElementWhitespaceBuffer() int {
	return displayElementWhitespaceBuffer
}

// GetDisplayIndexWhitespaceBuffer returns DisplayIndexWhitespaceBuffer.
// DisplayIndexWhitespaceBuffer is an option when printing a Series.
// It is the number of spaces between index labels. This applies only to a multi-level index.
//
// Default buffer: 1 space
func GetDisplayIndexWhitespaceBuffer() int {
	return displayIndexWhitespaceBuffer
}
