package master_failure_test

import (
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)


func TestMasterFailure(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Master failure suite")
}

var _ = BeforeSuite(func() {
	director := NewDirector()
	director.EnableResurrection(true)
})
