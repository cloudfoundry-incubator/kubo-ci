package workload_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var iaas string

func TestWorkerFailure(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WorkerFailure Suite")
}
