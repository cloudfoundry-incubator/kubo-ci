package persistence_failure_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var iaas string

func TestPersistenceFailure(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PersistenceFailure Suite")
}

var _ = BeforeSuite(func() {
	iaas = os.Getenv("TURBULENCE_IAAS")
	platforms := []string{"aws", "gcp"}
	message := fmt.Sprintf("Expected TURBULENCE_IAAS to be one of the following values: %#v", platforms)
	Expect(platforms).To(ContainElement(iaas), message)
})
