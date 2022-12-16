package service

import "testing"

func TestResponse_ToClientString(t *testing.T) {
	t.Run("Request String with Value", func(t *testing.T) {
		response := NewResponse("val", Key("testKey"), Value("veryLongValue"))
		message := response.ClientString()
		if message != "val213veryLongValue" {
			t.Error("Expected message 'val213veryLongValue', got", message)
		}
	})
	t.Run("Request String without Value", func(t *testing.T) {
		response := NewResponse("ack", Key("testKey"), Value("veryLongValue"))
		message := response.ClientString()
		if message != "ack" {
			t.Error("Expected message 'ack', got", message)
		}
	})
}

func TestNewRequest(t *testing.T) {
	request := NewRequest("get", "someKey", "someValue")
	if request.Type != "get" || request.Key != "someKey" || request.Value != "someValue" {
		t.Error("New Request could not be instantiated")
	}
}

func TestRequest_String(t *testing.T) {
	t.Run("Request with Value", func(t *testing.T) {
		request := NewRequest("get", "someKey", "someValue")
		if request.String() != "get[someKey]<someValue>" {
			t.Error("Request string incorrectly generated")
		}
	})
	t.Run("Request without Value", func(t *testing.T) {
		request := NewRequest("get", "someKey", "")
		if request.String() != "get[someKey]" {
			t.Error("New Request could not be instantiated")
		}
	})
}
