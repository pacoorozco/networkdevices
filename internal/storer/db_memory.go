package storer

import (
	"strings"
	"sync"

	"github.com/pacoorozco/networkdevices/internal/models"
)

// InMemoryDAL implements Storer in memory.
type InMemoryDAL struct {
	items map[string]models.Device
	mu    sync.RWMutex
}

// NewInMemory returns a Storer implemented in memory.
func NewInMemory() *InMemoryDAL {
	return &InMemoryDAL{
		items: map[string]models.Device{},
	}
}

// GetDevice returns the Device that corresponds with the key or ErrDeviceNotFound if it does not exist.
func (dal *InMemoryDAL) GetDevice(key string) (models.Device, error) {
	dal.mu.RLock()
	item, found := dal.items[strings.ToLower(key)]
	if !found {
		dal.mu.RUnlock()
		return models.Device{}, ErrDeviceNotFound
	}
	dal.mu.RUnlock()
	return item, nil
}

// GetAllDevices returns all the Devices existing in the storage.
func (dal *InMemoryDAL) GetAllDevices() ([]models.Device, error) {
	dal.mu.RLock()
	defer dal.mu.RUnlock()
	m := make([]models.Device, len(dal.items))
	var i int
	for _, v := range dal.items {
		m = append(m, v)
		i++
	}
	return m, nil
}

// SetDevice stores the provided Device into the storage.
func (dal *InMemoryDAL) SetDevice(item models.Device) error {
	dal.mu.Lock()
	dal.items[strings.ToLower(item.FQDN)] = item
	dal.mu.Unlock()
	return nil
}

// DeleteDevice removes the stored Device under the provided key.
func (dal *InMemoryDAL) DeleteDevice(key string) error {
	dal.mu.Lock()
	delete(dal.items, strings.ToLower(key))
	dal.mu.Unlock()
	return nil
}
