package as1130

import (
	"fmt"

	"golang.org/x/exp/io/i2c"

	"github.com/dcarley/as1130/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("as1130", func() {
	TestCommand := func(buf *gbytes.Buffer, register, subregister byte, binaryData string) {
		command := buf.Contents()

		Expect(command).To(HaveLen(4), "shoud have four bytes of commands")
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

		Describe("SetDisplayOption", func() {
			const (
				register    = RegisterControl
				subregister = ControlDisplayOption
			)

			It("should write defaults", func() {
				option := DisplayOption{
					ScanLimit: 1,
				}
				Expect(as.SetDisplayOption(option)).To(Succeed())
				TestCommand(writeBuf, register, subregister, "11100000")
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

			It("should write defaults", func() {
				config := Config{}
				Expect(as.SetConfig(config)).To(Succeed())
				TestCommand(writeBuf, register, subregister, "00000001")
			})

			It("should write non-defaults", func() {
				config := Config{
					LowVDDReset:         true,
					LowVDDStatus:        true,
					LEDErrorCorrection:  true,
					DotCorrection:       true,
					CommonAddress:       true,
					MemoryConfiguration: 6,
				}
				Expect(as.SetConfig(config)).To(Succeed())
				TestCommand(writeBuf, register, subregister, "11111110")
			})

			It("should error on out of range MemoryConfiguration", func() {
				config := Config{
					MemoryConfiguration: 7,
				}
				Expect(as.SetConfig(config)).To(MatchError("MemoryConfiguration out of range [1,6]: 7"))
				Expect(writeBuf.Contents()).To(BeEmpty())
			})
		})
	})
})
