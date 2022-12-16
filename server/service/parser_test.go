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

//func TestStringExtractor(t *testing.T) {
//	input := "13abc"
//	simplifiedString, remainingString, _ := ExtractArgument(input)
//	if simplifiedString != "abc" {
//		t.Error("Failed to extractArgument string from byte command segment")
//	}
//	if remainingString != "" {
//		t.Error("Expected empty remaining string, got", remainingString)
//	}
//}

func TestExtractArgument(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		result, _, err := ExtractArgument("")
		if result != "" || err != nil {
			t.Error("Expected empty string parsed from empty string")
		}
	})
	t.Run("arg size size not an integer", func(t *testing.T) {
		_, _, err := ExtractArgument("abcdefg")
		if err == nil {
			t.Error("expected error parsing a character to an integer")
		}
	})
	t.Run("argSize not long enough", func(t *testing.T) {
		_, _, err := ExtractArgument("21")
		if err == nil {
			t.Error("Expected Parsing Error")
		}
	})
	t.Run("argument size bytes mismatch with argument", func(t *testing.T) {
		_, _, err := ExtractArgument("21a")
		if err == nil {
			t.Error("arg size specified as 2 bytes, but was 1 byte")
		}
	})

	//t.Run("arg size does not match argument size size descriptor", func(t *testing.T) {
	//	_, _, err := ExtractArgument("12a")
	//	if err == nil {
	//		t.Error("arg size was 2 bytes, but arg was 1 byte")
	//	}
	//})
	//t.Run("remaining string", func(t *testing.T) {
	//	input := "13abcEXTRA"
	//	_, remainingString, _ := ExtractArgument(input)
	//	if remainingString != "EXTRA" {
	//		t.Error("Remaining string expected to be 'EXTRA', got", remainingString)
	//	}
	//})
	//
	//t.Run("missing arg", func(t *testing.T) {
	//	_, _, err := ExtractArgument("11")
	//	if err == nil {
	//		t.Error("Expected Parsing Error")
	//	}
	//})
	//t.Run("arg not long enough", func(t *testing.T) {
	//	_, _, err := ExtractArgument("13a")
	//	if err == nil {
	//		t.Error("Expected Parsing Error")
	//	}
	//})

}

func TestParseMessage(t *testing.T) {
	inputString := "put13key212stored value"
	command, _ := ParseMessage(inputString)
	t.Run("parse type segment", func(t *testing.T) {
		if command.task != "put" {
			t.Error("error getting task expected 'put', got", command.task)
		}
	})
	t.Run("parse key segment", func(t *testing.T) {
		if command.key != "key" {
			t.Error("error getting key expected 'key', got", command.key)
		}
	})
	t.Run("parse value segment", func(t *testing.T) {
		if command.value != "stored value" {
			t.Error("error getting task expected 'stored value', got", command.value)
		}
	})
}

func TestParseMessage_Fail(t *testing.T) {
	t.Run("short message", func(t *testing.T) {
		shortMessage := "put1a"
		_, err := ParseMessage(shortMessage)
		if err == nil {
			t.Error("Expected validation error due to invalid string length")
		}
	})
	t.Run("unable to parse", func(t *testing.T) {
		badMessage := "put232adg"
		_, err := ParseMessage(badMessage)
		if err == nil {
			t.Error("Expected parsing error")
		}
	})
}

func TestCommand_Valid(t *testing.T) {
	t.Run("missing key", func(t *testing.T) {
		command := NewCommand("someType", "", "someValue")
		if command.Valid() {
			t.Error("Command with No key should be invalid")
		}
	})
	t.Run("missing value when put type", func(t *testing.T) {
		command := NewCommand("put", "someKey", "")
		if command.Valid() {
			t.Error("Command with no value for put type should be invalid")
		}
	})
	t.Run("valid command", func(t *testing.T) {
		command := NewCommand("get", "someKey", "someValue")
		if !command.Valid() {
			t.Error("Valid command found to be invalid")
		}
	})
}

func TestTrimMessage(t *testing.T) {
	t.Run("long message", func(t *testing.T) {
		longMessage := "This message is longer than 30 characters"
		result := TrimMessage(longMessage)
		expectedMessage := longMessage[:MaximumLengthForLoggingCommand] + "..."
		if result != expectedMessage {
			t.Errorf("error getting key expected '%s...', got '%s'", expectedMessage, result)
		}
	})
	t.Run("short message", func(t *testing.T) {
		shortMessage := "This message is short"
		result := TrimMessage(shortMessage)
		if result != shortMessage {
			t.Errorf("error getting key expected '%s', got '%s'", shortMessage, result)
		}
	})
}

// True TDD for Variable Length Stretch Goal
func TestNewCommand_StoresVariableLengthArg(t *testing.T) {
	expectedCommand := NewCommand("get", "someKey", "", 2)
	if expectedCommand.VariableLength != 2 {
		t.Error
	}
}
