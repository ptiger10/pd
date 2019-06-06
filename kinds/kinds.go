package kinds

// Kind is a modified reflect.Kind, including time.Time
type Kind int

// Kind convenience options
const (
	None Kind = iota
	Float64
	Int64
	String
	Bool
	DateTime
	Interface
	PlaceholdervalueType
	Unsupported
)

func (kind Kind) String() string {
	kinds := []string{
		"none",
		"float64",
		"int64",
		"string",
		"bool",
		"time.Time",
		"interface",
		"placeholder",
		"unsupported",
	}

	if kind < None || kind > Unsupported {
		return "unknown"
	}
	return kinds[kind]
}
