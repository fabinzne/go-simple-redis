package storage

import "sync"

type DataStore struct {
	items map[string]any
	mutex sync.RWMutex
}

func NewDataStore() *DataStore {
	return &DataStore{
		items: make(map[string]interface{}),
	}
}

func (ds *DataStore) Set(key string, value any) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	ds.items[key] = value
}

func (ds *DataStore) Get(key string) (any, bool) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()
	val, ok := ds.items[key]
	return val, ok
}

func (ds *DataStore) Delete(key string) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	delete(ds.items, key)
}
