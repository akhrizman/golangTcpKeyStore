package service

import (
	"errors"
	"strconv"
)

var (
	ErrParseCommand        = errors.New("failed to parse input string")
	ErrParseNumberOfBytes  = errors.New("unable to parse number of bytes")
	ErrInvalidStringLength = errors.New("invalid string length, may have been truncated")
)

// Command struct to hold task, key, value, and size descriptors
type Command struct {
	Task  string // i.e. put/get/delete
	Key   string
	Value string
}

func NewCommand(task, key, value string) Command {
	return Command{
		Task:  task,
		Key:   key,
		Value: value,
	}
}

// ParseCommand returns a Command struct containing parsed task, key, value, and size descriptors
func ParseCommand(input string) (Command, error) {
	if !ValidCommandLength(input) {
		return Command{}, ErrInvalidStringLength
	}
	task := input[:3]
	key, remainingSegment, errKey := extractArgument(input[3:])
	if errKey != nil {
		return Command{}, ErrParseCommand
	}

	value, _, errValue := extractArgument(remainingSegment)
	if task == "put" && errValue != nil {
		value = ""
	}
	// Unnecessary to set keyArgSizeSize and keyArgSize, but may be useful for debugging
	return NewCommand(task, key, value), nil
}

// ValidCommandLength checks that input meets is of minimum length - example: "put11a"
func ValidCommandLength(input string) bool {
	return len(input) >= 6
}

func (command *Command) ToRequest() Request {
	return Request{
		Task:  command.Task,
		Key:   Key(command.Key),
		Value: Value(command.Value),
	}
}

// Valid checks any additional command requirements unrelated to parsing
func (command *Command) Valid() bool {
	if !TaskAllowed(command.Task) {
		return false
	}
	if command.Key == "" {
		return false
	}
	if command.Task == "put" && command.Value == "" {
		return false
	}
	return true
}

func TaskAllowed(task string) bool {
	return task == "put" || task == "get" || task == "del"
}

// extractArgument accepts a string of format "[argSizeSize(0-9)][argSize(#)][argString][remainingString]"
// and returns the validated string the remaining substring
// argString's size MUST equal it's preceding argSize descriptor
// argSize's size MUST equal it's preceding argSizeSize descriptor
func extractArgument(input string) (string, string, error) {
	if input == "" {
		return "", "", nil
	}

	argSizeSize, err := strconv.Atoi(input[:1])
	if err != nil {
		return "", "", ErrParseNumberOfBytes
	}

	if len(input) < 1+argSizeSize {
		return "", "", ErrInvalidStringLength
	}

	argSize, err := strconv.Atoi(input[1 : 1+argSizeSize])
	if err != nil {
		return "", "", ErrParseNumberOfBytes
	}

	if len(input) < 1+argSizeSize+argSize {
		return "", "", ErrInvalidStringLength
	}

	extractedString := input[1+argSizeSize : 1+argSizeSize+argSize]
	if len(extractedString) != argSize {
		return "", "", ErrInvalidStringLength
	}

	remainingString := input[1+argSizeSize+argSize:]
	return extractedString, remainingString, nil
}
