package kinds

// Kind is a modified reflect.Kind, including time.Time
type Kind int

// Kind convenience options
const (
	Invalid Kind = iota
	Float
	Int
	String
	Bool
	DateTime
	Interface
	Unsupported
)

func (kind Kind) String() string {
	kinds := []string{
		"invalid",
		"float64",
		"int64",
		"string",
		"bool",
		"time.Time",
		"interface",
		"unsupported",
	}

	if kind < Invalid || kind > Unsupported {
		return "unknown"
	}
	return kinds[kind]
}
