package service

import (
	"fmt"
	"strconv"
)

type Key string
type Value string

type Request struct {
	Type     string
	Key      Key
	Value    Value
	Response Response
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
	Acknowledgement string // val, del, ack, nil
	Key             string
	Value           string
}

func NewResponse(acknowledgement string, key Key, value Value) Response {
	return Response{
		Acknowledgement: acknowledgement,
		Key:             string(key),
		Value:           string(value),
	}
}

func (resp *Response) ClientString() string {
	switch resp.Acknowledgement {
	case "val":
		return fmt.Sprintf("val%d%d%s", resp.getValueArgSizeSize(), resp.getValueArgSize(), resp.Value)
	default:
		return resp.Acknowledgement
	}
}

func (resp *Response) getValueArgSize() int {
	return len(resp.Value)
}

func (resp *Response) getValueArgSizeSize() int {
	return len(strconv.Itoa(len(resp.Value)))
}
