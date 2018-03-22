package multiaz

import (
	"testing"
	"tests/config"
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestK8sMultiAZ(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "K8s Multi-AZ Suite")
}

var (
	runner     *test_helpers.KubectlRunner
	nginxSpec  = test_helpers.PathFromRoot("specs/nginx-daemonset.yml")
	testconfig *config.Config
)

var _ = BeforeSuite(func() {
	var err error
	testconfig, err = config.InitConfig()
	Expect(err).NotTo(HaveOccurred())

	runner = test_helpers.NewKubectlRunner(testconfig.Kubernetes.PathToKubeConfig)
	runner.RunKubectlCommand("create", "namespace", runner.Namespace()).Wait("60s")
})

var _ = AfterSuite(func() {
	if runner != nil && !CurrentGinkgoTestDescription().Failed {
		runner.RunKubectlCommand("delete", "namespace", runner.Namespace()).Wait("60s")
	}
})

func MultiAZDescribe(description string, callback func()) bool {
	return Describe("[multiaz]", func() {
		BeforeEach(func() {
			if !testconfig.IntegrationTests.IncludeMultiAZ {
				Skip(`Skipping this test suite because Config.IntegrationTests.IncludeMultiAZ is set to 'false'.`)
			}
		})
		Describe(description, callback)
	})
}
