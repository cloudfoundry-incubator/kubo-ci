package cluster_restart_test

import (
	"testing"

	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestClusterRestart(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ClusterRestart Suite")
}

var _ = BeforeSuite(func() {
	director := NewDirector()
	director.EnableResurrection(true)
})
