package persistence_failure_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPersistenceFailure(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PersistenceFailure Suite")
}
