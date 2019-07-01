package options

import (
	"reflect"
	"testing"
)

func TestSettableOptions(t *testing.T) {
	if GetDisplayMaxWidth() != defaultOptions.displayMaxWidth {
		t.Errorf("Default setting not reading for DisplayMaxWidth")
	}
	SetDisplayMaxWidth(15)
	if GetDisplayMaxWidth() != 15 {
		t.Error("Unable to set/get DisplayMaxWidth")
	}

	if GetDisplayFloatPrecision() != defaultOptions.displayFloatPrecision {
		t.Errorf("Default setting not reading for DisplayFloatPrecision")
	}
	SetDisplayFloatPrecision(10)
	if GetDisplayFloatPrecision() != 10 {
		t.Error("Unable to set/get DisplayFloatPrecision")
	}

	if GetDisplayMaxRows() != defaultOptions.displayMaxRows {
		t.Errorf("Default setting not reading for DisplayMaxRows")
	}
	SetDisplayMaxRows(10)
	if GetDisplayMaxRows() != 10 {
		t.Error("Unable to set/get DisplayMaxRows")
	}

	if GetDisplayMaxColumns() != defaultOptions.displayMaxColumns {
		t.Errorf("Default setting not reading for DisplayMaxColumns")
	}
	SetDisplayMaxColumns(10)
	if GetDisplayMaxColumns() != 10 {
		t.Error("Unable to set/get DisplayMaxColumns")
	}

	if GetDisplayRepeatedLabels() != defaultOptions.displayRepeatedLabels {
		t.Errorf("Default setting not reading for DisplayRepeatedLabels")
	}
	SetDisplayRepeatedLabels(true)
	if GetDisplayRepeatedLabels() != true {
		t.Error("Unable to set/get DisplayRepeatedLabels")
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

	if GetLogWarnings() != defaultOptions.logWarnings {
		t.Errorf("Default setting not reading for LogWarnings")
	}
	SetLogWarnings(false)
	if GetLogWarnings() != false {
		t.Error("Unable to set/get LogWarnings")
	}

	if GetAsync() != defaultOptions.async {
		t.Errorf("Default setting not reading for Async")
	}
	SetAsync(false)
	if GetAsync() != false {
		t.Error("Unable to set/get Async")
	}

	RestoreDefaults()
	if GetDisplayMaxWidth() != 35 {
		t.Error("Unable to restore default for DisplayMaxWidth")
	}
	if GetLogWarnings() != true {
		t.Error("Unable to restore default for LogWarnings")
	}

}
