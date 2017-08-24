package as1130_test

import (
	. "github.com/dcarley/as1130"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Text", func() {
	Describe("TextFrames", func() {
		const text = "hello world as1130 test"
		var frames []*Frame24x5

		BeforeEach(func() {
			var err error
			frames, err = TextFrames(text)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should have 1 frame per 4 characters on average", func() {
			const (
				framesPerChar = 4
				spacerFrames  = 1
			)

			Expect(frames).To(HaveLen(spacerFrames + len(text)/framesPerChar))
		})

		It("should start with an empty frame", func() {
			frame := frames[0]
			size := frame.Bounds()

			var nonEmptyPixels int
			for x := size.Min.X; x < size.Max.X; x++ {
				for y := size.Min.Y; y < size.Max.Y; y++ {
					if frame.GrayAt(x, y) != Off {
						nonEmptyPixels++
					}

				}
			}

			Expect(nonEmptyPixels).To(Equal(0), "frame 1 to have no pixels set")
		})

		It("should have have text in all other frames", func() {
			for i, frame := range frames {
				if i == 0 {
					continue
				}

				size := frame.Bounds()
				var nonEmptyPixels int
				for x := size.Min.X; x < size.Max.X; x++ {
					for y := size.Min.Y; y < size.Max.Y; y++ {
						if frame.GrayAt(x, y) == On {
							nonEmptyPixels++
						}
					}
				}

				Expect(nonEmptyPixels).To(BeNumerically(">", 10),
					"frame %d to contain enough non-empty pixels", i+1,
				)
			}
		})
	})
})
