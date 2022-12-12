package keystore

import (
	"errors"
	"fmt"
)

var (
	Store                  = NewKeyStore()
	ErrKvStoreDoesNotExist = errors.New("key value store has not been initialized")
	ErrKeyNotFound         = errors.New("key not found")
)

type Key string
type Value string

type Request struct {
	Task  string
	Key   Key
	Value Value
}

type Response struct {
	acknowledgement string // val, del, ack, nil
	key             string
	value           string
}

func (response *Response) ClientString() string {
	switch response.acknowledgement {
	case "val":
		return fmt.Sprint(response.acknowledgement, 123, response.key, response.value)
	case "del":
		return fmt.Sprint(response.acknowledgement, 456, response.key)
	default:
		return response.acknowledgement
	}
}

func NewResponse(acknowledgement string, key Key, value Value) Response {
	return Response{
		acknowledgement: acknowledgement,
		key:             string(key),
		value:           string(value),
	}
}

type KeyStore struct {
	keyStore        map[Key]Value
	putChannel      chan Request
	getChannel      chan Request
	deleteChannel   chan Request
	ResponseChannel chan Response
}

func NewKeyStore() KeyStore {
	return KeyStore{
		keyStore:        map[Key]Value{},
		putChannel:      make(chan Request),
		getChannel:      make(chan Request),
		deleteChannel:   make(chan Request),
		ResponseChannel: make(chan Response),
	}
}

func (ks *KeyStore) QueueRequest(request Request) {
	switch request.Task {
	case "put":
		ks.putChannel <- request
	case "get":
		ks.getChannel <- request
	case "del:":
		ks.deleteChannel <- request
	default:
		fmt.Println("Task invalid")
	}
}

func (ks *KeyStore) CommandMonitor() {
	for {
		select {
		case request := <-ks.putChannel:
			err := ks.Put(request)
			if err != nil {
				fmt.Println(err)
			} else {
				ks.ResponseChannel <- NewResponse("ack", request.Key, request.Value)
			}
		case request := <-ks.getChannel:
			value, err := ks.Get(request)
			if err != nil {
				fmt.Println(err)
				ks.ResponseChannel <- NewResponse("nil", request.Key, request.Value)
			} else {
				ks.ResponseChannel <- NewResponse("val", request.Key, value)
			}
		case request := <-ks.deleteChannel:
			err := ks.Delete(request)
			if err != nil {
				fmt.Println(err)
			} else {
				ks.ResponseChannel <- NewResponse("del", request.Key, request.Value)
			}
		}
	}
}

func (ks *KeyStore) Put(request Request) error {
	if ks.isClosed() {
		return ErrKvStoreDoesNotExist
	}
	ks.keyStore[Key(request.Key)] = Value(request.Value)
	return nil
}

func (ks *KeyStore) Get(request Request) (Value, error) {
	if ks.isClosed() {
		return "", ErrKvStoreDoesNotExist
	}
	value, ok := ks.keyStore[request.Key]
	if !ok {
		return "", ErrKeyNotFound
	}
	return value, nil
}

func (ks *KeyStore) Delete(request Request) error {
	if ks.isClosed() {
		return ErrKvStoreDoesNotExist
	}
	_, ok := ks.keyStore[request.Key]
	if !ok {
		return ErrKeyNotFound
	}
	delete(ks.keyStore, request.Key)
	return nil
}

func (ks *KeyStore) isClosed() bool {
	return ks.keyStore == nil
}
