package as1130

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

var (
	Off = color.Gray{0}   // Min off
	On  = color.Gray{255} // Max on
)

// Framer renders an image into frame data.
type Framer interface {
	draw.Image
	OnOffBytes() ([24]byte, error)
	PWMBytes() ([132]byte, error)
	SetPWMSet(uint8)
	PWMSet() uint8
}

// rect12x11 is the size of Frame12x11.
func rect12x11() image.Rectangle {
	return image.Rect(0, 0, 12, 11)
}

// Frame12x11 is a frame for a 12x11 matrix with every LED connected.
type Frame12x11 struct {
	*image.Gray
	pwmSet uint8
}

// NewFrame12x11 creates a new Frame12x11 of the correct size.
func NewFrame12x11() *Frame12x11 {
	return &Frame12x11{
		Gray: image.NewGray(rect12x11()),
	}
}

// OnOffBytes renders On/Off LED data. Pixels with colour values greater
// than 0 are considered on or blinking.
func (f *Frame12x11) OnOffBytes() ([24]byte, error) {
	data := [24]byte{}
	if actual, expected := f.Bounds(), rect12x11(); actual != expected {
		return data, fmt.Errorf("doesn't match size %s: %s", actual, expected)
	}

	size := f.Bounds()
	for x := size.Min.X; x < size.Max.X; x++ {
		var shift uint8
		segmentAddr := x * 2

		for y := size.Min.Y; y < size.Max.Y; y++ {
			if shift == 8 {
				shift = 0
				segmentAddr++
			}

			if f.GrayAt(x, y).Y > 0 {
				data[segmentAddr] |= 1 << shift
			}
			shift++
		}
	}

	return data, nil
}

// PWMBytes renders PWM (brightness) LED data. Each pixel has 255 steps.
func (f *Frame12x11) PWMBytes() ([132]byte, error) {
	data := [132]byte{}
	if f.Bounds() != rect12x11() {
		return data, fmt.Errorf("XXX")
	}

	var dataIndex uint8
	size := f.Bounds()
	for x := size.Min.X; x < size.Max.X; x++ {
		for y := size.Min.Y; y < size.Max.Y; y++ {
			data[dataIndex] = f.GrayAt(x, y).Y
			dataIndex++
		}
	}

	return data, nil
}

// PWMSet returns the PWM set associated with the frame.
func (f *Frame12x11) PWMSet() uint8 {
	return f.pwmSet
}

// SetPWMSet sets the PWM set to associate with the frame.
func (f *Frame12x11) SetPWMSet(set uint8) {
	f.pwmSet = set
}

// rect24x5 is the size of Frame24x5.
func rect24x5() image.Rectangle {
	return image.Rect(0, 0, 24, 5)
}

// Frame24x5 is a frame for a 24x5 matrix where the last LED in each segment
// is disconnected.
type Frame24x5 struct {
	*image.Gray
	pwmSet uint8
}

// NewFrame24x5 creates a new Frame24x5 of the correct size.
func NewFrame24x5() *Frame24x5 {
	return &Frame24x5{
		Gray: image.NewGray(rect24x5()),
	}
}

// OnOffBytes renders On/Off LED data. Pixels with colour values greater
// than 0 are considered on or blinking.
func (f *Frame24x5) OnOffBytes() ([24]byte, error) {
	data := [24]byte{}
	if actual, expected := f.Bounds(), rect24x5(); actual != expected {
		return data, fmt.Errorf("doesn't match size %s: %s", actual, expected)
	}

	size := f.Bounds()
	for x := size.Min.X; x < size.Max.X; x += 2 {
		var shift uint8
		segmentAddr := x

		for xOffset := 0; xOffset < 2; xOffset++ {
			for y := size.Min.Y; y < size.Max.Y; y++ {
				if shift == 8 {
					shift = 0
					segmentAddr++
				}

				if f.GrayAt(x+xOffset, y).Y > 0 {
					data[segmentAddr] |= 1 << shift
				}
				shift++
			}
		}
	}

	return data, nil
}

// PWMBytes renders PWM (brightness) LED data. Each pixel has 255 steps.
func (f *Frame24x5) PWMBytes() ([132]byte, error) {
	data := [132]byte{}
	if f.Bounds() != rect24x5() {
		return data, fmt.Errorf("XXX")
	}

	var dataIndex uint8
	size := f.Bounds()
	for x := size.Min.X; x < size.Max.X; x++ {
		for y := size.Min.Y; y < size.Max.Y; y++ {
			data[dataIndex] = f.GrayAt(x, y).Y
			dataIndex++
		}

		if x%2 == 1 {
			dataIndex++
		}
	}

	return data, nil
}

// PWMSet returns the PWM set associated with the frame.
func (f *Frame24x5) PWMSet() uint8 {
	return f.pwmSet
}

// SetPWMSet sets the PWM set to associate with the frame.
func (f *Frame24x5) SetPWMSet(set uint8) {
	f.pwmSet = set
}
