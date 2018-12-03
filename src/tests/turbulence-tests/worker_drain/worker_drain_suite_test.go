package worker_drain

import (
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestWorkerDrain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WorkerDrain Suite")
}

var _ = BeforeSuite(func() {
	director := NewDirector()
	director.EnableResurrection(true)
})
