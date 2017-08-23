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

	Describe("layout", func() {
		It("should print layout for default size", func() {
			command := exec.Command(cliPath, "leds", "-l")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(gbytes.Say(regexp.QuoteMeta(
				` | 1,1 | 2,1 | 3,1 | 4,1 | 5,1 | 6,1 | 7,1 | 8,1 | 9,1 | 10,1 | 11,1 | 12,1 | 13,1 | 14,1 | 15,1 | 16,1 | 17,1 | 18,1 | 19,1 | 20,1 | 21,1 | 22,1 | 23,1 | 24,1 |
 | 1,2 | 2,2 | 3,2 | 4,2 | 5,2 | 6,2 | 7,2 | 8,2 | 9,2 | 10,2 | 11,2 | 12,2 | 13,2 | 14,2 | 15,2 | 16,2 | 17,2 | 18,2 | 19,2 | 20,2 | 21,2 | 22,2 | 23,2 | 24,2 |
 | 1,3 | 2,3 | 3,3 | 4,3 | 5,3 | 6,3 | 7,3 | 8,3 | 9,3 | 10,3 | 11,3 | 12,3 | 13,3 | 14,3 | 15,3 | 16,3 | 17,3 | 18,3 | 19,3 | 20,3 | 21,3 | 22,3 | 23,3 | 24,3 |
 | 1,4 | 2,4 | 3,4 | 4,4 | 5,4 | 6,4 | 7,4 | 8,4 | 9,4 | 10,4 | 11,4 | 12,4 | 13,4 | 14,4 | 15,4 | 16,4 | 17,4 | 18,4 | 19,4 | 20,4 | 21,4 | 22,4 | 23,4 | 24,4 |
 | 1,5 | 2,5 | 3,5 | 4,5 | 5,5 | 6,5 | 7,5 | 8,5 | 9,5 | 10,5 | 11,5 | 12,5 | 13,5 | 14,5 | 15,5 | 16,5 | 17,5 | 18,5 | 19,5 | 20,5 | 21,5 | 22,5 | 23,5 | 24,5 |`,
			)))
		})

		It("should print layout for 12x11", func() {
			command := exec.Command(cliPath, "leds", "-l", "-s", "12x11")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))
			Eventually(session.Out).Should(gbytes.Say(regexp.QuoteMeta(
				` | 1,1  | 2,1  | 3,1  | 4,1  | 5,1  | 6,1  | 7,1  | 8,1  | 9,1  | 10,1  | 11,1  | 12,1  |
 | 1,2  | 2,2  | 3,2  | 4,2  | 5,2  | 6,2  | 7,2  | 8,2  | 9,2  | 10,2  | 11,2  | 12,2  |
 | 1,3  | 2,3  | 3,3  | 4,3  | 5,3  | 6,3  | 7,3  | 8,3  | 9,3  | 10,3  | 11,3  | 12,3  |
 | 1,4  | 2,4  | 3,4  | 4,4  | 5,4  | 6,4  | 7,4  | 8,4  | 9,4  | 10,4  | 11,4  | 12,4  |
 | 1,5  | 2,5  | 3,5  | 4,5  | 5,5  | 6,5  | 7,5  | 8,5  | 9,5  | 10,5  | 11,5  | 12,5  |
 | 1,6  | 2,6  | 3,6  | 4,6  | 5,6  | 6,6  | 7,6  | 8,6  | 9,6  | 10,6  | 11,6  | 12,6  |
 | 1,7  | 2,7  | 3,7  | 4,7  | 5,7  | 6,7  | 7,7  | 8,7  | 9,7  | 10,7  | 11,7  | 12,7  |
 | 1,8  | 2,8  | 3,8  | 4,8  | 5,8  | 6,8  | 7,8  | 8,8  | 9,8  | 10,8  | 11,8  | 12,8  |
 | 1,9  | 2,9  | 3,9  | 4,9  | 5,9  | 6,9  | 7,9  | 8,9  | 9,9  | 10,9  | 11,9  | 12,9  |
 | 1,10 | 2,10 | 3,10 | 4,10 | 5,10 | 6,10 | 7,10 | 8,10 | 9,10 | 10,10 | 11,10 | 12,10 |
 | 1,11 | 2,11 | 3,11 | 4,11 | 5,11 | 6,11 | 7,11 | 8,11 | 9,11 | 10,11 | 11,11 | 12,11 |`,
			)))
		})
	})
})
