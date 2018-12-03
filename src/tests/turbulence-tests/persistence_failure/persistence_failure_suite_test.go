package persistence_failure_test

import (
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPersistenceFailure(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PersistenceFailure Suite")
}

var _ = BeforeSuite(func() {
	director := NewDirector()
	director.EnableResurrection(true)
})
