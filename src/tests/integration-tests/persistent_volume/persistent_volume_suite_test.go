package persistent_volume_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"os"
	"testing"
)

var iaas string
var deploymentName string

func TestPersistentVolume(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PersistentVolume Suite")
}

var _ = BeforeSuite(func() {
	iaas = os.Getenv("INTEGRATIONTEST_IAAS")
	deploymentName = os.Getenv("DEPLOYMENT_NAME")
	platforms := []string{"aws", "gcp", "vsphere"}
	message := fmt.Sprintf("Expected IAAS to be one of the following values: %#v", platforms)
	Expect(platforms).To(ContainElement(iaas), message)
})
