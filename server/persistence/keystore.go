package persistence

import (
	"errors"
	"tcpstore/service"
)

var (
	ErrKvStoreDoesNotExist = errors.New("key value store has not been initialized")
	ErrKeyNotFound         = errors.New("key not found")
)

type KeyStore struct {
	name     string
	keyStore map[service.Key]service.Value
}

func NewKeyStore() KeyStore {
	return KeyStore{
		name:     "In-Memory Map",
		keyStore: map[service.Key]service.Value{},
	}
}

func (ks *KeyStore) Name() string {
	return ks.name
}

func (ks *KeyStore) CreateOrUpdate(key service.Key, value service.Value) error {
	if ks.IsClosed() {
		return ErrKvStoreDoesNotExist
	}
	ks.keyStore[key] = value
	return nil
}

func (ks *KeyStore) Read(key service.Key) (service.Value, error) {
	if ks.IsClosed() {
		return "", ErrKvStoreDoesNotExist
	}
	value, ok := ks.keyStore[key]
	if !ok {
		return "", ErrKeyNotFound
	}
	return value, nil
}

func (ks *KeyStore) Delete(key service.Key) error {
	if ks.IsClosed() {
		return ErrKvStoreDoesNotExist
	}
	delete(ks.keyStore, key)
	return nil
}

func (ks *KeyStore) IsClosed() bool {
	return ks.keyStore == nil
}
