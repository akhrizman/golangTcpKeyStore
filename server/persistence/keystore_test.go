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

func TestKeyStore_Name(t *testing.T) {
	if ks.Name() != "In-Memory Map" {
		t.Error("Keystore should be named 'In-Memory Map'")
	}
}

func shutdown() {
	ks.keyStore = nil
}

func TestKeyStore_CreateOrUpdate(t *testing.T) {
	err := ks.CreateOrUpdate(putRequest.Key, putRequest.Value)
	if err != nil {
		t.Error("Could not insert value into key store")
	}
}

func TestKeyStore_Read(t *testing.T) {
	ks.keyStore[getRequest.Key] = getRequest.Value
	t.Run("Key exists", func(t *testing.T) {
		value, _ := ks.Read(getRequest.Key)
		if value != getRequest.Value {
			t.Errorf("Expected %s but got %s", getRequest.Value, value)
		}
	})
	t.Run("Key missing", func(t *testing.T) {
		value, err := ks.Read("ThisKeyDoesn'tExist")
		if err == nil || value != "" {
			t.Error("Expected key not found error")
		}
	})

}

func TestKeyStore_Delete(t *testing.T) {
	ks.keyStore[deleteRequest.Key] = deleteRequest.Value
	err := ks.Delete(deleteRequest.Key)
	if err != nil {
		t.Error("key deletion failed")
	}
}

func TestKeyStore_CrudMethodsWhenStoreClosed(t *testing.T) {
	ks.keyStore = nil
	t.Run("CreateOrUpdate", func(t *testing.T) {
		err := ks.CreateOrUpdate(putRequest.Key, putRequest.Value)
		if err == nil {
			t.Error("Expected store closed error")
		}
	})
	t.Run("Read", func(t *testing.T) {
		_, err := ks.Read(getRequest.Key)
		if err == nil {
			t.Error("Expected store closed error")
		}
	})
	t.Run("Delete", func(t *testing.T) {
		err := ks.Delete(deleteRequest.Key)
		if err == nil {
			t.Error("Expected store closed error")
		}
	})
}

func TestKeyStoreIsClosed(t *testing.T) {
	ks.keyStore = nil
	closed := ks.IsClosed()
	if !closed {
		t.Error("Expected true error got ", false)
	}
}
