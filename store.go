package main

import "sync"

// KeyValueStorer represents a generic key value store
type KeyValueStorer interface {
	// Put should set the value for the provided key.
	// If the key does not exist then it should be created.
	// If the key already exists, then it should be updated
	Put(key string, value []byte) error

	// Get should retrieve the value for the provided key
	// if the key is not present in the database then it should return nil for the data
	Get(key string) ([]byte, error)

	// Delete should remove the key.
	// If the key does not exist it should return false
	// If the key does exist it should return true
	Delete(key string) (bool, error)

	// GetKeys should return a list of all keys in the store
	GetKeys() ([]string, error)
}

// KeyValueStore is a simple in memory store. It handles concurrent
// reads and writes safely. All data is copied on read and write to
// avoid references outside the store updating values unexpectedly
type KeyValueStore struct {
	keyValues map[string][]byte
	lock      *sync.RWMutex
}

func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{
		keyValues: make(map[string][]byte),
		lock:      &sync.RWMutex{},
	}
}

// Put copies the provided value to the key
func (s *KeyValueStore) Put(key string, value []byte) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.keyValues[key] = make([]byte, len(value))
	copy(s.keyValues[key], value)

	return nil
}

// Get returns a copy of the value held at key. If the key does not
// exist, then a nil slice is returned
func (s *KeyValueStore) Get(key string) ([]byte, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	val, exists := s.keyValues[key]
	if !exists {
		return nil, nil
	}

	dst := make([]byte, len(val))
	copy(dst, val)

	return dst, nil
}

// Delete removes a key from the store
func (s *KeyValueStore) Delete(key string) (bool, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, exists := s.keyValues[key]; !exists {
		return false, nil
	}

	delete(s.keyValues, key)

	return true, nil
}

// GetKeys returns a slice of available keys in the store
func (s *KeyValueStore) GetKeys() ([]string, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var keys []string = make([]string, 0)
	for key := range s.keyValues {
		keys = append(keys, key)
	}

	return keys, nil
}
