package service

import "testing"

func TestValidCommandLength(t *testing.T) {
	if ValidCommandLength("123456") == false {
		t.Error("Expected valid command length")
	}
	if ValidCommandLength("12345") == true {
		t.Error("Expected invalid command length")
	}
}

func TestCommand_ToRequest(t *testing.T) {
	command := NewCommand("get", "testKey", "testValue")
	request := command.AsRequest()
	if request.Type != "get" {
		t.Error("expected task 'get', got", request.Type)
	}
	if request.Key != "testKey" {
		t.Error("expected task 'get', got", request.Key)
	}
	if request.Value != "testValue" {
		t.Error("expected task 'get', got", request.Value)
	}
}

func TestStringExtractor(t *testing.T) {
	input := "13abc"
	simplifiedString, remainingString, _ := extractArgument(input)
	if simplifiedString != "abc" {
		t.Error("Failed to extractArgument string from byte command segment")
	}
	if remainingString != "" {
		t.Error("Expected empty remaining string, got", remainingString)
	}
}

func TestStringExtractorWhenRemainingString(t *testing.T) {
	input := "13abcEXTRA"
	_, remainingString, _ := extractArgument(input)
	if remainingString != "EXTRA" {
		t.Error("Remaining string expected to be 'EXTRA', got", remainingString)
	}
}

func TestStringExtractorParseEmptyString(t *testing.T) {
	result, _, err := extractArgument("")
	if result != "" || err != nil {
		t.Error("Expected empty string parsed from empty string")
	}
}

func TestStringExtractorArgMissing(t *testing.T) {
	_, _, err := extractArgument("11")
	if err == nil {
		t.Error("Expected Parsing Error")
	}
}

func TestStringExtractorArgNotLongEnough(t *testing.T) {
	_, _, err := extractArgument("13a")
	if err == nil {
		t.Error("Expected Parsing Error")
	}
}

func TestStringExtractorArgSizeNotLongEnough(t *testing.T) {
	_, _, err := extractArgument("21")
	if err == nil {
		t.Error("Expected Parsing Error")
	}
}

func TestParseCommand(t *testing.T) {
	inputString := "put13key212stored value"
	command, _ := ParseMessage(inputString)
	if command.task != "put" {
		t.Error("error getting task expected 'put', got", command.task)
	}
	if command.key != "key" {
		t.Error("error getting key expected 'key', got", command.key)
	}
	if command.value != "stored value" {
		t.Error("error getting task expected 'stored value', got", command.value)
	}
}
