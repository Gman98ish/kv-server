package main

import (
	"testing"
)

func TestCanGetAndStoreData(t *testing.T) {
	store := NewKeyValueStore()

	err := store.Put("my-key", []byte("Some value"))
	if err != nil {
		t.Errorf("Failed to put key: %v", err)
		return
	}

	data, err := store.Get("my-key")
	if err != nil {
		t.Errorf("Failed to get key: %v", err)
		return
	}

	if string(data) != "Some value" {
		t.Errorf("Expected 'Some value' got '%s'", string(data))
	}
}

func TestCanRemoveKeys(t *testing.T) {
	store := NewKeyValueStore()

	err := store.Put("my-key", []byte("Some value"))
	if err != nil {
		t.Errorf("Failed to put key: %v", err)
		return
	}

	exists, err := store.Delete("my-key")
	if err != nil {
		t.Errorf("Failed to delete key: %v", err)
		return
	}

	if !exists {
		t.Errorf("Expected key to exist")
		return
	}

	data, err := store.Get("my-key")
	if err != nil {
		t.Errorf("Failed to get key: %v", err)
		return
	}

	if data != nil {
		t.Errorf("Expected 'my-key' to not have any data, but got %v", data)
	}
}

func TestCanListKeys(t *testing.T) {
	store := NewKeyValueStore()

	store.Put("my-key1", []byte("Some value"))
	store.Put("my-key2", []byte("Some value"))
	store.Put("my-key3", []byte("Some value"))

	keys, err := store.GetKeys()
	if err != nil {
		t.Errorf("Failed to get keys, %v", err)
		return
	}

	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

}

func parallelReadTest(kv *KeyValueStore) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		_, err := kv.Get("my-cool-key")
		if err != nil {
			t.Errorf("Failed to read data: %v", err)
		}

	}
}

func parallelWriteTest(kv *KeyValueStore) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		err := kv.Put("my-cool-key", []byte("Some cool data"))
		if err != nil {
			t.Errorf("Failed to write data: %v", err)
		}

	}
}
func TestCanSafelyReadAndWriteConcurrently(t *testing.T) {
	kv := NewKeyValueStore()
	t.Run("read1", parallelReadTest(kv))
	t.Run("write1", parallelWriteTest(kv))
	t.Run("write2", parallelWriteTest(kv))
	t.Run("write3", parallelWriteTest(kv))
	t.Run("read2", parallelReadTest(kv))

}
