package as1130_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAs1130(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "As1130 Suite")
}
