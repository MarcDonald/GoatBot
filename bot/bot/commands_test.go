package bot

import "testing"

func TestCheckParametersForReservedKeyword_NoReservedKeywords(t *testing.T) {
	err := checkParametersForReservedKeyword(InvokableCommand{Parameters: []CommandParameter{{Name: "something"}}})
	if err != nil {
		t.Error("Test Failed: Expected error to be nil but was: " + err.Error())
	}
}

func TestCheckParametersForReservedKeyword_Username(t *testing.T) {
	err := checkParametersForReservedKeyword(InvokableCommand{Parameters: []CommandParameter{{Name: "username"}}})
	if err == nil {
		t.Error("Test Failed: Expected error to not be nil")
	} else {
		if err.Error() != "reserved keyword 'username' cannot be used as a parameter name" {
			t.Error("Test Failed: Expected error to be 'reserved keyword 'username' cannot be used as a parameter name' but was: " + err.Error())
		}
	}
}
