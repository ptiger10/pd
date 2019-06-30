package options

var defaultOptions = struct {
	displayMaxWidth         int
	displayFloatPrecision   int
	displayRepeatedLabels   bool
	displayStringNullFiller string
	displayTimeFormat       string
	stringNullValues        []string
	logWarnings             bool
	async                   bool
}{
	displayMaxWidth,
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
	SetDisplayFloatPrecision(defaultOptions.displayFloatPrecision)
	SetDisplayRepeatedLabels(defaultOptions.displayRepeatedLabels)
	SetDisplayStringNullFiller(defaultOptions.displayStringNullFiller)
	SetDisplayTimeFormat(defaultOptions.displayTimeFormat)
	SetStringNullValues(defaultOptions.stringNullValues)
	SetLogWarnings(defaultOptions.logWarnings)
	SetAsync(defaultOptions.async)
}

var displayMaxWidth = 35
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
//
// Default width: 35 characters
func SetDisplayMaxWidth(n int) {
	displayMaxWidth = n
}

// GetDisplayMaxWidth returns DisplayMaxWidth.
func GetDisplayMaxWidth() int {
	return displayMaxWidth
}

// SetDisplayFloatPrecision sets DisplayFloatPrecision to n decimal places.
// DisplayFloatPrecision is an option when printing a Series.
// It is the number of decimal points in floating point values and index labels.
//
// Default precision: 2 decimal points
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
//
// Default: false
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
//
// Default: "NaN"
func SetDisplayStringNullFiller(s string) {
	displayStringNullFiller = s
}

// GetDisplayStringNullFiller returns DisplayStringNullFiller.
func GetDisplayStringNullFiller() string {
	return displayStringNullFiller
}

// SetDisplayTimeFormat formats how datetimes are displayed, using the syntax specified in package time.Time.
//
// Default: "1/2/2006T15:04:05"
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
//
// default: []string{"NaN", "n/a", "N/A", "", "nil"}
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
//
// default: true
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
//
// default: true
func SetAsync(boolean bool) {
	async = boolean
}

// GetAsync returns Async.
func GetAsync() bool {
	return async
}
