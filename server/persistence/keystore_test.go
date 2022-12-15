package persistence

import (
	"os"
	"tcpstore/service"
	"testing"
)

var (
	ks            KeyStore
	putRequest    service.Request
	getRequest    service.Request
	deleteRequest service.Request
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	ks = NewKeyStore()
	putRequest = service.NewRequest("", "testKeyPut", "testValuePut")
	getRequest = service.NewRequest("", "testKeyGet", "testValueGet")
	deleteRequest = service.NewRequest("", "testKeyDel", "testValueDel")
}

func shutdown() {
	ks.keyStore = nil
}

func TestKeyStore_CreateOrUpdate(t *testing.T) {
	err := ks.CreateOrUpdate(putRequest.Key, putRequest.Value)
	if err != nil {
		t.Error("Could not insert value into key store")
	}
	value, ok := ks.keyStore[putRequest.Key]
	if !ok || value != putRequest.Value {
		t.Error("value was not stored")
	}
}

func TestKeyStore_Read(t *testing.T) {
	ks.keyStore[getRequest.Key] = getRequest.Value
	value, err := ks.Read(getRequest.Key)
	if err != nil {
		t.Error("Expected key, but key not found")
	} else if value != getRequest.Value {
		t.Errorf("Expected %s but got %s", getRequest.Value, value)
	}
}

func TestKeyStore_Delete(t *testing.T) {
	ks.keyStore[deleteRequest.Key] = deleteRequest.Value
	err := ks.Delete(deleteRequest.Key)
	if err != nil {
		t.Error("key deletion failed")
	}
	_, ok := ks.keyStore[deleteRequest.Key]
	if ok {
		t.Error("key remains after deletion")
	}
}

func TestPutWhenStoreClosed(t *testing.T) {
	ks.keyStore = nil
	err := ks.CreateOrUpdate(putRequest.Key, putRequest.Value)
	if err == nil {
		t.Error("Expected store closed error")
	}
}

func TestGetWhenStoreClosed(t *testing.T) {
	ks.keyStore = nil
	_, err := ks.Read(getRequest.Key)
	if err == nil {
		t.Error("Expected store closed error")
	}
}

func TestDeleteWhenStoreClosed(t *testing.T) {
	ks.keyStore = nil
	err := ks.Delete(deleteRequest.Key)
	if err == nil {
		t.Error("Expected store closed error")
	}
}

func TestKeyStoreIsClosed(t *testing.T) {
	ks.keyStore = nil
	closed := ks.IsClosed()
	if !closed {
		t.Error("Expected true error got ", false)
	}
}
