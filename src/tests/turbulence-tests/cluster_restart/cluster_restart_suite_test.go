package cluster_restart_test

import (
	"testing"
	"tests/config"

	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestClusterRestart(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ClusterRestart Suite")
}

var (
	testconfig *config.Config
)
var _ = BeforeSuite(func() {
	var err error
	testconfig, err = config.InitConfig()
	Expect(err).NotTo(HaveOccurred())
	director := NewDirector(testconfig.Bosh)
	director.EnableResurrection(true)
})
