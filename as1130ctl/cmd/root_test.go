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
})
