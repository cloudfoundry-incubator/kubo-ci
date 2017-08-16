package workload_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
	"testing"
	"fmt"
)

var iaas string

func TestWorkerFailure(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WorkerFailure Suite")
}

var _ = BeforeSuite(func() {
	iaas = os.Getenv("TURBULENCE_IAAS")
	platforms := []string{"aws", "gcp"}
	message := fmt.Sprintf("Expected TURBULENCE_IAAS to be one of the following values: %#v", platforms)
	Expect(platforms).To(ContainElement(iaas), message)
})