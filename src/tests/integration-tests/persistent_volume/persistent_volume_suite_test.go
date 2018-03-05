package persistent_volume_test

import (
	"tests/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"testing"
)

var iaas string
var deploymentName string
var testconfig *config.Config

func TestPersistentVolume(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PersistentVolume Suite")
}

var _ = BeforeSuite(func() {
	var err error
	testconfig, err = config.InitConfig()
	Expect(err).NotTo(HaveOccurred())

	platforms := []string{"aws", "gcp", "vsphere", "openstack"}
	message := fmt.Sprintf("Expected IAAS to be one of the following values: %#v", platforms)
	Expect(platforms).To(ContainElement(testconfig.Bosh.Iaas), message)
})

func PersistentVolumeDescribe(description string, callback func()) bool {
	return Describe("[PersistentVolume]", func() {
		BeforeEach(func() {
			if !testconfig.TestSuites.IncludePersistentVolume {
				Skip(`Skipping this test suite because Config.TestSuites.IncludePersistentVolume is set to 'false'.`)
			}
		})
		Describe(description, callback)
	})
}
