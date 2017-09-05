package as1130

import (
	"fmt"
	"time"

	"golang.org/x/exp/io/i2c"
)

const (
	DeviceDefault        = "/dev/i2c-1"
	AddressDefault       = 0x30
	CurrentSourceDefault = 5
)

// Register Selection Address Map (datasheet fig. 31)
const (
	RegisterOnOffFrameFirst    byte = 0x01
	RegisterOnOffFrameLast     byte = 0x24
	RegisterBlinkPWMFrameFirst byte = 0x40
	RegisterBlinkPWMFrameLast  byte = 0x45
	RegisterControl            byte = 0xC0
	RegisterSelect             byte = 0xFD
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

// ClockFreq is a clock frequency value.
type ClockFreq byte

// Clock synchronization frequencies.
const (
	Clock1MHz ClockFreq = iota
	Clock500kHz
	Clock125kHz
	Clock32kHz
)

// LEDs On/Off Frame Register Format (fig. 34)
// and
// LEDs Blink Frame Register Format (fig. 35)
const (
	FrameSegmentFirst byte = 0x00
	FrameSegmentLast  byte = 0x17
)

// LEDs PWM Register Format (fig. 36)
const (
	PWMSegmentFirst byte = 0x18
	PWMSegmentLast  byte = 0x9B
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

// byteToBool converts a 1bit value to a boolean.
func byteToBool(v byte) bool {
	return v == 1
}

// validateClockFreq checks whether a clock out frequency is valid.
func validateClockFreq(freq ClockFreq) bool {
	validFreqs := []ClockFreq{Clock1MHz, Clock500kHz, Clock125kHz, Clock32kHz}
	for _, validFreq := range validFreqs {
		if freq == validFreq {
			return true
		}
	}

	return false
}

// AS1130 is a connected controller.
type AS1130 struct {
	conn            *i2c.Device
	blinkAndPWMSets uint8
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

// Read a register from the AS1130.
func (a *AS1130) Read(register, subregister byte) (byte, error) {
	err := a.conn.WriteReg(RegisterSelect, []byte{register})
	if err != nil {
		return byte(0), err
	}

	buf := make([]byte, 1)
	err = a.conn.ReadReg(subregister, buf)

	return buf[0], err
}

// Init performs the startup sequence with default settings. You still need
// to call Start() when all frames and related settings have been set.
func (a *AS1130) Init(blinkAndPWMSets uint8) error {
	if err := a.SetConfig(Config{BlinkAndPWMSets: blinkAndPWMSets}); err != nil {
		return err
	}
	if err := a.SetCurrentSource(CurrentSourceDefault); err != nil {
		return err
	}
	if err := a.SetDisplayOption(DisplayOption{}); err != nil {
		return err
	}

	return nil
}

// Start takes the device out of shutdown mode, enables the display, and
// starts the state machine.
func (a *AS1130) Start() error {
	return a.SetShutdown(Shutdown{})
}

// Reset performs a soft reset of the device. It shuts down the outputs and
// state machine, then waits for the device to become ready again. Not all
// registers are reset. You will need to call Init() afterwards.
func (a *AS1130) Reset() error {
	reset := Shutdown{
		Initialise: true,
		Shutdown:   true,
	}
	err := a.SetShutdown(reset)
	a.blinkAndPWMSets = 0
	time.Sleep(5 * time.Millisecond)

	return err
}

// MaxFrames returns the total amount of frames that you can use after
// Config.BlinkAndPWMSets has been set.
func (a *AS1130) MaxFrames() (uint8, error) {
	if a.blinkAndPWMSets == 0 {
		return 0, fmt.Errorf("must set Config.BlinkAndPWMSets first")
	}

	const reservedFramesPerSet = 6
	reservedFrames := (a.blinkAndPWMSets - 1) * reservedFramesPerSet

	return RegisterOnOffFrameLast - reservedFrames, nil
}

// Picture Register Format (datasheet fig. 39)
type Picture struct {
	Blink   bool  // All LEDs in blink mode during display picture
	Display bool  // Display picture
	Frame   uint8 // Number of picture frame, 1 if unset
}

// SetPicture sets the picture register.
func (a *AS1130) SetPicture(p Picture) error {
	if p.Frame == 0 {
		p.Frame = 1
	}
	if max, err := a.MaxFrames(); err != nil {
		return err
	} else if p.Frame > max {
		return fmt.Errorf("Frame out of range [1,%d]: %d", max, p.Frame)
	}

	data := boolToByte(p.Blink)<<7 |
		boolToByte(p.Display)<<6 |
		p.Frame - 1

	return a.Write(RegisterControl, ControlPicture, data)
}

// Movie Register Format (datasheet fig. 40)
type Movie struct {
	Blink   bool  // All LEDs in blink mode during play movie
	Display bool  // Display movie
	Frame   uint8 // Number of first frame in movie, 1 if unset
}

// SetMovie sets the movie register.
func (a *AS1130) SetMovie(m Movie) error {
	if m.Frame == 0 {
		m.Frame = 1
	}
	if max, err := a.MaxFrames(); err != nil {
		return err
	} else if m.Frame > max {
		return fmt.Errorf("Frame out of range [1,%d]: %d", max, m.Frame)
	}

	data := boolToByte(m.Blink)<<7 |
		boolToByte(m.Display)<<6 |
		m.Frame - 1

	return a.Write(RegisterControl, ControlMovie, data)
}

// MovieMode Register Format (datasheet fig. 41)
type MovieMode struct {
	Blink   bool  // All LEDs in blink mode during play movie
	EndLast bool  // End movie with last frame instead of first
	Frames  uint8 // Number of frames to play in movie, 1 if unset
}

// SetMovieMode sets the movie mode register.
func (a *AS1130) SetMovieMode(m MovieMode) error {
	if m.Frames == 0 {
		m.Frames = 1
	}
	if max, err := a.MaxFrames(); err != nil {
		return err
	} else if m.Frames > max {
		return fmt.Errorf("Frames out of range [1,%d]: %d", max, m.Frames)
	}

	data := boolToByte(m.Blink)<<7 |
		boolToByte(m.EndLast)<<6 |
		m.Frames - 1

	return a.Write(RegisterControl, ControlMovieMode, data)
}

// FrameTime & Scroll Register Format (datasheet fig. 42)
type FrameTime struct {
	Fade        bool  // Fade at end of frame
	ScrollRight bool  // Scroll right instead of left
	Scroll12x11 bool  // Scroll in 12x11 mode instead of 24x5
	Scrolling   bool  // Scroll digits at play movie
	Delay       uint8 // Delay between frame change in a movie, multiple of 32.5ms
}

// SetFrameTime sets the frame time register.
func (a *AS1130) SetFrameTime(f FrameTime) error {
	if v := f.Delay; v > 15 {
		return fmt.Errorf("Delay out of range [0,15]: %d", v)
	}

	data := boolToByte(f.Fade)<<7 |
		boolToByte(!f.ScrollRight)<<6 |
		boolToByte(!f.Scroll12x11)<<5 |
		boolToByte(f.Scrolling)<<4 |
		f.Delay

	return a.Write(RegisterControl, ControlFrameTime, data)
}

// DisplayOption Register Format (datasheet fig. 43)
type DisplayOption struct {
	Loops          uint8 // Number of loops played in one movie, forever (7) if unset
	BlinkFrequency bool  // Blink every 3s instead of 1.5s
	ScanLimit      uint8 // Number of displayed segments in one frame, all (12) if unset
}

// SetDisplayOption sets the display option.
func (a *AS1130) SetDisplayOption(d DisplayOption) error {
	if d.Loops == 0 {
		d.Loops = 7
	}
	if v := d.Loops; v > 7 {
		return fmt.Errorf("Loops out of range [1,7]: %d", v)
	}
	if d.ScanLimit == 0 {
		d.ScanLimit = 12
	}
	if v := d.ScanLimit; v > 12 {
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
	LowVDDReset        bool  // Reset LowVDD at end of movie or picture
	LowVDDStatus       bool  // Map LowVDD to IRQ pin
	LEDErrorCorrection bool  // Disable open LEDs
	DotCorrection      bool  // Analog current DotCorrection
	CommonAddress      bool  // I2C common address for all AS1130
	BlinkAndPWMSets    uint8 // Number of blink and PWM sets, 1 if unset, each uses 6 On/Off frames
}

// SetConfig sets the config register. The config cannot be changed once you
// have written any frame data, you will need to hard reset the device.
func (a *AS1130) SetConfig(c Config) error {
	if c.BlinkAndPWMSets == 0 {
		c.BlinkAndPWMSets = 1
	}
	if v, max := c.BlinkAndPWMSets, 6; v > 6 {
		return fmt.Errorf("BlinkAndPWMSets out of range [1,%d]: %d", max, v)
	}

	existingConfig, err := a.Read(RegisterControl, ControlConfig)
	if err != nil {
		return err
	}
	existingBlinkAndPWMSets := existingConfig & 0x7
	if existingBlinkAndPWMSets != 0 && existingBlinkAndPWMSets != c.BlinkAndPWMSets {
		return fmt.Errorf("you must hard reset the device to change BlinkAndPWMSets")
	}

	data := boolToByte(c.LowVDDReset)<<7 |
		boolToByte(c.LowVDDStatus)<<6 |
		boolToByte(c.LEDErrorCorrection)<<5 |
		boolToByte(c.DotCorrection)<<4 |
		boolToByte(c.CommonAddress)<<3 |
		c.BlinkAndPWMSets

	err = a.Write(RegisterControl, ControlConfig, data)
	if err == nil {
		a.blinkAndPWMSets = c.BlinkAndPWMSets
	}

	return err
}

// Interrupt Mask Register Format (datasheet fig. 46)
type InterruptMask struct {
	Frame        bool // Trigger IRQ when frame defined by SetInterruptFrame is displayed
	Watchdog     bool // Trigger IRQ when the I²C watchdog triggers
	PowerOrReset bool // Trigger IRQ when power or reset has occurred
	OverTemp     bool // Trigger IRQ when the overtemperature limit is reached
	LowVDD       bool // Trigger IRQ when VDD is too low for used LEDs
	OpenError    bool // Trigger IRQ when an error on the open test occurs
	ShortError   bool // Trigger IRQ when an error on the short test occurs
	MovieFinish  bool // Trigger IRQ when movie is finished
}

// SetInterruptMask sets the interrupt mask register.
func (a *AS1130) SetInterruptMask(i InterruptMask) error {
	data := boolToByte(i.Frame)<<7 |
		boolToByte(i.Watchdog)<<6 |
		boolToByte(i.PowerOrReset)<<5 |
		boolToByte(i.OverTemp)<<4 |
		boolToByte(i.LowVDD)<<3 |
		boolToByte(i.OpenError)<<2 |
		boolToByte(i.ShortError)<<1 |
		boolToByte(i.MovieFinish)

	return a.Write(RegisterControl, ControlInterruptMask, data)
}

// SetInterruptFrame sets the interrupt frame register. This should be used
// in combination with InterruptMask.Frame
func (a *AS1130) SetInterruptFrame(lastFrame uint8) error {
	if max, err := a.MaxFrames(); err != nil {
		return err
	} else if lastFrame < 1 || lastFrame > max {
		return fmt.Errorf("frame out of range [1,%d]: %d", max, lastFrame)
	}

	return a.Write(RegisterControl, ControlInterruptFrame, lastFrame-1)
}

// Shutdown & Open/Short Register Format (datasheet fig. 48)
type Shutdown struct {
	TestAll    bool // LED open/short test is performed on all LED locations
	AutoTest   bool // Automatic LED open/short test is started when picture or movie is displayed
	ManualTest bool // Manual LED open/short test is started after updating shutdown register
	Initialise bool // Initialise control logic (internal state machine is reset again)
	Shutdown   bool // Put device in shutdown mode (outputs are turned off, internal state machine stops)
}

// SetShutdown sets the shutdown register.
func (a *AS1130) SetShutdown(s Shutdown) error {
	data := boolToByte(s.TestAll)<<4 |
		boolToByte(s.AutoTest)<<3 |
		boolToByte(s.ManualTest)<<2 |
		boolToByte(!s.Initialise)<<1 |
		boolToByte(!s.Shutdown)

	return a.Write(RegisterControl, ControlShutdown, data)
}

// CLK Synchronization Register Format (datasheet fig. 50)
type ClockSync struct {
	ClockOutFreq ClockFreq // Adjustable clock out frequency
	SyncOut      bool      // Internal oscillator is used and made available on the SYNC pin
	SyncIn       bool      // Internal oscillator is disabled and input from SYNC pin is used
}

// SetClockSync sets the Clock Synchronization register.
func (a *AS1130) SetClockSync(c ClockSync) error {
	if !validateClockFreq(c.ClockOutFreq) {
		return fmt.Errorf("invalid ClockOutFreq: %d", c.ClockOutFreq)
	}
	if c.SyncOut && c.SyncIn {
		return fmt.Errorf("SyncOut and SyncIn cannot be set at the same time")
	}

	data := byte(c.ClockOutFreq)<<2 |
		boolToByte(c.SyncOut)<<1 |
		boolToByte(c.SyncIn)

	return a.Write(RegisterControl, ControlClockSync, data)
}

// InterruptStatus returns the contents of the interrupt status register.
func (a *AS1130) InterruptStatus() (InterruptMask, error) {
	data, err := a.Read(RegisterControl, ControlInterruptStatus)
	if err != nil {
		return InterruptMask{}, err
	}

	interrupts := InterruptMask{
		Frame:        byteToBool((data & 0x80) >> 7),
		Watchdog:     byteToBool((data & 0x40) >> 6),
		PowerOrReset: byteToBool((data & 0x20) >> 5),
		OverTemp:     byteToBool((data & 0x10) >> 4),
		LowVDD:       byteToBool((data & 0x08) >> 3),
		OpenError:    byteToBool((data & 0x04) >> 2),
		ShortError:   byteToBool((data & 0x02) >> 1),
		MovieFinish:  byteToBool((data & 0x01) >> 0),
	}

	return interrupts, nil
}

// SetFrame sets an On/Off frame.
func (a *AS1130) SetFrame(frame uint8, img Framer) error {
	if max, err := a.MaxFrames(); err != nil {
		return err
	} else if frame < 1 || frame > max {
		return fmt.Errorf("frame out of range [1,%d]: %d", max, frame)
	}

	data, err := img.OnOffBytes()
	if err != nil {
		return err
	}

	pwmSet := img.PWMSet()
	if pwmSet == 0 {
		pwmSet = 1
	}
	if max := a.blinkAndPWMSets; pwmSet > max {
		return fmt.Errorf("PWM set out of range [1,%d]: %d", max, pwmSet)
	}
	data[1] |= (pwmSet - 1) << 5

	registerAddr := RegisterOnOffFrameFirst + (frame - 1)
	for segmentAddr, segmentData := range data {
		if err := a.Write(registerAddr, byte(segmentAddr), segmentData); err != nil {
			return err
		}
	}

	return nil
}

// SetBlinkAndPWMSet sets a blink and PWM set.
func (a *AS1130) SetBlinkAndPWMSet(set uint8, blink Framer, pwm Framer) error {
	if max := a.blinkAndPWMSets; max == 0 {
		return fmt.Errorf("must set Config.BlinkAndPWMSets first")
	} else if set < 1 || set > max {
		return fmt.Errorf("set out of range [1,%d]: %d", max, set)
	}

	blinkData, err := blink.OnOffBytes()
	if err != nil {
		return err
	}
	pwmData, err := pwm.PWMBytes()
	if err != nil {
		return err
	}

	registerAddr := RegisterBlinkPWMFrameFirst + (set - 1)

	var blinkAddr byte
	for i, blinkData := range blinkData {
		blinkAddr = FrameSegmentFirst + byte(i)
		if err := a.Write(registerAddr, blinkAddr, blinkData); err != nil {
			return err
		}
	}

	var pwmAddr byte
	for i, pwmData := range pwmData {
		pwmAddr = PWMSegmentFirst + byte(i)
		if err := a.Write(registerAddr, pwmAddr, pwmData); err != nil {
			return err
		}
	}

	return nil
}
