package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

// KVServer is a HTTP handler for a key value store
type KVServer struct {
	kvStore KeyValueStorer
}

func (s *KVServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	_, key, _ := strings.Cut(req.URL.Path, "/")

	if key == "" && req.Method == "GET" {
		s.HandleList(rw)
		return
	}

	if req.Method == http.MethodGet && key != "" {
		s.HandleGet(key, rw)
		return
	}

	if req.Method == http.MethodPut && key != "" {
		s.HandlePut(key, rw, req)
		return
	}

	if req.Method == http.MethodDelete && key != "" {
		s.HandleDelete(key, rw)
		return
	}

	encode(map[string]string{
		"message": "not found",
	}, rw, 404)
}

// HandleList gets all keys and returns them in a JSON array
func (s *KVServer) HandleList(rw http.ResponseWriter) {
	keys, err := s.kvStore.GetKeys()
	if err != nil {
		ServerError(rw, fmt.Errorf("error fetching keys: %w", err))
		return
	}

	encode(keys, rw, 200)
}

// HandleGet gets a specific key and returns its content directly
func (s *KVServer) HandleGet(key string, rw http.ResponseWriter) {
	data, err := s.kvStore.Get(key)
	if err != nil {
		ServerError(rw, fmt.Errorf("error getting key '%s': %w", key, err))
		return
	}

	if data == nil {
		encode(map[string]string{
			"message": fmt.Sprintf("No such key %s", key),
		}, rw, 404)
	}

	rw.Write(data)
}

// HandleDelete removes a specific key from the store. If the key doesn't
// exist then a not found response will be returned
func (s *KVServer) HandleDelete(key string, rw http.ResponseWriter) {
	exists, err := s.kvStore.Delete(key)
	if !exists {
		encode(map[string]string{
			"message": fmt.Sprintf("No such key %s", key),
		}, rw, 404)
		return
	}

	if err != nil {
		ServerError(rw, fmt.Errorf("error deleting key '%s': %w", key, err))
	}
}

// HandlePut will update a key in place if the key exists. If not the key will
// be created
func (s *KVServer) HandlePut(key string, rw http.ResponseWriter, req *http.Request) {
	data, err := io.ReadAll(req.Body)
	if err != nil {
		ServerError(rw, fmt.Errorf("failed to ready request body: %w", err))
		return
	}

	err = s.kvStore.Put(key, data)
	if err != nil {
		ServerError(rw, fmt.Errorf("failed to store key '%s': %w", key, err))
	}
}

// ServerError writes a 500 response with a json body
// indicating there has been a server error. it also logs the error
// using slog
func ServerError(rw http.ResponseWriter, err error) {
	slog.Error("server error: %w", err)
	encode(map[string]string{
		"message": "server error",
	}, rw, 500)
}

// encode is a generic function for writing JSON responses
func encode[T any](data T, rw http.ResponseWriter, status int) error {
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(status)

	encoder := json.NewEncoder(rw)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to write response body: %w", err)
	}

	return nil
}
