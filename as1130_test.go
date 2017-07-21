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
	})
})
