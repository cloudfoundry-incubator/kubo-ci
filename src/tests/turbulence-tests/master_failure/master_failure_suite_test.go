package master_failure_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"os"
)

var iaas string

func TestMasterFailure(t *testing.T) {
	RegisterFailHandler(Fail)
	iaas = os.Getenv("TURBULENCE_IAAS")
	RunSpecs(t, "Master failure suite")
}
