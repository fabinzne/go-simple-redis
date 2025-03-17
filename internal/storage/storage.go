package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type DataStore struct {
	items       map[string]any
	expirations map[string]time.Time
	timers      map[string]*time.Timer
	mutex       sync.RWMutex
}

type StorageEntry struct {
	Value      interface{} `json:"value"`
	Expiration int64       `json:"expiration,omitempty"` // Timestamp UnixNano
}

func NewDataStore() *DataStore {
	return &DataStore{
		items:       make(map[string]interface{}),
		expirations: make(map[string]time.Time),
		timers:      make(map[string]*time.Timer),
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
	fmt.Println(val, ok)
	return val, ok
}

func (ds *DataStore) Delete(key string) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	delete(ds.items, key)
}

func (ds *DataStore) SetWithTTL(key string, value any, ttl time.Duration) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	ds.items[key] = value
	if ttl > 0 {
		if timer, ok := ds.timers[key]; ok {
			timer.Stop()
		}
		ds.timers[key] = time.AfterFunc(ttl*time.Second, func() {
			ds.Delete(key)
		})
	}
}

func (ds *DataStore) Expire(key string, ttl time.Duration) bool {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	if _, ok := ds.items[key]; !ok {
		return false
	}
	if timer, exists := ds.timers[key]; exists {
		timer.Stop()
		delete(ds.timers, key)
	}

	expiration := time.Now().Add(ttl * time.Second)
	ds.expirations[key] = expiration

	ds.timers[key] = time.AfterFunc(ttl*time.Second, func() {
		ds.Delete(key)
	})

	return true
}

func (ds *DataStore) isExpired(key string) bool {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	expiration, exists := ds.expirations[key]
	return exists && time.Now().After(expiration)
}

func (ds *DataStore) Flush() {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	ds.items = make(map[string]interface{})
	ds.expirations = make(map[string]time.Time)
	for key, timer := range ds.timers {
		timer.Stop()
		delete(ds.timers, key)
	}
}

func (ds *DataStore) SafeGetKeyInfo(key string) (interface{}, bool, time.Time) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	val, ok := ds.items[key]
	exp, expOk := ds.expirations[key]
	if !expOk {
		return val, ok, time.Time{}
	}
	return val, ok, exp
}

func (ds *DataStore) SaveToFile(filename string) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	data := make(map[string]StorageEntry)

	for key, value := range ds.items {
		entry := StorageEntry{
			Value: value,
		}
		if exp, exists := ds.expirations[key]; exists {
			entry.Expiration = exp.UnixNano()
		}
		data[key] = entry
	}

	file, err := os.CreateTemp(filepath.Dir(filename), "temp-")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return os.Rename(file.Name(), filename)
}

func (ds *DataStore) LoadFromFile(filename string) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var data map[string]StorageEntry
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return err
	}

	now := time.Now()
	for key, entry := range data {
		ds.items[key] = entry.Value
		if entry.Expiration > 0 {
			expTime := time.Unix(0, entry.Expiration)
			if expTime.Before(now) {
				continue
			}

			ds.expirations[key] = expTime
			ttl := expTime.Sub(now)

			ds.timers[key] = time.AfterFunc(ttl, func() {
				ds.Delete(key)
			})
		}
	}

	return nil
}
