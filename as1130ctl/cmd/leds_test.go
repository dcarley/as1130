package cmd_test

import (
	"os/exec"
	"regexp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Leds", func() {
	It("should print help message", func() {
		command := exec.Command(cliPath, "leds", "--help")
		session, err := Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(Exit(0))
	})

	It("should reject invalid LED positions", func() {
		command := exec.Command(cliPath, "leds",
			"1,1",   // valid
			"0,0",   // zero indexed
			"1",     // no y
			"1,",    // no y
			"1,2,3", // too many axis
			"a,b",   // NaN
			"25,1",  // x too large
			"1,6",   // y too large
			"24,5",  // valid
		)
		session, err := Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(Exit(1))
		Eventually(session.Err).Should(gbytes.Say(regexp.QuoteMeta(
			`invalid LED positions: 0,0 1 1, 1,2,3 a,b 25,1 1,6`,
		)))
	})
})
