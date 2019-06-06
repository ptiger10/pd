package series

import (
	"github.com/ptiger10/pd/kinds"
)

// Index contains index selection and conversion
type Index struct {
	s  *Series
	To To
}

// To contains conversion methods
type To struct {
	s   *Series
	idx bool
}

// Float converts Series values to float64.
func (t To) Float() Series {
	if t.idx {
		t.s.index = t.s.index.Copy()
		t.s.index.Levels[0] = t.s.index.Levels[0].ToFloat()
	} else {
		t.s.values = t.s.values.ToFloat()
		t.s.kind = kinds.Float
	}
	return *t.s
}

// Int converts Series values to int64.
func (t To) Int() Series {
	if t.idx {
		t.s.index = t.s.index.Copy()
		t.s.index.Levels[0] = t.s.index.Levels[0].ToInt()
	} else {
		t.s.values = t.s.values.ToInt()
		t.s.kind = kinds.Int
	}
	return *t.s
}

// String converts Series values to string.
func (t To) String() Series {
	if t.idx {
		t.s.index = t.s.index.Copy()
		t.s.index.Levels[0] = t.s.index.Levels[0].ToString()
	} else {
		t.s.values = t.s.values.ToString()
		t.s.kind = kinds.String
	}
	return *t.s
}

// Bool converts Series values to bool.
func (t To) Bool() Series {
	if t.idx {
		t.s.index = t.s.index.Copy()
		t.s.index.Levels[0] = t.s.index.Levels[0].ToBool()
	} else {
		t.s.values = t.s.values.ToBool()
		t.s.kind = kinds.Bool
	}
	return *t.s
}

// DateTime converts Series values to time.Time.
func (t To) DateTime() Series {
	if t.idx {
		t.s.index = t.s.index.Copy()
		t.s.index.Levels[0] = t.s.index.Levels[0].ToDateTime()
	} else {
		t.s.values = t.s.values.ToDateTime()
		t.s.kind = kinds.DateTime
	}
	return *t.s
}

// Interface converts Series values to interface.
func (t To) Interface() Series {
	if t.idx {
		t.s.index = t.s.index.Copy()
		t.s.index.Levels[0] = t.s.index.Levels[0].ToInterface()
	} else {
		t.s.values = t.s.values.ToInterface()
		t.s.kind = kinds.Interface
	}
	return *t.s
}
