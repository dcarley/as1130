package as1130_test

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	. "github.com/dcarley/as1130"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Frames", func() {
	TestFrameOnOff := func(frameData Framer, expected [24]string) {
		data, err := frameData.OnOffBytes()
		Expect(err).ToNot(HaveOccurred())

		for segmentAddr, _ := range expected {
			By(fmt.Sprintf("frame segment address: 0x%02X", segmentAddr))

			Expect(
				fmt.Sprintf("%08b", data[segmentAddr]),
			).To(
				Equal(expected[segmentAddr]), "segment should match binary representation",
			)
		}
	}

	turnOnEveryOther := func(img draw.Image, odd bool) {
		var mod int
		if odd {
			mod = 1
		}

		var count int
		size := img.Bounds()
		for x := size.Min.X; x < size.Max.X; x++ {
			for y := size.Min.Y; y < size.Max.Y; y++ {
				if count%2 == mod {
					img.Set(x, y, On)
				}
				count++
			}
		}
	}

	Describe("Frame12x11", func() {
		var frame Frame12x11

		BeforeEach(func() {
			frame = NewFrame12x11()
		})

		Describe("OnOffBytes", func() {
			It("should error on incorrect image size", func() {
				frame.Rect.Max.X++

				data, err := frame.OnOffBytes()
				Expect(err).To(MatchError("doesn't match size (0,0)-(13,11): (0,0)-(12,11)"))
				Expect(data).To(Equal([24]byte{}))
			})

			It("should turn all LEDs off by default", func() {
				TestFrameOnOff(frame, [24]string{
					"00000000", "00000000",
					"00000000", "00000000",
					"00000000", "00000000",
					"00000000", "00000000",
					"00000000", "00000000",
					"00000000", "00000000",
					"00000000", "00000000",
					"00000000", "00000000",
					"00000000", "00000000",
					"00000000", "00000000",
					"00000000", "00000000",
					"00000000", "00000000",
				})
			})

			It("should turn all LEDs on", func() {
				draw.Draw(frame, frame.Bounds(), &image.Uniform{On}, image.ZP, draw.Src)

				TestFrameOnOff(frame, [24]string{
					"11111111", "00000111",
					"11111111", "00000111",
					"11111111", "00000111",
					"11111111", "00000111",
					"11111111", "00000111",
					"11111111", "00000111",
					"11111111", "00000111",
					"11111111", "00000111",
					"11111111", "00000111",
					"11111111", "00000111",
					"11111111", "00000111",
					"11111111", "00000111",
				})
			})

			It("should turn on every other LED starting with the first", func() {
				turnOnEveryOther(frame, false)

				TestFrameOnOff(frame, [24]string{
					"01010101", "00000101",
					"10101010", "00000010",
					"01010101", "00000101",
					"10101010", "00000010",
					"01010101", "00000101",
					"10101010", "00000010",
					"01010101", "00000101",
					"10101010", "00000010",
					"01010101", "00000101",
					"10101010", "00000010",
					"01010101", "00000101",
					"10101010", "00000010",
				})
			})

			It("should turn on every other LED starting with the second", func() {
				turnOnEveryOther(frame, true)

				TestFrameOnOff(frame, [24]string{
					"10101010", "00000010",
					"01010101", "00000101",
					"10101010", "00000010",
					"01010101", "00000101",
					"10101010", "00000010",
					"01010101", "00000101",
					"10101010", "00000010",
					"01010101", "00000101",
					"10101010", "00000010",
					"01010101", "00000101",
					"10101010", "00000010",
					"01010101", "00000101",
				})
			})
		})

		Describe("PWMBytes", func() {
			It("should set all LEDs to midpoint", func() {
				draw.Draw(frame, frame.Bounds(), &image.Uniform{color.Gray{128}}, image.ZP, draw.Src)

				data, err := frame.PWMBytes()
				Expect(err).ToNot(HaveOccurred())
				Expect(data).To(Equal([132]byte{
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 128,
				}))
			})
		})
	})

	Describe("FrameLEDs24x5", func() {
		var frame Frame24x5

		BeforeEach(func() {
			frame = NewFrame24x5()
		})

		Describe("OnOffBytes", func() {
			It("should error on incorrect image size", func() {
				frame.Rect.Max.X++
				data, err := frame.OnOffBytes()
				Expect(err).To(MatchError("doesn't match size (0,0)-(25,5): (0,0)-(24,5)"))
				Expect(data).To(Equal([24]byte{}))
			})

			It("should turn all LEDs off by default", func() {
				TestFrameOnOff(frame, [24]string{
					"00000000", "00000000",
					"00000000", "00000000",
					"00000000", "00000000",
					"00000000", "00000000",
					"00000000", "00000000",
					"00000000", "00000000",
					"00000000", "00000000",
					"00000000", "00000000",
					"00000000", "00000000",
					"00000000", "00000000",
					"00000000", "00000000",
					"00000000", "00000000",
				})
			})

			It("should turn all LEDs on", func() {
				draw.Draw(frame, frame.Bounds(), &image.Uniform{On}, image.ZP, draw.Src)

				TestFrameOnOff(frame, [24]string{
					"11111111", "00000011",
					"11111111", "00000011",
					"11111111", "00000011",
					"11111111", "00000011",
					"11111111", "00000011",
					"11111111", "00000011",
					"11111111", "00000011",
					"11111111", "00000011",
					"11111111", "00000011",
					"11111111", "00000011",
					"11111111", "00000011",
					"11111111", "00000011",
				})
			})

			It("should turn on every other LED starting with the first", func() {
				turnOnEveryOther(frame, false)

				TestFrameOnOff(frame, [24]string{
					"01010101", "00000001",
					"01010101", "00000001",
					"01010101", "00000001",
					"01010101", "00000001",
					"01010101", "00000001",
					"01010101", "00000001",
					"01010101", "00000001",
					"01010101", "00000001",
					"01010101", "00000001",
					"01010101", "00000001",
					"01010101", "00000001",
					"01010101", "00000001",
				})
			})

			It("should turn on every other LED starting with the second", func() {
				turnOnEveryOther(frame, true)

				TestFrameOnOff(frame, [24]string{
					"10101010", "00000010",
					"10101010", "00000010",
					"10101010", "00000010",
					"10101010", "00000010",
					"10101010", "00000010",
					"10101010", "00000010",
					"10101010", "00000010",
					"10101010", "00000010",
					"10101010", "00000010",
					"10101010", "00000010",
					"10101010", "00000010",
					"10101010", "00000010",
				})
			})
		})

		Describe("PWMBytes", func() {
			It("should set all LEDs to midpoint ignoring last LED in each segment", func() {
				draw.Draw(frame, frame.Bounds(), &image.Uniform{color.Gray{128}}, image.ZP, draw.Src)

				data, err := frame.PWMBytes()
				Expect(err).ToNot(HaveOccurred())
				Expect(data).To(Equal([132]byte{
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 0,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 0,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 0,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 0,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 0,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 0,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 0,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 0,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 0,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 0,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 0,
					128, 128, 128, 128, 128, 128, 128, 128, 128, 128, 0,
				}))
			})
		})
	})
})
