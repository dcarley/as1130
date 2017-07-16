package as1130

import (
	"golang.org/x/exp/io/i2c"
)

const (
	DeviceDefault  = "/dev/i2c-1"
	AddressDefault = 0x30
)

// AS1130 is a connected controller.
type AS1130 struct {
	conn *i2c.Device
}

// NewAS1130 opens a connection to an AS1130 controller. If `device` or
// `address` are zero values for their type then the defaults will be used.
// You must call Close() when you're done.
func NewAS1130(device string, address int) (*AS1130, error) {
	if device == "" {
		device = DeviceDefault
	}
	if address == 0 {
		address = AddressDefault
	}

	conn, err := i2c.Open(&i2c.Devfs{Dev: device}, address)
	if err != nil {
		return nil, err
	}

	return &AS1130{conn: conn}, nil
}

// Close closes the connection to an AS1130 controller.
func (a *AS1130) Close() error {
	return a.conn.Close()
}
