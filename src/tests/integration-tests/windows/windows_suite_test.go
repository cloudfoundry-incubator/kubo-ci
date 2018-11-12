package windows_test

import (
	"testing"
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = BeforeSuite(func() {
	hasWindowsWorkers, err := test_helpers.HasWindowsWorkers()
	Expect(err).To(BeNil())
	if !hasWindowsWorkers {
		Skip("skipping Windows tests since no Windows nodes were detected")
	}
})

func TestWindows(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Windows Suite")
}
