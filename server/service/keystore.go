package service

import (
	"errors"
	"fmt"
)

var (
	Store                  = NewKeyStore()
	ErrKvStoreDoesNotExist = errors.New("key value store has not been initialized")
	ErrKeyNotFound         = errors.New("key not found")
	ResponseStreams        map[int]chan Response
)

type KeyStore struct {
	keyStore      map[Key]Value
	putChannel    chan Request
	getChannel    chan Request
	deleteChannel chan Request
	//ResponseChannel chan Response
}

func NewKeyStore() KeyStore {
	return KeyStore{
		keyStore:      map[Key]Value{},
		putChannel:    make(chan Request),
		getChannel:    make(chan Request),
		deleteChannel: make(chan Request),
		//ResponseChannel: make(chan Response),
	}
}

func (ks *KeyStore) Put(request Request) {
	fmt.Println("Sending Request to put channel: ", request)
	ks.putChannel <- request
}

func (ks *KeyStore) Get(request Request) {
	fmt.Println("Sending Request to get channel: ", request)
	ks.getChannel <- request
}

func (ks *KeyStore) Del(request Request) {
	fmt.Println("Sending Request to delete channel: ", request)
	ks.deleteChannel <- request
}

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
	//if ks.isClosed() {
	//	return ErrKvStoreDoesNotExist
	//}
	delete(ks.keyStore, request.Key)
}

func (ks *KeyStore) isClosed() bool {
	return ks.keyStore == nil
}

func (ks *KeyStore) QueueRequest(request Request) {
	switch request.Task {
	case "put":
		ks.Put(request)
	case "get":
		ks.Get(request)
	case "del":
		ks.Del(request)
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
				ResponseStreams[request.ConnectionId] <- NewResponse("ack", request.Key, request.Value, request.ConnectionId)
			}
		case request := <-ks.getChannel:
			fmt.Println("Received request from get channel: ", request)
			value, err := ks.Read(request)
			if err != nil {
				fmt.Println(err)
				ResponseStreams[request.ConnectionId] <- NewResponse("nil", request.Key, request.Value, request.ConnectionId)
			} else {
				fmt.Println("Get Response generated and passed to response channel: ", request)
				ResponseStreams[request.ConnectionId] <- NewResponse("val", request.Key, value, request.ConnectionId)
			}
		case request := <-ks.deleteChannel:
			fmt.Println("Received request from delete channel: ", request)
			ks.Delete(request)
			fmt.Println("Delete Response generated and passed to response channel: ", request)
			ResponseStreams[request.ConnectionId] <- NewResponse("ack", request.Key, request.Value, request.ConnectionId)
		}
	}
}
