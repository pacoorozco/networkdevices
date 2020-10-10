package models

import (
	"errors"
)

type Device struct {
	FQDN    string `json:"fqdn"`
	Model   string `json:"model"`
	Version string `json:"version"`
}

// Validate validates the Device data.
func (d *Device) Validate() error {
	if !isFQDN(d.FQDN) {
		return errors.New("attribute FQDN is invalid")
	}
	if !isModel(d.Model) {
		return errors.New("attribute Model is invalid")
	}
	return nil
}

// isFQDN returns true if the provided string is a valid FQDN.
// @see https://en.wikipedia.org/wiki/Fully_qualified_domain_name
func isFQDN(fqdn string) bool {
	// TODO: Validate fqdn is a valid FQDN.
	return len(fqdn) > 0
}

// isModel returns true if the provided string is a valid Model.
func isModel(model string) bool {
	return (model == "ios-xr") || (model == "ios-xe") || (model == "nx-os")
}

func (d *Device) Presenter() Device {
	version := d.Version
	if len(version) == 0 {
		version = "unknown"
	}
	return Device{
		FQDN:    d.FQDN,
		Model:   d.Model,
		Version: version,
	}
}
