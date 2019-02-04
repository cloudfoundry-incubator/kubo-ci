package windows_test

import (
	"testing"
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var hasWindowsWorkers bool

var _ = BeforeSuite(func() {
	var err error
	hasWindowsWorkers, err = test_helpers.HasWindowsWorkers()
	Expect(err).To(BeNil())
})

func TestWindows(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Windows Suite")
}
