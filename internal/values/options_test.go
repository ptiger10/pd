package values

import "testing"

func TestNonsettableOptions(t *testing.T) {
	if GetDisplayValuesWhitespaceBuffer() != 4 {
		t.Error("Default setting not reading for DisplayValuesWhitespaceBuffer")
	}
	if GetDisplayColumnsWhitespaceBuffer() != 2 {
		t.Error("Default setting not reading for DisplayColumnsWhitespaceBuffer")
	}
	if GetDisplayElementWhitespaceBuffer() != 1 {
		t.Errorf("Default setting not reading for DisplayElementWhitespaceBuffer")
	}
	if GetDisplayIndexWhitespaceBuffer() != 1 {
		t.Errorf("Default setting not reading for DisplayIndexWhitespaceBuffer")
	}
	if GetMultiColNameSeparator() != " | " {
		t.Errorf("Default setting not reading for MultiColNameSeparator")
	}
}
