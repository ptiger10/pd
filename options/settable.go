package options

var defaultOptions = struct {
	displayMaxWidth                int
	displayValuesWhitespaceBuffer  int
	displayElementWhitespaceBuffer int
	displayIndexWhitespaceBuffer   int
	displayFloatPrecision          int
	displayRepeatedLabels          bool
	displayStringNullFiller        string
	displayTimeFormat              string
	stringNullValues               []string
	logWarnings                    bool
	async                          bool
}{
	displayMaxWidth,
	displayValuesWhitespaceBuffer,
	displayElementWhitespaceBuffer,
	displayIndexWhitespaceBuffer,
	displayFloatPrecision,
	displayRepeatedLabels,
	displayStringNullFiller,
	displayTimeFormat,
	stringNullValues,
	logWarnings,
	async,
}

// RestoreDefaults resets options back to their default setting
func RestoreDefaults() {
	SetDisplayMaxWidth(defaultOptions.displayMaxWidth)
	SetDisplayValuesWhitespaceBuffer(defaultOptions.displayValuesWhitespaceBuffer)
	SetDisplayElementWhitespaceBuffer(defaultOptions.displayElementWhitespaceBuffer)
	SetDisplayIndexWhitespaceBuffer(defaultOptions.displayIndexWhitespaceBuffer)
	SetDisplayFloatPrecision(defaultOptions.displayFloatPrecision)
	SetDisplayRepeatedLabels(defaultOptions.displayRepeatedLabels)
	SetDisplayStringNullFiller(defaultOptions.displayStringNullFiller)
	SetDisplayTimeFormat(defaultOptions.displayTimeFormat)
	SetStringNullValues(defaultOptions.stringNullValues)
	SetLogWarnings(defaultOptions.logWarnings)
	SetAsync(defaultOptions.async)
}

var displayMaxWidth = 35
var displayValuesWhitespaceBuffer = 4
var displayElementWhitespaceBuffer = 1
var displayIndexWhitespaceBuffer = 1
var displayFloatPrecision = 2
var displayRepeatedLabels = false
var displayStringNullFiller = "NaN"
var displayTimeFormat = "1/2/2006T15:04:05"
var stringNullValues = []string{"NaN", "n/a", "N/A", "", "nil"}
var logWarnings = true
var async = true

// SetDisplayMaxWidth sets DisplayMaxWidth to n characters.
// DisplayMaxWidth is an option when printing a Series.
// It is the widest allowable character width for an index label or value.
// If a label is longer than the max, it will be elided at the end.
func SetDisplayMaxWidth(n int) {
	displayMaxWidth = n
}

// GetDisplayMaxWidth returns DisplayMaxWidth.
func GetDisplayMaxWidth() int {
	return displayMaxWidth
}

// SetDisplayValuesWhitespaceBuffer sets DisplayValuesWhitespaceBuffer to n spaces.
// DisplayValuesWhitespaceBuffer is an option when printing a Series or DataFrame.
// It is the number of spaces between the last level of index labels
// and the first collection of values. In a Series, there is only one collection of values.
// In a DataFrame, the first collection of values is the first Series.
func SetDisplayValuesWhitespaceBuffer(n int) {
	displayValuesWhitespaceBuffer = n
}

// GetDisplayValuesWhitespaceBuffer returns DisplayValuesWhitespaceBuffer.
func GetDisplayValuesWhitespaceBuffer() int {
	return displayValuesWhitespaceBuffer
}

// SetDisplayElementWhitespaceBuffer sets DisplayElementWhitespaceBuffer to n spaces.
// DisplayElementWhitespaceBuffer is an option when printing an Element.
// It is the number of spaces between the last level of index labels and the first value.
func SetDisplayElementWhitespaceBuffer(n int) {
	displayElementWhitespaceBuffer = n
}

// GetDisplayElementWhitespaceBuffer returns DisplayElementWhitespaceBuffer.
func GetDisplayElementWhitespaceBuffer() int {
	return displayElementWhitespaceBuffer
}

// SetDisplayIndexWhitespaceBuffer sets DisplayIndexWhitespaceBuffer to n spaces.
// DisplayIndexWhitespaceBuffer is an option when printing a Series.
// It is the number of spaces between index labels. This applies only to a multi-level index.
func SetDisplayIndexWhitespaceBuffer(n int) {
	displayIndexWhitespaceBuffer = n
}

// GetDisplayIndexWhitespaceBuffer returns DisplayIndexWhitespaceBuffer.
func GetDisplayIndexWhitespaceBuffer() int {
	return displayIndexWhitespaceBuffer
}

// SetDisplayFloatPrecision sets DisplayFloatPrecision to n decimal places.
// DisplayFloatPrecision is an option when printing a Series.
// It is the number of decimal points in floating point values and index labels.
func SetDisplayFloatPrecision(n int) {
	displayFloatPrecision = n
}

// GetDisplayFloatPrecision returns DisplayFloatPrecision.
func GetDisplayFloatPrecision() int {
	return displayFloatPrecision
}

// SetDisplayRepeatedLabels sets DisplayRepeatedLabels to boolean.
// DisplayRepeatedLabels is an option when printing a Series.
// If true, all index labels will be shown, like so:
//
// A 0    foo
//
// B 0    bar
//
// C 1    baz
//
// If false, repeated index labels in the same level will be excluded, like so:
//
// A 0    foo
//
// B ... bar
//
// C 1    baz
//
// NB: ellipsis not included in actual printing
func SetDisplayRepeatedLabels(boolean bool) {
	displayRepeatedLabels = boolean
}

// GetDisplayRepeatedLabels returns DisplayRepeatedLabels.
func GetDisplayRepeatedLabels() bool {
	return displayRepeatedLabels
}

// SetDisplayStringNullFiller sets DisplayStringNullFiller to "s".
// DisplayStringNullFiller is an option when printing a Series.
// It is how null string values are represented.
func SetDisplayStringNullFiller(s string) {
	displayStringNullFiller = s
}

// GetDisplayStringNullFiller returns DisplayStringNullFiller.
func GetDisplayStringNullFiller() string {
	return displayStringNullFiller
}

// SetDisplayTimeFormat formats how datetimes are displayed, using the syntax specified in package time.Time.
func SetDisplayTimeFormat(s string) {
	displayTimeFormat = s
}

// GetDisplayTimeFormat returns DisplayTimeFormat.
func GetDisplayTimeFormat() string {
	return displayTimeFormat
}

// SetStringNullValues sets StringNullValues to include only those items contained in nullList.
// StringNullValues is an option when constructing or converting a Series.
// It is the list of string values that are considered null.
func SetStringNullValues(nullList []string) {
	stringNullValues = nullList
}

// GetStringNullValues returns StringNullValues.
func GetStringNullValues() []string {
	return stringNullValues
}

// SetLogWarnings sets LogWarnings to boolean.
// LogWarnings is an option when executing functions within this module.
// If true, non-returned errors are logged to stderr.
// This is relevant for many common exploratory methods, which are often chained together and therefore not designed to return an error value.
func SetLogWarnings(boolean bool) {
	logWarnings = boolean
}

// GetLogWarnings returns LogWarnings.
func GetLogWarnings() bool {
	return logWarnings
}

// SetAsync sets Async to boolean.
// Async is an option for executing certain operations over multiple groups (e.g., math on Groupings or Columns) as goroutines instead of synchronously.
// If true, eligible operations are split into goroutines and merged back together.
func SetAsync(boolean bool) {
	async = boolean
}

// GetAsync returns Async.
func GetAsync() bool {
	return async
}
