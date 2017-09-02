// Package as1130 is a library for controlling the AS1130 LED driver, as
// used by The Matrix from Boldport.
//
// Configuration
//
// The device is configured using registers. Each register is represented as
// a struct with descriptive field names, comments, and default values. The
// struct can be passed to the corresponding Set method.
//
// You may also wish to refer to the datasheet:
// http://ams.com/eng/content/download/185846/834724/105034
//
// Connecting
//
// To connect using the default device path and address:
//
//		as, err := as1130.NewAS1130("", 0)
//
// If you need to use another path or address:
//
//		as, err := as1130.NewAS1130("/dev/i2c-2", 0x31)
//
// Initialisation
//
// It's advisable to reset the device before using it:
//
//		err = as.Reset()
//
// After which you must choose how many blink and PWM sets you want to use
// and perform the startup sequence. There's a helper method to make this
// easier:
//
//		err = as.Init(1)
//
// Or you can call SetConfig(), SetCurrentSource() and SetDisplayOption()
// yourself if you need more control.
//
// Frames
//
// Frames are defined as grayscale images. They can be manipulated using
// standard Go image libraries:
//
//		ledsOn := as1130.NewFrame24x5()
//		draw.Draw(ledsOn, ledsOn.Bounds(), &image.Uniform{as1130.On}, image.ZP, draw.Src)
//		err = as.SetFrame(1, ledsOn)
//
// Each frame is associated with blink and PWM set, which defines flashing
// and intensity for individual LEDs. You should set the first (default)
// blink and PWM set at a minimum:
//
//		noBlink , pwmFull := as1130.NewFrame24x5(), as1130.NewFrame24x5()
//		draw.Draw(pwm, pwm.Bounds(), &image.Uniform{as1130.On}, image.ZP, draw.Src)
//		err = as.SetBlinkAndPWMSet(1, noBlink, pwmFull)
//
// Pictures
//
// An individual frame can be displayed in picture mode:
//
//		err = as.SetPicture(as1130.Picture{Display: true})
//
// Movies
//
// Multiple frames can be displayed in movie mode:
//
//		err = as.SetMovie(as1130.Movie{Display: true})
//		err = as.SetMovieMode(as1130.MovieMode{Frames: 2})
//
// Start
//
// When you have finished setting all of the frames you can turn on the
// display:
//
//		err = as.Start()
//
// Examples
//
// Have a look at the as1130ctl sub-package for more complete examples.
//
package as1130
