package service

import (
	"fmt"
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
	Close()
}

type DataHandler struct {
	store         Datasource
	PutChannel    chan requestWrapper
	GetChannel    chan requestWrapper
	DeleteChannel chan requestWrapper
}

func NewDataHandler(datasource Datasource) DataHandler {
	return DataHandler{
		store:         datasource,
		PutChannel:    make(chan requestWrapper),
		GetChannel:    make(chan requestWrapper),
		DeleteChannel: make(chan requestWrapper),
	}
}

type requestWrapper struct {
	request         Request
	responseChannel chan Response
}

func (handler *DataHandler) CloseStore() {
	handler.store.Close()
}

func (handler *DataHandler) StoreName() string {
	return handler.store.Name()
}

func (handler *DataHandler) ProcessRequest(request Request) Response {
	rw := requestWrapper{request, make(chan Response)}

	switch request.Type {
	case Put:
		handler.PutChannel <- rw
		return <-rw.responseChannel
	case Get:
		handler.GetChannel <- rw
		return <-rw.responseChannel
	case Del:
		handler.DeleteChannel <- rw
		return <-rw.responseChannel
	default:
		return NewResponse(Err, "", "")
	}
}

func (handler *DataHandler) RequestMonitor() {
	for {
		select {
		case req := <-handler.PutChannel:
			err := handler.store.CreateOrUpdate(req.request.Key, req.request.Value)
			if err != nil {
				logg.Error.Printf("For %s - %s", req.request.String(), err)
				req.responseChannel <- NewResponse(Err, "", "")
			} else {
				req.responseChannel <- NewResponse(Ack, req.request.Key, req.request.Value)
			}
		case req := <-handler.GetChannel:
			value, err := handler.store.Read(req.request.Key)
			if err != nil {
				fmt.Println(req)
				logg.Error.Printf("For %s - %s", req.request.String(), err)
				req.responseChannel <- NewResponse(Nil, req.request.Key, req.request.Value)
			} else {
				req.responseChannel <- NewResponse(Val, req.request.Key, value)
			}
		case req := <-handler.DeleteChannel:
			err := handler.store.Delete(req.request.Key)
			if err != nil {
				logg.Error.Printf("For %s - %s", req.request.String(), err)
				req.responseChannel <- NewResponse(Err, "", "")
			} else {
				req.responseChannel <- NewResponse(Ack, req.request.Key, req.request.Value)
			}
		}
	}
}
