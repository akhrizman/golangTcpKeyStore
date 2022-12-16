package service_test

import (
	"fmt"
	"os"
	"strings"
	"tcpstore/persistence"
	"tcpstore/service"
	"testing"
)

var (
	keyStore      persistence.KeyStore
	dataHandler   service.DataHandler
	putRequest    service.Request
	getRequest    service.Request
	deleteRequest service.Request
	badRequest    service.Request
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	keyStore = persistence.NewKeyStore()
	dataHandler = service.NewDataHandler(&keyStore)
	go dataHandler.RequestMonitor()
	putRequest = service.NewRequest("put", "testKey", "testValue")
	getRequest = service.NewRequest("get", "testKey", "")
	deleteRequest = service.NewRequest("del", "testKey", "")
	badRequest = service.NewRequest("bad", "", "")
}

func shutdown() {
	//close(dataHandler.PutChannel)
	//close(dataHandler.GetChannel)
	//close(dataHandler.DeleteChannel)
}

func TestNewDataHandler(t *testing.T) {
	if dataHandler.PutChannel == nil || dataHandler.GetChannel == nil || dataHandler.DeleteChannel == nil {
		t.Error("One or more Request channels not instantiated")
	}
}

func TestDataHandler_ProcessRequest(t *testing.T) {
	t.Run("queue put request", func(t *testing.T) {
		response := dataHandler.ProcessRequest(putRequest)
		if response.Acknowledgement != "ack" {
			t.Error("Expected 'ack' response to put request")
		}
	})
	t.Run("queue get request", func(t *testing.T) {
		response := dataHandler.ProcessRequest(getRequest)
		if response.Acknowledgement != "val" || response.Value != string(putRequest.Value) {
			fmt.Println(response)
			t.Error("Expected 'val' response to get request")
		}
	})
	t.Run("queue delete request", func(t *testing.T) {
		response := dataHandler.ProcessRequest(deleteRequest)
		if response.Acknowledgement != "ack" {
			t.Error("Expected 'ack' response to delete request")
		}
	})
	t.Run("queue unknown request type", func(t *testing.T) {
		response := dataHandler.ProcessRequest(badRequest)
		if !strings.Contains(response.ClientString(), "err") {
			t.Error("Expected to receive put request through channel")
		}
	})
}

func TestDataHandler_ProcessRequestFailures(t *testing.T) {
	//log.SetOutput(ioutil.Discard)
	//dataHandler.CloseStore()
	//fmt.Println(dataHandler.StoreName())
	//t.Run("queue put request and fail", func(t *testing.T) {
	//	response := dataHandler.ProcessRequest(putRequest)
	//	if response.Acknowledgement != "err" {
	//		t.Error("Expected 'err' response to put request")
	//	}
	//})
	//t.Run("queue get request", func(t *testing.T) {
	//	response := dataHandler.ProcessRequest(getRequest)
	//	if response.Acknowledgement != "val" || response.Value != string(putRequest.Value) {
	//		fmt.Println(response)
	//		t.Error("Expected 'val' response to get request")
	//	}
	//})
	//t.Run("queue delete request", func(t *testing.T) {
	//	response := dataHandler.ProcessRequest(deleteRequest)
	//	if response.Acknowledgement != "ack" {
	//		t.Error("Expected 'ack' response to delete request")
	//	}
	//})
	//t.Run("queue unknown request type", func(t *testing.T) {
	//	response := dataHandler.ProcessRequest(badRequest)
	//	if !strings.Contains(response.ClientString(), "err") {
	//		t.Error("Expected to receive put request through channel")
	//	}
	//})

}
