package options

import "strings"

// DataType identifies the type of a data object.
// For most values it is interchangeable with the reflect.Type value, but it supports custom identifiers as well (e.g., DateTime).
type DataType int

// datatype convenience options
const (
	None DataType = iota
	Float64
	Int64
	String
	Bool
	DateTime
	Interface
	PlaceholdervalueType
	Unsupported
)

// DT returns the DataType associated with a string.
func DT(datatype string) DataType {
	datatype = strings.ToLower(datatype)
	switch datatype {
	case "string":
		return String
	case "bool":
		return Bool
	case "datetime":
		return DateTime
	case "interface":
		return Interface
	default:
		if strings.Contains(datatype, "float") {
			return Float64
		}
		if strings.Contains(datatype, "int") {
			return Int64
		}
		return Unsupported
	}
}

func (datatype DataType) String() string {
	datatypes := []string{
		"none",
		"float64",
		"int64",
		"string",
		"bool",
		"dateTime",
		"interface",
		"placeholder",
		"unsupported",
	}

	if datatype < None || datatype > Unsupported {
		return "unknown"
	}
	return datatypes[datatype]
}
