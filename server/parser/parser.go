package parser

import "tcpstore/keystore"

type Command struct {
	Task             string // i.e. put/get/delete
	keyArgSizeSize   int
	keyArgSize       int
	Key              string
	valueArgSizeSize int
	valueArgSize     int
	Value            string
}

func ParseCommand(input string) Command {
	// TODO parse the string and return a command Object
	return Command{
		Task:             "",
		keyArgSizeSize:   0,
		keyArgSize:       0,
		Key:              "",
		valueArgSizeSize: 0,
		valueArgSize:     0,
		Value:            "",
	}
}

func (command *Command) ToRequest() keystore.Request {
	return keystore.Request{
		Task:  command.Task,
		Key:   keystore.Key(command.Key),
		Value: keystore.Value(command.Value),
	}
}

func (command *Command) Valid() bool {
	// TODO
	return true
}
