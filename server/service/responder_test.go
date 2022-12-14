package service

import "testing"

func TestResponse_ToClientString(t *testing.T) {
	response := NewResponse("val", Key("testKey"), Value("veryLongValue"), 1)
	message := response.ClientString()
	if message != "val213veryLongValue" {
		t.Error("Expected message 'val213veryLongValue', got", message)
	}
}

func TestResponse_ToClientStringAcknowledged(t *testing.T) {
	response := NewResponse("ack", Key("testKey"), Value("veryLongValue"), 1)
	message := response.ClientString()
	if message != "ack" {
		t.Error("Expected message 'ack', got", message)
	}
}

func TestResponse_ToClientStringFailed(t *testing.T) {
	response := NewResponse("nil", Key("testKey"), Value("veryLongValue"), 1)
	message := response.ClientString()
	if message != "nil" {
		t.Error("Expected message 'nil', got", message)
	}
}
