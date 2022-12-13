package service

import (
	"fmt"
	"os"
	"testing"
)

var (
	ks            KeyStore
	putRequest    Request
	getRequest    Request
	deleteRequest Request
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	ks = NewKeyStore()
	putRequest = NewRequest("", "testKeyPut", "testValuePut")
	getRequest = NewRequest("", "testKeyGet", "testValueGet")
	deleteRequest = NewRequest("", "testKeyDel", "testValueDel")
}

func shutdown() {
	ks.keyStore = nil
	close(ks.putChannel)
	close(ks.getChannel)
	close(ks.deleteChannel)
	close(ks.ResponseChannel)
}

func TestKeyStore_CreateOrUpdate(t *testing.T) {
	err := ks.CreateOrUpdate(putRequest)
	if err != nil {
		t.Error("Could not insert value into key store")
	}
	value, ok := ks.keyStore[putRequest.Key]
	if !ok || value != putRequest.Value {
		t.Error("Value was not stored")
	}
}

func TestKeyStore_Read(t *testing.T) {
	ks.keyStore[getRequest.Key] = getRequest.Value
	value, err := ks.Read(getRequest)
	if err != nil {
		t.Error("Expected key, but key not found")
	} else if value != getRequest.Value {
		t.Error(fmt.Sprintf("Expected %s but got %s", getRequest.Value, value))
	}
}

func TestKeyStore_Delete(t *testing.T) {
	ks.keyStore[deleteRequest.Key] = deleteRequest.Value
	err := ks.Delete(deleteRequest)
	if err != nil {
		t.Error("Key deletion failed")
	}
	_, ok := ks.keyStore[deleteRequest.Key]
	if ok {
		t.Error("Key remains after deletion")
	}
}

func TestPutWhenStoreClosed(t *testing.T) {
	ks.keyStore = nil
	err := ks.CreateOrUpdate(putRequest)
	if err == nil {
		t.Error("Expected store closed error")
	}
}

func TestGetWhenStoreClosed(t *testing.T) {
	ks.keyStore = nil
	_, err := ks.Read(getRequest)
	if err == nil {
		t.Error("Expected store closed error")
	}
}

func TestDeleteWhenStoreClosed(t *testing.T) {
	ks.keyStore = nil
	err := ks.Delete(deleteRequest)
	if err == nil {
		t.Error("Expected store closed error")
	}
}

func TestKeyStoreIsClosed(t *testing.T) {
	ks.keyStore = nil
	closed := ks.isClosed()
	if !closed {
		t.Error("Expected true error got ", false)
	}
}
