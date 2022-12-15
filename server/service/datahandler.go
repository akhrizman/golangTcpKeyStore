package service

import (
	"tcpstore/logg"
)

const (
	Ack = "ack"
	Val = "val"
	Nil = "nil"
	Err = "err"
	Put = "put"
	Get = "get"
	Del = "del"
)

// Datasource interface
type Datasource interface {
	Name() string
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
	switch request.Type {
	case Put:
		handler.putChannel <- request
	case Get:
		handler.getChannel <- request
	case Del:
		handler.deleteChannel <- request
	default:
		request.ResponseChannel <- NewResponse(Err, "", "")
	}
}

func (handler *DataHandler) RequestMonitor() {
	for {
		select {
		case request := <-handler.putChannel:
			err := handler.store.CreateOrUpdate(request.Key, request.Value)
			if err != nil {
				logg.Error.Printf("For %s - %s", request.String(), err)
			} else {
				request.ResponseChannel <- NewResponse(Ack, request.Key, request.Value)
			}
		case request := <-handler.getChannel:
			value, err := handler.store.Read(request.Key)
			if err != nil {
				logg.Error.Printf("For %s - %s", request.String(), err)
				request.ResponseChannel <- NewResponse(Nil, request.Key, request.Value)
			} else {
				request.ResponseChannel <- NewResponse(Val, request.Key, value)
			}
		case request := <-handler.deleteChannel:
			err := handler.store.Delete(request.Key)
			if err != nil {
				logg.Error.Printf("For %s - %s", request.String(), err)
				request.ResponseChannel <- NewResponse(Err, request.Key, request.Value)
			} else {
				request.ResponseChannel <- NewResponse(Ack, request.Key, request.Value)
			}
		}
	}
}
