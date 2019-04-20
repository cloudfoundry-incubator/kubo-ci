package etcd_test

import (
	"os"
	"testing"
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEtcd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Etcd Suite")
}

var (
	deploymentName string
)

var _ = BeforeSuite(func() {
	test_helpers.CheckRequiredEnvs([]string{
		"BOSH_DEPLOYMENT",
		"BOSH_ENVIRONMENT",
		"BOSH_CLIENT",
		"BOSH_CLIENT_SECRET",
		"BOSH_CA_CERT"})

	deploymentName = os.Getenv("BOSH_DEPLOYMENT")
})
