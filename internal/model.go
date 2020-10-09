package internal

type Device struct {
	FQDN    string `json:"fqdn"`
	Model   string `json:"model"`
	Version string `json:"version"`
}

func (d *Device) getDevice(dal DatabaseAccessLayer) error {
	data, err := dal.Get(d.FQDN)
	if err != nil {
		return err
	}
	d.Model = data.Model
	d.Version = data.Version
	return nil
}

func (d *Device) updateDevice(dal DatabaseAccessLayer) error {
	return dal.Set(d.FQDN, *d)

}

func (d *Device) deleteDevice(dal DatabaseAccessLayer) error {
	return dal.Delete(d.FQDN)
}

func (d *Device) createDevice(dal DatabaseAccessLayer) error {
	return dal.Set(d.FQDN, *d)
}

func getDevices(dal DatabaseAccessLayer) ([]Device, error) {
	data, err := dal.GetAll()
	if err != nil {
		return make([]Device, 0), err
	}
	res := make([]Device, len(data))
	for _, v := range data {
		res = append(res, v)
	}
	return res, nil
}
