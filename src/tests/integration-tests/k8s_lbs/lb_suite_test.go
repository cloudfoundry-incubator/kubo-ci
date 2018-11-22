package k8s_lbs_test

import (
	"fmt"
	"testing"
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestK8sLb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "K8sLb Suite")
}

var (
	kubectl     *test_helpers.KubectlRunner
	nginxLBSpec = test_helpers.PathFromRoot("specs/nginx-lb.yml")

	internalNginxLBSpec string
	iaas                string
)

var _ = BeforeSuite(func() {
	kubectl = test_helpers.NewKubectlRunner()
	kubectl.Setup()

	var err error
	iaas, err = test_helpers.IaaS()
	Expect(err).ToNot(HaveOccurred())
	internalNginxLBSpec = fmt.Sprintf(test_helpers.PathFromRoot("specs/nginx-internal-lb-%s.yml"), iaas)
})

var _ = AfterSuite(func() {
	if kubectl != nil {
		kubectl.Teardown()
	}
})
