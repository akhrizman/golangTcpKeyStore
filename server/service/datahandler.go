package service

import "fmt"

// Datasource interface
type Datasource interface {
	CreateOrUpdate(key Key, value Value) error
	Read(key Key) (Value, error)
	Delete(key Key) error
	IsClosed() bool
}

type DataHandler struct {
	store         Datasource
	putChannel    chan Request
	getChannel    chan Request
	deleteChannel chan Request
}

func NewDataHandler(datasource Datasource) DataHandler {
	return DataHandler{
		store:         datasource,
		putChannel:    make(chan Request),
		getChannel:    make(chan Request),
		deleteChannel: make(chan Request),
	}
}

func (handler *DataHandler) QueueRequest(request Request) {
	switch request.Task {
	case "put":
		fmt.Println("Sending Request to put channel: ", request)
		handler.putChannel <- request
	case "get":
		fmt.Println("Sending Request to get channel: ", request)
		handler.getChannel <- request
	case "del":
		fmt.Println("Sending Request to delete channel: ", request)
		handler.deleteChannel <- request
	default:
		request.ResponseChannel <- NewResponse("err", "", "")
	}
}

func (handler *DataHandler) RequestMonitor() {
	for {
		select {
		case request := <-handler.putChannel:
			fmt.Println("Received request from put channel: ", request)
			err := handler.store.CreateOrUpdate(request.Key, request.Value)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Put Response generated and passed to response channel: ", request)
				request.ResponseChannel <- NewResponse("ack", request.Key, request.Value)
			}
		case request := <-handler.getChannel:
			fmt.Println("Received request from get channel: ", request)
			value, err := handler.store.Read(request.Key)
			if err != nil {
				fmt.Println(err)
				request.ResponseChannel <- NewResponse("nil", request.Key, request.Value)
			} else {
				fmt.Println("Get Response generated and passed to response channel: ", request)
				request.ResponseChannel <- NewResponse("val", request.Key, value)
			}
		case request := <-handler.deleteChannel:
			fmt.Println("Received request from delete channel: ", request)
			handler.store.Delete(request.Key)
			fmt.Println("Delete Response generated and passed to response channel: ", request)
			request.ResponseChannel <- NewResponse("ack", request.Key, request.Value)
		}
	}
}
