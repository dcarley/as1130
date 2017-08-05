//go:generate go run embedfont/main.go -p as1130 -f "embedfont/CG pixel 4x5.ttf"

package as1130

import (
	"image"
	"image/draw"
	"math"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func loadFont() (font.Face, error) {
	const size = 5

	file, err := fontBytes()
	if err != nil {
		return nil, err
	}
	ttf, err := truetype.Parse(file)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(ttf, &truetype.Options{
		Size: float64(size),
	})

	return face, nil
}

// TextFrames renders some text into a slice of frames that can be scrolled
// across a 24x5 display.
func TextFrames(text string) ([]Frame24x5, error) {
	const spacerFrames = 1

	face, err := loadFont()
	if err != nil {
		return []Frame24x5{}, err
	}

	rect := rect24x5()
	f := &font.Drawer{
		Src:  image.NewUniform(On),
		Face: face,
		Dot:  fixed.P(rect.Max.X*spacerFrames, rect.Max.Y),
	}

	partialFrames := float64(f.MeasureString(text)>>6) / float64(rect.Max.X)
	wholeFrames := int(math.Ceil(partialFrames)) + spacerFrames

	wideFrame := Frame24x5{Gray: image.NewGray(
		image.Rect(0, 0, rect.Max.X*wholeFrames, rect.Max.Y),
	)}

	f.Dst = wideFrame
	f.DrawString(text)

	advance := rect.Max.X
	frames := make([]Frame24x5, wholeFrames)
	for i, _ := range frames {
		frame := NewFrame24x5()
		draw.Draw(frame, frame.Bounds(), wideFrame.SubImage(rect), rect.Min, draw.Src)
		frames[i] = frame

		rect.Min.X += advance
		rect.Max.X += advance
	}

	return frames, nil
}
