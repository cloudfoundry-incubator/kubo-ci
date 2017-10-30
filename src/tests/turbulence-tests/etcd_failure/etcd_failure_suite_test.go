package etcd_failure_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"os"
	"testing"
)

var iaas string

func TestEtcdFailure(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EtcdFailure Suite")
}

var _ = BeforeSuite(func() {
	iaas = os.Getenv("TURBULENCE_IAAS")
	platforms := []string{"aws", "gcp", "vsphere", "openstack"}
	message := fmt.Sprintf("Expected TURBULENCE_IAAS to be one of the following values: %#v", platforms)
	Expect(platforms).To(ContainElement(iaas), message)

})
