package options

import (
	"reflect"
	"testing"
)

func TestSetting(t *testing.T) {
	if GetDisplayIndexMaxWidth() != defaultOptions.displayIndexMaxWidth {
		t.Errorf("Default setting not reading for DisplayIndexMaxWidth")
	}
	SetDisplayIndexMaxWidth(15)
	if GetDisplayIndexMaxWidth() != 15 {
		t.Error("Unable to set/get DisplayIndexMaxWidth")
	}

	if GetDisplayIndexWhitespaceBuffer() != defaultOptions.displayIndexWhitespaceBuffer {
		t.Errorf("Default setting not reading for DisplayIndexWhitespaceBuffer")
	}
	SetDisplayIndexWhitespaceBuffer(5)
	if GetDisplayIndexWhitespaceBuffer() != 5 {
		t.Error("Unable to set/get DisplayIndexWhitespaceBuffer")
	}

	if GetDisplayValuesWhitespaceBuffer() != defaultOptions.displayValuesWhitespaceBuffer {
		t.Errorf("Default setting not reading for DisplayValuesWhitespaceBuffer")
	}
	SetDisplayValuesWhitespaceBuffer(10)
	if GetDisplayValuesWhitespaceBuffer() != 10 {
		t.Error("Unable to set/get DisplayValuesWhitespaceBuffer")
	}

	if GetDisplayElementWhitespaceBuffer() != defaultOptions.displayElementWhitespaceBuffer {
		t.Errorf("Default setting not reading for DisplayElementWhitespaceBuffer")
	}
	SetDisplayElementWhitespaceBuffer(10)
	if GetDisplayElementWhitespaceBuffer() != 10 {
		t.Error("Unable to set/get DisplayElementWhitespaceBuffer")
	}

	if GetDisplayFloatPrecision() != defaultOptions.displayFloatPrecision {
		t.Errorf("Default setting not reading for DisplayFloatPrecision")
	}
	SetDisplayFloatPrecision(10)
	if GetDisplayFloatPrecision() != 10 {
		t.Error("Unable to set/get DisplayFloatPrecision")
	}

	if GetDisplayRepeatedIndexLabels() != defaultOptions.displayRepeatedIndexLabels {
		t.Errorf("Default setting not reading for DisplayRepeatedIndexLabels")
	}
	SetDisplayRepeatedIndexLabels(true)
	if GetDisplayRepeatedIndexLabels() != true {
		t.Error("Unable to set/get DisplayRepeatedIndexLabels")
	}

	if GetDisplayStringNullFiller() != defaultOptions.displayStringNullFiller {
		t.Errorf("Default setting not reading for DisplayStringNullFiller")
	}
	SetDisplayStringNullFiller("Nothing")
	if GetDisplayStringNullFiller() != "Nothing" {
		t.Error("Unable to set/get DisplayStringNullFiller")
	}

	if GetDisplayTimeFormat() != defaultOptions.displayTimeFormat {
		t.Errorf("Default setting not reading for DisplayTimeFormat")
	}
	SetDisplayTimeFormat("2006")
	if GetDisplayTimeFormat() != "2006" {
		t.Error("Unable to set/get DisplayTimeFormat")
	}

	SetStringNullValues([]string{"Nada", "Nothing"})
	if !reflect.DeepEqual(GetStringNullValues(), []string{"Nada", "Nothing"}) {
		t.Error("Unable to set/get StringNullValues")
	}

	SetLogWarnings(false)
	if GetLogWarnings() != false {
		t.Error("Unable to set/get LogWarnings")
	}
	RestoreDefaults()
	if GetDisplayIndexMaxWidth() != 25 {
		t.Error("Unable to restore default for DisplayIndexMaxWidth")
	}

}
