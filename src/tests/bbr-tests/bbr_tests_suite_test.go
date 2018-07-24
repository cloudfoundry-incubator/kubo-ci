package bbr_tests_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBbrTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BbrTests Suite")
}
