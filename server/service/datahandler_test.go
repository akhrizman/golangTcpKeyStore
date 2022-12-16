package service_test

import (
	"os"
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
	putRequest = service.NewRequest("put", "testKeyPut", "testValuePut")
	getRequest = service.NewRequest("get", "testKeyGet", "testValueGet")
	deleteRequest = service.NewRequest("del", "testKeyDel", "testValueDel")
	badRequest = service.NewRequest("bad", "testKeyDel", "testValueDel")
}

func shutdown() {
	close(dataHandler.PutChannel)
	close(dataHandler.GetChannel)
	close(dataHandler.DeleteChannel)
}

func TestNewDataHandler(t *testing.T) {
	if dataHandler.PutChannel == nil || dataHandler.GetChannel == nil || dataHandler.DeleteChannel == nil {
		t.Error("One or more Request channels not instantiated")
	}
}

func TestDataHandler_QueueRequest(t *testing.T) {
	//t.Run("queue put request", func(t *testing.T) {
	//	defer close(dataHandler.PutChannel)
	//	dataHandler.QueueRequest(putRequest)
	//	//request := <-dataHandler.PutChannel
	//	response := <-putRequest.ResponseChannel
	//	fmt.Println("I'm In the test code")
	//	if !strings.Contains(response.ClientString(), "ack") {
	//		t.Error("Expected to receive put request through channel")
	//	}
	//})
	//t.Run("queue get request", func(t *testing.T) {
	//	dataHandler.QueueRequest(getRequest)
	//	request := <-dataHandler.GetChannel
	//	if request != getRequest {
	//		t.Error("Expected to receive put request through channel")
	//	}
	//})
	//t.Run("queue delete request", func(t *testing.T) {
	//	dataHandler.QueueRequest(deleteRequest)
	//	request := <-dataHandler.DeleteChannel
	//	if request != deleteRequest {
	//		t.Error("Expected to receive put request through channel")
	//	}
	//})
	//t.Run("queue unknown request type", func(t *testing.T) {
	//	dataHandler.QueueRequest(badRequest)
	//	response := <-badRequest.ResponseChannel
	//	exampleResponse := service.NewResponse("", "", "")
	//	if reflect.TypeOf(exampleResponse) == reflect.TypeOf(response) {
	//		t.Error("Expected to receive put request through channel")
	//	}
	//})
}
