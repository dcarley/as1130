package as1130

import (
	"fmt"
	"image"
	"image/color"
)

var (
	Off = color.Gray{0}   // Min off
	On  = color.Gray{255} // Max on
)

// Framer renders an image into frame data.
type Framer interface {
	OnOffBytes() ([24]byte, error)
}

// rect12x11 is the size of Frame12x11.
func rect12x11() image.Rectangle {
	return image.Rect(0, 0, 12, 11)
}

// Frame12x11 is a frame for a 12x11 matrix with every LED connected.
type Frame12x11 struct {
	*image.Gray
}

// NewFrame12x11 creates a new Frame12x11 of the correct size.
func NewFrame12x11() Frame12x11 {
	return Frame12x11{image.NewGray(rect12x11())}
}

// OnOffBytes renders On/Off LED data. Pixels with colour values greater
// than 0 are considered on.
func (f Frame12x11) OnOffBytes() ([24]byte, error) {
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

// rect24x5 is the size of Frame24x5.
func rect24x5() image.Rectangle {
	return image.Rect(0, 0, 24, 5)
}

// Frame24x5 is a frame for a 24x5 matrix where the last LED in each segment
// is disconnected.
type Frame24x5 struct {
	*image.Gray
}

// NewFrame24x5 creates a new Frame24x5 of the correct size.
func NewFrame24x5() Frame24x5 {
	return Frame24x5{image.NewGray(rect24x5())}
}

// OnOffBytes renders On/Off LED data. Pixels with colour values greater
// than 0 are considered on.
func (f Frame24x5) OnOffBytes() ([24]byte, error) {
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
