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
	PutChannel    chan Request
	GetChannel    chan Request
	DeleteChannel chan Request
}

func NewDataHandler(datasource Datasource) DataHandler {
	return DataHandler{
		store:         datasource,
		PutChannel:    make(chan Request),
		GetChannel:    make(chan Request),
		DeleteChannel: make(chan Request),
	}
}

func (handler *DataHandler) QueueRequest(request Request) {
	switch request.Type {
	case Put:
		handler.PutChannel <- request
	case Get:
		handler.GetChannel <- request
	case Del:
		handler.DeleteChannel <- request
	default:
		request.ResponseChannel <- NewResponse(Err, "", "")
	}
}

func (handler *DataHandler) RequestMonitor() {
	for {
		select {
		case request := <-handler.PutChannel:
			err := handler.store.CreateOrUpdate(request.Key, request.Value)
			if err != nil {
				logg.Error.Printf("For %s - %s", request.String(), err)
			} else {
				request.ResponseChannel <- NewResponse(Ack, request.Key, request.Value)
			}
		case request := <-handler.GetChannel:
			value, err := handler.store.Read(request.Key)
			if err != nil {
				logg.Error.Printf("For %s - %s", request.String(), err)
				request.ResponseChannel <- NewResponse(Nil, request.Key, request.Value)
			} else {
				request.ResponseChannel <- NewResponse(Val, request.Key, value)
			}
		case request := <-handler.DeleteChannel:
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
