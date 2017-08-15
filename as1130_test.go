package as1130

import (
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"

	"golang.org/x/exp/io/i2c"

	"github.com/dcarley/as1130/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("as1130", func() {
	TestReadCommand := func(buf *gbytes.Buffer, register, subregister byte) {
		const size = 3
		command := make([]byte, size)
		read, err := buf.Read(command)
		Expect(err).ToNot(HaveOccurred())
		Expect(read).To(Equal(size), "shoud have three bytes of command")

		Expect(command[0]).To(Equal(RegisterSelect), "should call RegisterSelect")
		Expect(command[1]).To(Equal(register), "should select the register")
		Expect(command[2]).To(Equal(subregister), "should select the subregister")
	}

	TestRemainingCommands := func(buf *gbytes.Buffer, commands int) {
		const bytesPerCommand = 4
		allCommands, err := ioutil.ReadAll(buf)
		Expect(err).ToNot(HaveOccurred())
		Expect(allCommands).To(HaveLen(commands * bytesPerCommand))
	}

	TestCommand := func(buf *gbytes.Buffer, register, subregister byte, binaryData string) {
		const bytesPerCommand = 4
		command := make([]byte, bytesPerCommand)
		read, err := buf.Read(command)
		Expect(err).ToNot(HaveOccurred())
		Expect(read).To(Equal(bytesPerCommand), "shoud have four bytes of command")

		Expect(command[0]).To(Equal(RegisterSelect), "should call RegisterSelect")
		Expect(command[1]).To(Equal(register), "should select the register")
		Expect(command[2]).To(Equal(subregister), "should select the subregister")
		Expect(
			fmt.Sprintf("%08b", command[3]),
		).To(
			Equal(binaryData), "should write data to the register",
		)
	}

	Describe("AS1130", func() {
		var (
			as                *AS1130
			writeBuf, readBuf *gbytes.Buffer
		)

		BeforeEach(func() {
			writeBuf, readBuf = gbytes.NewBuffer(), gbytes.NewBuffer()
			opener := &fakes.FakeOpener{
				W: writeBuf,
				R: readBuf,
			}
			conn, err := i2c.Open(opener, 0x30)
			Expect(err).ToNot(HaveOccurred())

			as = &AS1130{conn: conn}
		})

		AfterEach(func() {
			TestRemainingCommands(writeBuf, 0)

			Expect(as.Close()).To(Succeed())
			Expect(writeBuf.Closed()).To(BeTrue())
			Expect(readBuf.Closed()).To(BeTrue())
		})

		Describe("Write", func() {
			It("should write a command to the device", func() {
				Expect(as.Write(RegisterControl, ControlConfig, 1)).To(Succeed())
				TestCommand(writeBuf, RegisterControl, ControlConfig, "00000001")
			})
		})

		Describe("MaxFrames", func() {
			BeforeEach(func() {
				_, err := readBuf.Write([]byte{0x0})
				Expect(err).ToNot(HaveOccurred())
			})

			Context("Config.BlinkAndPWMSets is unset", func() {
				It("MaxFrames returns an error", func() {
					max, err := as.MaxFrames()
					Expect(err).To(MatchError("must set Config.BlinkAndPWMSets first"))
					Expect(max).To(Equal(uint8(0)))
				})
			})

			DescribeTable("Config.BlinkAndPWMSets is set",
				func(blinkAndPWMSets, maxFrames int) {
					config := Config{BlinkAndPWMSets: uint8(blinkAndPWMSets)}
					Expect(as.SetConfig(config)).To(Succeed())

					max, err := as.MaxFrames()
					Expect(err).ToNot(HaveOccurred())
					Expect(max).To(Equal(uint8(maxFrames)))

					TestReadCommand(writeBuf, RegisterControl, ControlConfig)
					TestRemainingCommands(writeBuf, 1)
				},
				Entry("to 1", 1, 36),
				Entry("to 2", 2, 30),
				Entry("to 3", 3, 24),
				Entry("to 4", 4, 18),
				Entry("to 5", 5, 12),
				Entry("to 6", 6, 6),
			)
		})

		Describe("SetPicture", func() {
			const (
				register    = RegisterControl
				subregister = ControlPicture
			)

			Context("Config.BlinkAndPWMSets has not been set", func() {
				It("should error", func() {
					Expect(as.SetPicture(Picture{})).To(MatchError("must set Config.BlinkAndPWMSets first"))
					Expect(writeBuf.Contents()).To(BeEmpty())
				})
			})

			Context("Config.BlinkAndPWMSets has been set", func() {
				BeforeEach(func() {
					as.blinkAndPWMSets = 1
				})

				It("should write defaults", func() {
					picture := Picture{}
					Expect(as.SetPicture(picture)).To(Succeed())
					TestCommand(writeBuf, register, subregister, "00000000")
				})

				It("should write non-defaults", func() {
					picture := Picture{
						Blink:   true,
						Display: true,
						Frame:   36,
					}
					Expect(as.SetPicture(picture)).To(Succeed())
					TestCommand(writeBuf, register, subregister, "11100011")
				})

				It("should error on too high frame", func() {
					picture := Picture{
						Frame: 37,
					}
					Expect(as.SetPicture(picture)).To(MatchError("Frame out of range [1,36]: 37"))
					Expect(writeBuf.Contents()).To(BeEmpty())
				})
			})
		})

		Describe("SetMovie", func() {
			const (
				register    = RegisterControl
				subregister = ControlMovie
			)

			Context("Config.BlinkAndPWMSets has not been set", func() {
				It("should error", func() {
					Expect(as.SetMovie(Movie{})).To(MatchError("must set Config.BlinkAndPWMSets first"))
					Expect(writeBuf.Contents()).To(BeEmpty())
				})
			})

			Context("Config.BlinkAndPWMSets has been set", func() {
				BeforeEach(func() {
					as.blinkAndPWMSets = 1
				})

				It("should write defaults", func() {
					movie := Movie{}
					Expect(as.SetMovie(movie)).To(Succeed())
					TestCommand(writeBuf, register, subregister, "00000000")
				})

				It("should write non-defaults", func() {
					movie := Movie{
						Blink:   true,
						Display: true,
						Frame:   36,
					}
					Expect(as.SetMovie(movie)).To(Succeed())
					TestCommand(writeBuf, register, subregister, "11100011")
				})

				It("should error on too high frame", func() {
					picture := Movie{
						Frame: 37,
					}
					Expect(as.SetMovie(picture)).To(MatchError("Frame out of range [1,36]: 37"))
					Expect(writeBuf.Contents()).To(BeEmpty())
				})
			})
		})

		Describe("SetMovieMode", func() {
			const (
				register    = RegisterControl
				subregister = ControlMovieMode
			)

			Context("Config.BlinkAndPWMSets has not been set", func() {
				It("should error", func() {
					Expect(as.SetMovieMode(MovieMode{})).To(MatchError("must set Config.BlinkAndPWMSets first"))
					Expect(writeBuf.Contents()).To(BeEmpty())
				})
			})

			Context("Config.BlinkAndPWMSets has been set", func() {
				BeforeEach(func() {
					as.blinkAndPWMSets = 1
				})

				It("should write defaults", func() {
					movieMode := MovieMode{}
					Expect(as.SetMovieMode(movieMode)).To(Succeed())
					TestCommand(writeBuf, register, subregister, "00000000")
				})

				It("should write non-defaults", func() {
					movieMode := MovieMode{
						Blink:   true,
						EndLast: true,
						Frames:  36,
					}
					Expect(as.SetMovieMode(movieMode)).To(Succeed())
					TestCommand(writeBuf, register, subregister, "11100011")
				})

				It("should error on too high frame", func() {
					picture := MovieMode{
						Frames: 37,
					}
					Expect(as.SetMovieMode(picture)).To(MatchError("Frames out of range [1,36]: 37"))
					Expect(writeBuf.Contents()).To(BeEmpty())
				})
			})
		})

		Describe("SetFrameTime", func() {
			const (
				register    = RegisterControl
				subregister = ControlFrameTime
			)

			It("should write defaults", func() {
				frameTime := FrameTime{}
				Expect(as.SetFrameTime(frameTime)).To(Succeed())
				TestCommand(writeBuf, register, subregister, "01100000")
			})

			It("should write non-defaults", func() {
				frameTime := FrameTime{
					Fade:        true,
					ScrollRight: true,
					Scroll12x11: true,
					Scrolling:   true,
					Delay:       15,
				}
				Expect(as.SetFrameTime(frameTime)).To(Succeed())
				TestCommand(writeBuf, register, subregister, "10011111")
			})

			It("should error on out of range Delay", func() {
				frameTime := FrameTime{
					Delay: 16,
				}
				Expect(as.SetFrameTime(frameTime)).To(MatchError("Delay out of range [0,15]: 16"))
				Expect(writeBuf.Contents()).To(BeEmpty())
			})
		})

		Describe("SetDisplayOption", func() {
			const (
				register    = RegisterControl
				subregister = ControlDisplayOption
			)

			It("should write defaults", func() {
				option := DisplayOption{}
				Expect(as.SetDisplayOption(option)).To(Succeed())
				TestCommand(writeBuf, register, subregister, "11101011")
			})

			It("should write non-defaults", func() {
				option := DisplayOption{
					Loops:          1,
					BlinkFrequency: true,
					ScanLimit:      12,
				}
				Expect(as.SetDisplayOption(option)).To(Succeed())
				TestCommand(writeBuf, register, subregister, "00111011")
			})

			It("should error on out of range Loops", func() {
				option := DisplayOption{
					Loops: 8,
				}
				Expect(as.SetDisplayOption(option)).To(MatchError("Loops out of range [1,7]: 8"))
				Expect(writeBuf.Contents()).To(BeEmpty())
			})

			It("should error on out of range ScanLimit", func() {
				option := DisplayOption{
					ScanLimit: 13,
				}
				Expect(as.SetDisplayOption(option)).To(MatchError("ScanLimit out of range [1,12]: 13"))
				Expect(writeBuf.Contents()).To(BeEmpty())
			})
		})

		Describe("SetCurrentSource", func() {
			const (
				register    = RegisterControl
				subregister = ControlCurrentSource
			)

			It("should set to min", func() {
				Expect(as.SetCurrentSource(0)).To(Succeed())
				TestCommand(writeBuf, register, subregister, fmt.Sprintf("%08b", byte(0)))
			})

			It("should set to half way", func() {
				Expect(as.SetCurrentSource(15)).To(Succeed())
				TestCommand(writeBuf, register, subregister, fmt.Sprintf("%08b", byte(127)))
			})

			It("should set to max", func() {
				Expect(as.SetCurrentSource(30)).To(Succeed())
				TestCommand(writeBuf, register, subregister, fmt.Sprintf("%08b", byte(255)))
			})

			It("should error on invalid current", func() {
				Expect(as.SetCurrentSource(31)).To(MatchError("current out of range [0,30]: 31"))
			})
		})

		Describe("SetConfig", func() {
			const (
				register    = RegisterControl
				subregister = ControlConfig
			)

			Context("BlinkAndPWMSets has not been previously set", func() {
				BeforeEach(func() {
					var err error
					_, err = readBuf.Write([]byte{0x0})
					Expect(err).ToNot(HaveOccurred())
				})

				It("should write defaults", func() {
					config := Config{}
					Expect(as.SetConfig(config)).To(Succeed())

					TestReadCommand(writeBuf, register, subregister)
					TestCommand(writeBuf, register, subregister, "00000001")
				})

				It("should write non-defaults", func() {
					config := Config{
						LowVDDReset:        true,
						LowVDDStatus:       true,
						LEDErrorCorrection: true,
						DotCorrection:      true,
						CommonAddress:      true,
						BlinkAndPWMSets:    6,
					}
					Expect(as.SetConfig(config)).To(Succeed())

					TestReadCommand(writeBuf, register, subregister)
					TestCommand(writeBuf, register, subregister, "11111110")
				})

				It("should error on out of range BlinkAndPWMSets", func() {
					config := Config{
						BlinkAndPWMSets: 7,
					}
					Expect(as.SetConfig(config)).To(MatchError("BlinkAndPWMSets out of range [1,6]: 7"))
					Expect(writeBuf.Contents()).To(BeEmpty())
				})
			})

			Context("BlinkAndPWMSets has been previously set to 2", func() {
				BeforeEach(func() {
					var err error
					_, err = readBuf.Write([]byte{0x2})
					Expect(err).ToNot(HaveOccurred())
				})

				It("should write the same BlinkAndPWMSets", func() {
					config := Config{
						BlinkAndPWMSets: 2,
					}
					Expect(as.SetConfig(config)).To(Succeed())

					TestReadCommand(writeBuf, register, subregister)
					TestCommand(writeBuf, register, subregister, "00000010")
				})

				It("should error on different BlinkAndPWMSets", func() {
					config := Config{
						BlinkAndPWMSets: 1,
					}
					Expect(as.SetConfig(config)).To(MatchError(
						"you must hard reset the device to change BlinkAndPWMSets",
					))

					TestReadCommand(writeBuf, register, subregister)
					TestRemainingCommands(writeBuf, 0)
				})
			})
		})

		Describe("SetInterruptMask", func() {
			const (
				register    = RegisterControl
				subregister = ControlInterruptMask
			)

			It("should write defaults", func() {
				mask := InterruptMask{}
				Expect(as.SetInterruptMask(mask)).To(Succeed())
				TestCommand(writeBuf, register, subregister, "00000000")
			})

			It("should write non-defaults", func() {
				mask := InterruptMask{
					Frame:        true,
					Watchdog:     true,
					PowerOrReset: true,
					OverTemp:     true,
					LowVDD:       true,
					OpenError:    true,
					ShortError:   true,
					MovieFinish:  true,
				}
				Expect(as.SetInterruptMask(mask)).To(Succeed())
				TestCommand(writeBuf, register, subregister, "11111111")
			})
		})

		Describe("SetInterruptFrame", func() {
			const (
				register    = RegisterControl
				subregister = ControlInterruptFrame
			)

			Context("Config.BlinkAndPWMSets has not been set", func() {
				It("should error", func() {
					Expect(as.SetInterruptFrame(1)).To(MatchError("must set Config.BlinkAndPWMSets first"))
					Expect(writeBuf.Contents()).To(BeEmpty())
				})
			})

			Context("Config.BlinkAndPWMSets set to 1", func() {
				BeforeEach(func() {
					as.blinkAndPWMSets = 1
				})

				It("should write first frame", func() {
					Expect(as.SetInterruptFrame(1)).To(Succeed())
					TestCommand(writeBuf, register, subregister, "00000000")
				})

				It("should write last frame", func() {
					Expect(as.SetInterruptFrame(36)).To(Succeed())
					TestCommand(writeBuf, register, subregister, "00100011")
				})

				It("should error on zero indexed frame", func() {
					Expect(as.SetInterruptFrame(0)).To(MatchError("frame out of range [1,36]: 0"))
					Expect(writeBuf.Contents()).To(BeEmpty())
				})

				It("should error on too high frame", func() {
					Expect(as.SetInterruptFrame(37)).To(MatchError("frame out of range [1,36]: 37"))
					Expect(writeBuf.Contents()).To(BeEmpty())
				})
			})
		})

		Describe("SetShutdown", func() {
			const (
				register    = RegisterControl
				subregister = ControlShutdown
			)

			It("should write defaults", func() {
				shutdown := Shutdown{}
				Expect(as.SetShutdown(shutdown)).To(Succeed())
				TestCommand(writeBuf, register, subregister, "00000011")
			})

			It("should write non-defaults", func() {
				shutdown := Shutdown{
					TestAll:    true,
					AutoTest:   true,
					ManualTest: true,
					Initialise: true,
					Shutdown:   true,
				}
				Expect(as.SetShutdown(shutdown)).To(Succeed())
				TestCommand(writeBuf, register, subregister, "00011100")
			})
		})

		Describe("SetClockSync", func() {
			const (
				register    = RegisterControl
				subregister = ControlClockSync
			)

			It("should write defaults", func() {
				clockSync := ClockSync{}
				Expect(as.SetClockSync(clockSync)).To(Succeed())
				TestCommand(writeBuf, register, subregister, "00000000")
			})

			It("should write ClockOutFreq as 32kHz", func() {
				clockSync := ClockSync{
					ClockOutFreq: Clock32kHz,
				}
				Expect(as.SetClockSync(clockSync)).To(Succeed())
				TestCommand(writeBuf, register, subregister, "00001100")
			})

			It("should error on invalid ClockOutFreq", func() {
				clockSync := ClockSync{
					ClockOutFreq: Clock32kHz + 1,
				}
				Expect(as.SetClockSync(clockSync)).To(MatchError("invalid ClockOutFreq: 4"))
				Expect(writeBuf.Contents()).To(BeEmpty())
			})

			It("should write SyncOut", func() {
				clockSync := ClockSync{
					SyncOut: true,
				}
				Expect(as.SetClockSync(clockSync)).To(Succeed())
				TestCommand(writeBuf, register, subregister, "00000010")
			})

			It("should write SyncIn", func() {
				clockSync := ClockSync{
					SyncIn: true,
				}
				Expect(as.SetClockSync(clockSync)).To(Succeed())
				TestCommand(writeBuf, register, subregister, "00000001")
			})

			It("should error on conflicting SyncOut and SyncIn", func() {
				clockSync := ClockSync{
					SyncOut: true,
					SyncIn:  true,
				}
				Expect(as.SetClockSync(clockSync)).To(MatchError("SyncOut and SyncIn cannot be set at the same time"))
				Expect(writeBuf.Contents()).To(BeEmpty())
			})
		})

		Describe("SetFrame", func() {
			const commandsPerFrame = int(FrameSegmentLast - FrameSegmentFirst + 1)
			var frame *Frame12x11

			BeforeEach(func() {
				frame = NewFrame12x11()
				draw.Draw(frame, frame.Bounds(), &image.Uniform{On}, image.ZP, draw.Src)
			})

			Context("Config.BlinkAndPWMSets has not been set", func() {
				It("should error", func() {
					Expect(as.SetFrame(0, frame)).To(MatchError("must set Config.BlinkAndPWMSets first"))
					Expect(writeBuf.Contents()).To(BeEmpty())
				})
			})

			Context("Config.BlinkAndPWMSets set to 1", func() {
				BeforeEach(func() {
					as.blinkAndPWMSets = 1
				})

				It("should write first frame", func() {
					Expect(as.SetFrame(1, frame)).To(Succeed())
					TestRemainingCommands(writeBuf, commandsPerFrame)
				})

				It("should write last frame", func() {
					Expect(as.SetFrame(36, frame)).To(Succeed())
					TestRemainingCommands(writeBuf, commandsPerFrame)
				})

				It("should error on zero indexed frame", func() {
					Expect(as.SetFrame(0, frame)).To(MatchError("frame out of range [1,36]: 0"))
					Expect(writeBuf.Contents()).To(BeEmpty())
				})

				It("should error on too high frame", func() {
					Expect(as.SetFrame(37, frame)).To(MatchError("frame out of range [1,36]: 37"))
					Expect(writeBuf.Contents()).To(BeEmpty())
				})
			})

			Context("Config.BlinkAndPWMSets set to 6", func() {
				BeforeEach(func() {
					as.blinkAndPWMSets = 6
				})

				DescribeTable("PWM set",
					func(pwmSet int, secondSegmentData string) {
						frame.SetPWMSet(uint8(pwmSet))
						Expect(as.SetFrame(1, frame)).To(Succeed())
						Expect(
							fmt.Sprintf("%08b", writeBuf.Contents()[7]),
						).To(
							Equal(secondSegmentData),
						)

						TestRemainingCommands(writeBuf, commandsPerFrame)
					},
					Entry("uses set 1 when default 0", 0, "00000111"),
					Entry("uses set 1", 1, "00000111"),
					Entry("uses set 2", 2, "00100111"),
					Entry("uses set 3", 3, "01000111"),
					Entry("uses set 4", 4, "01100111"),
					Entry("uses set 5", 5, "10000111"),
					Entry("uses set 6", 6, "10100111"),
				)

				It("should error on too high PWM set", func() {
					frame.SetPWMSet(7)
					Expect(as.SetFrame(1, frame)).To(MatchError("PWM set out of range [1,6]: 7"))
					Expect(writeBuf.Contents()).To(BeEmpty())
				})
			})
		})

		Describe("SetBlinkAndPWMSet", func() {
			const (
				blinkSegments = int(FrameSegmentLast-FrameSegmentFirst) + 1
				pwmSegments   = int(PWMSegmentLast-PWMSegmentFirst) + 1
				totalSegments = blinkSegments + pwmSegments
			)

			var blink, pwm *Frame12x11

			BeforeEach(func() {
				blink = NewFrame12x11()
				pwm = NewFrame12x11()
				draw.Draw(pwm, pwm.Bounds(), &image.Uniform{On}, image.ZP, draw.Src)
			})

			Context("Config.BlinkAndPWMSets has not been set", func() {
				It("should error", func() {
					Expect(as.SetBlinkAndPWMSet(0, blink, pwm)).To(MatchError("must set Config.BlinkAndPWMSets first"))
					Expect(writeBuf.Contents()).To(BeEmpty())
				})
			})

			Context("Config.BlinkAndPWMSets set to 6", func() {
				BeforeEach(func() {
					as.blinkAndPWMSets = 6
				})

				It("should write first set", func() {
					Expect(as.SetBlinkAndPWMSet(1, blink, pwm)).To(Succeed())
					TestRemainingCommands(writeBuf, totalSegments)
				})

				It("should write last set", func() {
					Expect(as.SetBlinkAndPWMSet(6, blink, pwm)).To(Succeed())
					TestRemainingCommands(writeBuf, totalSegments)
				})

				It("should error on zero indexed frame", func() {
					Expect(as.SetBlinkAndPWMSet(0, blink, pwm)).To(MatchError("set out of range [1,6]: 0"))
					Expect(writeBuf.Contents()).To(BeEmpty())
				})

				It("should error on too high frame", func() {
					Expect(as.SetBlinkAndPWMSet(7, blink, pwm)).To(MatchError("set out of range [1,6]: 7"))
					Expect(writeBuf.Contents()).To(BeEmpty())
				})
			})
		})
	})
})
