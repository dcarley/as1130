package cmd_test

import (
	"os/exec"
	"regexp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Root", func() {
	It("should print help message", func() {
		command := exec.Command(cliPath, "help")
		session, err := Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(Exit(0))
	})

	Describe("device argument", func() {
		It("should print default", func() {
			command := exec.Command(cliPath, "help")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(gbytes.Say(regexp.QuoteMeta(
				`I2C device address (default 0x30)`,
			)))
		})
	})

	Describe("address argument", func() {
		It("should print default", func() {
			command := exec.Command(cliPath, "help")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(gbytes.Say(regexp.QuoteMeta(
				`I2C device path (default "/dev/i2c-1")`,
			)))
		})
	})

	Describe("size argument", func() {
		It("should print default", func() {
			command := exec.Command(cliPath, "help")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(gbytes.Say(regexp.QuoteMeta(
				`size of the frame: 24x5, 12x11 (default 24x5)`,
			)))
		})

		It("should accept valid size argument", func() {
			command := exec.Command(cliPath, "-s", "12x11")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))
		})

		It("should reject invalid size argument", func() {
			command := exec.Command(cliPath, "-s", "1x1")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(gbytes.Say(regexp.QuoteMeta(
				`invalid argument "1x1" for "-s, --size" flag: invalid size`,
			)))
		})
	})
})
