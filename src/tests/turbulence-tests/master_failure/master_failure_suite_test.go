package master_failure_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var iaas string

func TestMasterFailure(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Master failure suite")
}
