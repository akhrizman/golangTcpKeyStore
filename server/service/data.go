package service

import (
	"fmt"
	"strconv"
)

type Key string
type Value string

type Request struct {
	Type            string
	Key             Key
	Value           Value
	ResponseChannel chan Response
}

func NewRequest(task, key, value string) Request {
	return Request{
		Type:  task,
		Key:   Key(key),
		Value: Value(value),
	}
}

func (req *Request) String() string {
	if req.Value == "" {
		return fmt.Sprintf("%s[%s]", req.Type, req.Key)
	} else {
		return fmt.Sprintf("%s[%s]<%s>", req.Type, req.Key, req.Value)
	}
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

func (resp *Response) String() string {
	return fmt.Sprintf("%s[%s]<%s>", resp.acknowledgement, resp.key, resp.value)
}

func (resp *Response) ClientString() string {
	switch resp.acknowledgement {
	case "val":
		return fmt.Sprintf("val%d%d%s", resp.getValueArgSizeSize(), resp.getValueArgSize(), resp.value)
	default:
		return resp.acknowledgement
	}
}

func (resp *Response) getValueArgSize() int {
	return len(resp.value)
}

func (resp *Response) getValueArgSizeSize() int {
	return len(strconv.Itoa(len(resp.value)))
}
