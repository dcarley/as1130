package as1130

import (
	"golang.org/x/exp/io/i2c"
)

const (
	DeviceDefault  = "/dev/i2c-1"
	AddressDefault = 0x30
)

// Register Selection Address Map (datasheet fig. 31)
const (
	RegisterControl byte = 0xC0
	RegisterSelect  byte = 0xFD
)

// Control Register Address Map (datasheet fig. 38)
const (
	ControlPicture byte = iota
	ControlMovie
	ControlMovieMode
	ControlFrameTime
	ControlDisplayOption
	ControlCurrentSource
	ControlConfig
	ControlInterruptMask
	ControlInterruptFrame
	ControlShutdown
	ControlI2CMonitoring
	ControlClockSync
	ControlInterruptStatus
	ControlStatus
)

// AS1130 is a connected controller.
type AS1130 struct {
	Conn *i2c.Device
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

	return &AS1130{Conn: conn}, nil
}

// Close closes the connection to an AS1130 controller.
func (a *AS1130) Close() error {
	return a.Conn.Close()
}

// Write sends a command to the AS1130.
func (a *AS1130) Write(register, subregister, data byte) error {
	err := a.Conn.WriteReg(RegisterSelect, []byte{register})
	if err != nil {
		return err
	}

	return a.Conn.WriteReg(subregister, []byte{data})
}
