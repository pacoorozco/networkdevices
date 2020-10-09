package internal

import "github.com/pkg/errors"

var (
	ErrNotFound = errors.New("device not found")
)

// DatabaseAccessLayer represents the storage where you can keep your devices.
type DatabaseAccessLayer interface {
	// Get returns the Device that corresponds with the key or ErrNotFound if it does not exist.
	Get(key string) (Device, error)
	// GetAll returns all the Devices existing in the storage.
	GetAll() (map[int]Device, error)
	// Set stores the provided Device into the storage.
	Set(key string, device Device) error
	// Delete removes the stored Device under the provided key.
	Delete(key string) error
}
