package service

import (
	"fmt"
	"strconv"
)

type Key string
type Value string

type Request struct {
	Task  string
	Key   Key
	Value Value
}

func NewRequest(task, key, value string) Request {
	return Request{
		Task:  task,
		Key:   Key(key),
		Value: Value(value),
	}
}

func (request *Request) String() string {
	return fmt.Sprintf("%s[%s]<%s>", request.Task, request.Key, request.Value)
}

type Response struct {
	acknowledgement string // val, del, ack, nil
	key             string
	value           string
}

func NewResponse(acknowledgement string, key Key, value Value) Response {
	return Response{
		acknowledgement: acknowledgement,
		key:             string(key),
		value:           string(value),
	}
}

func (response *Response) String() string {
	return fmt.Sprintf("%s[%s]<%s>", response.acknowledgement, response.key, response.value)
}

func (response *Response) ClientString() string {
	switch response.acknowledgement {
	case "val":
		return fmt.Sprintf("val%d%d%s", response.getValueArgSizeSize(), response.getValueArgSize(), response.value)
	default:
		return response.acknowledgement
	}
}

func (response *Response) getValueArgSize() int {
	return len(response.value)
}

func (response *Response) getValueArgSizeSize() int {
	return len(strconv.Itoa(len(response.value)))
}
