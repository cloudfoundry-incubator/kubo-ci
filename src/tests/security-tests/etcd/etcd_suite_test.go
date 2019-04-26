package etcd_test

import (
	"fmt"
	"os"
	"testing"
	"tests/test_helpers"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEtcd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Etcd Suite")
}

var (
	deploymentName string
	workers        []boshdir.VMInfo
	masters        []boshdir.VMInfo
	master         boshdir.VMInfo
	director       boshdir.Director
)

var _ = BeforeSuite(func() {
	test_helpers.CheckRequiredEnvs([]string{
		"DEPLOYMENT_NAME",
		"BOSH_ENVIRONMENT",
		"BOSH_CLIENT",
		"BOSH_CLIENT_SECRET",
		"BOSH_CA_CERT"})

	deploymentName = os.Getenv("DEPLOYMENT_NAME")

	director = test_helpers.NewDirector()
	deployment, err := director.FindDeployment(deploymentName)
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "Failed getting deployment %s: %v", deploymentName, err)
		os.Exit(1)
	}

	workers = test_helpers.DeploymentVmsOfType(deployment, test_helpers.WorkerVMType, "")

	masters = test_helpers.DeploymentVmsOfType(deployment, test_helpers.MasterVMType, "")
	Expect(len(masters) > 0).To(BeTrue())
	master = masters[0]

})
