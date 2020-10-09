package internal

import (
	"sync"
)

// InMemoryDAL implements DataAccessLayer interface in memory.
type InMemoryDAL struct {
	items map[string]Device
	mu    sync.RWMutex
}

// NewInMemory returns a DataAccessLayer implemented in memory.
func NewInMemory() *InMemoryDAL {
	return &InMemoryDAL{
		items: map[string]Device{},
	}
}

// Get returns the Device that corresponds with the key or ErrNotFound if it does not exist.
func (dal *InMemoryDAL) Get(key string) (Device, error) {
	dal.mu.RLock()
	item, found := dal.items[key]
	if !found {
		dal.mu.RUnlock()
		return Device{}, ErrNotFound
	}
	dal.mu.RUnlock()
	return item, nil
}

// GetAll returns all the Devices existing in the storage.
func (dal *InMemoryDAL) GetAll() (map[int]Device, error) {
	dal.mu.RLock()
	defer dal.mu.RUnlock()
	m := make(map[int]Device, len(dal.items))
	var i int
	for _, v := range dal.items {
		m[i] = v
		i++
	}
	return m, nil
}

// Set stores the provided Device into the storage.
func (dal *InMemoryDAL) Set(key string, item Device) error {
	dal.mu.Lock()
	dal.items[key] = item
	dal.mu.Unlock()
	return nil
}

// Delete removes the stored Device under the provided key.
func (dal *InMemoryDAL) Delete(key string) error {
	dal.mu.Lock()
	delete(dal.items, key)
	dal.mu.Unlock()
	return nil
}
