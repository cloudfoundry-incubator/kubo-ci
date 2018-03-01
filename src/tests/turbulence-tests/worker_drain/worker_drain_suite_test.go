package worker_drain

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var iaas string

func TestWorkerDrain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WorkerDrain Suite")
}
