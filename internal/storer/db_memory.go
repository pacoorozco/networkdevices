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
		items: make(map[string]models.Device, 0),
	}
}

// GetDevice returns the Device that corresponds with the key or ErrDeviceNotFound if it does not exist.
func (dal *InMemoryDAL) GetDevice(key string) (models.Device, error) {
	dal.mu.RLock()
	defer dal.mu.RUnlock()
	item, found := dal.items[strings.ToLower(key)]
	if !found {
		return models.Device{}, ErrDeviceNotFound
	}
	return item, nil
}

// GetAllDevices returns all the Devices existing in the storage.
func (dal *InMemoryDAL) GetAllDevices() ([]models.Device, error) {
	dal.mu.RLock()
	defer dal.mu.RUnlock()
	m := make([]models.Device, 0)
	for _, v := range dal.items {
		m = append(m, v.Presenter())
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

// AddDevice stores the provided Device into the storage if it doesn't exist.
// Returns ErrDeviceAlreadyCreated if device exists before creation.
func (dal *InMemoryDAL) AddDevice(item models.Device) error {
	dal.mu.Lock()
	defer dal.mu.Unlock()
	if _, found := dal.items[strings.ToLower(item.FQDN)]; found {
		return ErrDeviceAlreadyCreated
	}
	dal.items[strings.ToLower(item.FQDN)] = item
	return nil
}

// DeleteDevice removes the stored Device under the provided key.
func (dal *InMemoryDAL) DeleteDevice(key string) error {
	dal.mu.Lock()
	delete(dal.items, strings.ToLower(key))
	dal.mu.Unlock()
	return nil
}
