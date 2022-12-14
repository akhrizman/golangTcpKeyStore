package service

import (
	"errors"
	"fmt"
)

var (
	Store                  = NewKeyStore()
	ErrKvStoreDoesNotExist = errors.New("key value store has not been initialized")
	ErrKeyNotFound         = errors.New("key not found")
)

type KeyStore struct {
	keyStore      map[Key]Value
	putChannel    chan Request
	getChannel    chan Request
	deleteChannel chan Request
}

func NewKeyStore() KeyStore {
	return KeyStore{
		keyStore:      map[Key]Value{},
		putChannel:    make(chan Request),
		getChannel:    make(chan Request),
		deleteChannel: make(chan Request),
	}
}

//func (ks *KeyStore) Put(request Request) {
//	fmt.Println("Sending Request to put channel: ", request)
//	ks.putChannel <- request
//}
//
//func (ks *KeyStore) Get(request Request) {
//	fmt.Println("Sending Request to get channel: ", request)
//	ks.getChannel <- request
//}
//
//func (ks *KeyStore) Del(request Request) {
//	fmt.Println("Sending Request to delete channel: ", request)
//	ks.deleteChannel <- request
//}

func (ks *KeyStore) CreateOrUpdate(request Request) error {
	if ks.isClosed() {
		return ErrKvStoreDoesNotExist
	}
	ks.keyStore[request.Key] = request.Value
	return nil
}

func (ks *KeyStore) Read(request Request) (Value, error) {
	if ks.isClosed() {
		return "", ErrKvStoreDoesNotExist
	}
	value, ok := ks.keyStore[request.Key]
	if !ok {
		return "", ErrKeyNotFound
	}
	return value, nil
}

func (ks *KeyStore) Delete(request Request) {
	delete(ks.keyStore, request.Key)
}

func (ks *KeyStore) isClosed() bool {
	return ks.keyStore == nil
}

func (ks *KeyStore) QueueRequest(request Request) {
	switch request.Task {
	case "put":
		fmt.Println("Sending Request to put channel: ", request)
		ks.putChannel <- request
	case "get":
		fmt.Println("Sending Request to get channel: ", request)
		ks.getChannel <- request
	case "del":
		fmt.Println("Sending Request to delete channel: ", request)
		ks.deleteChannel <- request
	default:
		request.ResponseChannel <- NewResponse("err", "", "")
	}
}

func (ks *KeyStore) RequestMonitor() {
	for {
		select {
		case request := <-ks.putChannel:
			fmt.Println("Received request from put channel: ", request)
			err := ks.CreateOrUpdate(request)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Put Response generated and passed to response channel: ", request)
				request.ResponseChannel <- NewResponse("ack", request.Key, request.Value)
			}
		case request := <-ks.getChannel:
			fmt.Println("Received request from get channel: ", request)
			value, err := ks.Read(request)
			if err != nil {
				fmt.Println(err)
				request.ResponseChannel <- NewResponse("nil", request.Key, request.Value)
			} else {
				fmt.Println("Get Response generated and passed to response channel: ", request)
				request.ResponseChannel <- NewResponse("val", request.Key, value)
			}
		case request := <-ks.deleteChannel:
			fmt.Println("Received request from delete channel: ", request)
			ks.Delete(request)
			fmt.Println("Delete Response generated and passed to response channel: ", request)
			request.ResponseChannel <- NewResponse("ack", request.Key, request.Value)
		}
	}
}
