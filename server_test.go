package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func newServer() (*KVServer, *httptest.ResponseRecorder) {
	return &KVServer{kvStore: NewKeyValueStore()}, httptest.NewRecorder()
}

func TestListKeysEndpoint(t *testing.T) {
	t.Run("it returns an empty list when the store is empty", func(t *testing.T) {
		server, w := newServer()
		server.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))

		if w.Result().StatusCode != 200 {
			t.Errorf("Expected 200, got %d", w.Result().StatusCode)
		}

		if w.Body.String() != "[]\n" {
			t.Errorf("Expected an emtpy array, but got '%v'", w.Body.String())
		}
	})

	t.Run("it returns a list of keys", func(t *testing.T) {
		server, w := newServer()

		server.kvStore.Put("some-key", []byte("some-value"))
		server.kvStore.Put("some-other-key", []byte("some-value"))

		server.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
		if w.Result().StatusCode != 200 {
			t.Errorf("Expected 200, got %d", w.Result().StatusCode)
			return
		}

		decoder := json.NewDecoder(w.Body)
		var keys []string
		err := decoder.Decode(&keys)
		if err != nil {
			t.Errorf("Failed to read response json: %v", err)
		}

		if len(keys) != 2 {
			t.Errorf("Expected 2 keys, got %d", len(keys))
		}
	})
}

func TestGetKeyEndpoint(t *testing.T) {
	t.Run("it returns a 404 when the key does not exist", func(t *testing.T) {
		server, w := newServer()
		server.ServeHTTP(w, httptest.NewRequest("GET", "/key", nil))

		if w.Result().StatusCode != 404 {
			t.Errorf("Expected 404, got %d", w.Result().StatusCode)
		}
	})

	t.Run("it returns the key data if the key exists", func(t *testing.T) {
		server, w := newServer()

		server.kvStore.Put("some-key", []byte("some-value"))
		server.kvStore.Put("some-other-key", []byte("some-value"))

		server.ServeHTTP(w, httptest.NewRequest("GET", "/some-key", nil))
		if w.Result().StatusCode != 200 {
			t.Errorf("Expected 200, got %d", w.Result().StatusCode)
			return
		}

		if w.Body.String() != "some-value" {
			t.Errorf("Expected 'some-value', got '%s'", w.Body.String())
		}
	})
}

func TestDeleteKeyEndpoint(t *testing.T) {
	t.Run("it returns a 404 when the key does not exist", func(t *testing.T) {
		server, w := newServer()
		server.ServeHTTP(w, httptest.NewRequest("DELETE", "/key", nil))

		if w.Result().StatusCode != 404 {
			t.Errorf("Expected 404, got %d", w.Result().StatusCode)
		}
	})

	t.Run("it removes the key", func(t *testing.T) {
		server, w := newServer()

		server.kvStore.Put("some-key", []byte("some-value"))
		server.kvStore.Put("some-other-key", []byte("some-value"))

		server.ServeHTTP(w, httptest.NewRequest("DELETE", "/some-key", nil))
		if w.Result().StatusCode != 200 {
			t.Errorf("Expected 200, got %d", w.Result().StatusCode)
			return
		}

		keys, _ := server.kvStore.GetKeys()

		if len(keys) != 1 {
			t.Errorf("Expected store to have 1 key, but got %d", len(keys))
		}
	})
}

func TestPutKeyEndpoint(t *testing.T) {
	t.Run("it creates a key if it did not exist", func(t *testing.T) {
		server, w := newServer()
		body := bytes.NewBufferString("some-value")
		server.ServeHTTP(w, httptest.NewRequest("PUT", "/key", body))

		if w.Result().StatusCode != 200 {
			t.Errorf("Expected 200, got %d", w.Result().StatusCode)
			return
		}

		keys, _ := server.kvStore.GetKeys()
		if len(keys) != 1 {
			t.Errorf("Expected 1 key, got %d", len(keys))
			return
		}

		value, _ := server.kvStore.Get("key")
		if string(value) != "some-value" {
			t.Errorf("Expected some-value got '%s'", string(value))
		}
	})
	t.Run("it updates a key if it does exist", func(t *testing.T) {
		server, w := newServer()

		server.kvStore.Put("key", []byte("some-value"))

		body := bytes.NewBufferString("some-other-value")
		server.ServeHTTP(w, httptest.NewRequest("PUT", "/key", body))

		if w.Result().StatusCode != 200 {
			t.Errorf("Expected 200, got %d", w.Result().StatusCode)
			return
		}

		value, _ := server.kvStore.Get("key")
		if string(value) != "some-other-value" {
			t.Errorf("Expected some-other-value got '%s'", string(value))
		}
	})

}
