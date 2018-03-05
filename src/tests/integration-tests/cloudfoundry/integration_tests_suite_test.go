package cloudfoundry_test

import (
	"testing"
	"tests/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIntegrationTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IntegrationTests Suite")
}

var testconfig *config.Config

var _ = BeforeSuite(func() {
	var err error
	testconfig, err = config.InitConfig()
	Expect(err).NotTo(HaveOccurred())
})

func CloudFoundryDescribe(description string, callback func()) bool {
	return Describe("[cloudfoundry]", func() {
		BeforeEach(func() {
			if !testconfig.TestSuites.IncludeCloudFoundry {
				Skip(`Skipping this test suite because Config.TestSuites.IncludeCloudFoundry is set to 'false'.`)
			}
		})
		Describe(description, callback)
	})
}
