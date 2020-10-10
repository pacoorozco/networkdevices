package storer

import (
	"errors"

	"github.com/pacoorozco/networkdevices/internal/models"
)

var (
	ErrDeviceNotFound = errors.New("device not found")
	ErrDeviceAlreadyCreated = errors.New("device already created")
)

// Storer represents the storage where you can keep your devices.
type Storer interface {
	// AddDevice stores the provided Device into the storage if it doesn't exist.
	AddDevice(device models.Device) error
	// GetDevice returns the Device that corresponds with the key or ErrNotFound if it does not exist.
	GetDevice(key string) (models.Device, error)
	// GetAllDevices returns all the Devices existing in the storage.
	GetAllDevices() ([]models.Device, error)
	// SetDevice stores the provided Device into the storage.
	SetDevice(device models.Device) error
	// DeleteDevice removes the stored Device under the provided key.
	DeleteDevice(key string) error
}

