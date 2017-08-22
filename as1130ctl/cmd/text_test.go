package cmd_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Text", func() {
	It("should print help message", func() {
		command := exec.Command(cliPath, "text", "--help")
		session, err := Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(Exit(0))
	})

	It("should reject size that isn't 24x5", func() {
		command := exec.Command(cliPath, "text", "foo", "-s", "12x11")
		session, err := Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())
		Eventually(session).Should(Exit(1))
		Eventually(session.Err).Should(gbytes.Say(`only 24x5 size is supported for this command`))
	})
})
