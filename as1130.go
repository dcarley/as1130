package as1130

import (
	"fmt"

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

// boolToByte converts a boolean to a 1bit value. It allows us to provide
// human friendly structs and not have to do runtime validation on 1bit
// binary options.
func boolToByte(v bool) byte {
	if v {
		return 1
	}

	return 0
}

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

// Write sends a command to the AS1130.
func (a *AS1130) Write(register, subregister, data byte) error {
	err := a.conn.WriteReg(RegisterSelect, []byte{register})
	if err != nil {
		return err
	}

	return a.conn.WriteReg(subregister, []byte{data})
}

// DisplayOption Register Format (datasheet fig. 43)
type DisplayOption struct {
	Loops          uint8 // Number of loops played in one movie
	BlinkFrequency bool  // Blink every 3s instead of 1.5s
	ScanLimit      uint8 // Number of displayed segments in one frame
}

// SetDisplayOption sets the display option.
func (a *AS1130) SetDisplayOption(d DisplayOption) error {
	if v := d.Loops; v > 7 {
		return fmt.Errorf("Loops out of range [0,7]: %d", v)
	}
	if v := d.ScanLimit; v < 1 || v > 12 {
		return fmt.Errorf("ScanLimit out of range [1,12]: %d", v)
	}

	data := d.Loops<<5 |
		boolToByte(d.BlinkFrequency)<<4 |
		d.ScanLimit - 1

	return a.Write(RegisterControl, ControlDisplayOption, data)
}

// SetCurrentSource sets the current (mA) for all LEDs.
func (a *AS1130) SetCurrentSource(milliAmps byte) error {
	if milliAmps > 30 {
		return fmt.Errorf("current out of range [0,30]: %d", milliAmps)
	}

	return a.Write(RegisterControl, ControlCurrentSource, byte(int(milliAmps)*255/30.0))
}

// Config Register Format (datasheet fig. 45)
type Config struct {
	LowVDDReset         bool  // Reset LowVDD at end of movie or picture
	LowVDDStatus        bool  // Map LowVDD to IRQ pin
	LEDErrorCorrection  bool  // Disable open LEDs
	DotCorrection       bool  // Analog current DotCorrection
	CommonAddress       bool  // I2C common address for all AS1130
	MemoryConfiguration uint8 // RAM Configuration, 1 if unset
}

// SetConfig sets the config register.
func (a *AS1130) SetConfig(c Config) error {
	if c.MemoryConfiguration == 0 {
		c.MemoryConfiguration = 1
	}
	if v := c.MemoryConfiguration; v > 6 {
		return fmt.Errorf("MemoryConfiguration out of range [1,6]: %d", v)
	}

	data := boolToByte(c.LowVDDReset)<<7 |
		boolToByte(c.LowVDDStatus)<<6 |
		boolToByte(c.LEDErrorCorrection)<<5 |
		boolToByte(c.DotCorrection)<<4 |
		boolToByte(c.CommonAddress)<<3 |
		c.MemoryConfiguration

	return a.Write(RegisterControl, ControlConfig, data)
}
