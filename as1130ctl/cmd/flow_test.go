package cmd_test

import (
	"os/exec"
	"regexp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Flow", func() {
	It("should print help message", func() {
		command := exec.Command(cliPath, "flow", "--help")
		session, err := Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(Exit(0))
	})

	Describe("frames", func() {
		It("should reject larger than 5", func() {
			command := exec.Command(cliPath, "flow", "-w", "6")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(gbytes.Say(regexp.QuoteMeta(
				`width cannot be greater than 5: 6`,
			)))
		})
	})
})
