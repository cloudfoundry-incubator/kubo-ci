package addons_test

import (
	"tests/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAddons(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Addons Suite")
}

var testconfig *config.Config

var _ = BeforeSuite(func() {
	var err error
	testconfig, err = config.InitConfig()
	Expect(err).NotTo(HaveOccurred())
})

func AddonsDescribe(description string, callback func()) bool {
	return Describe("[addons]", func() {
		BeforeEach(func() {
			if !testconfig.IntegrationTests.IncludeAddons {
				Skip(`Skipping this test suite because Config.IntegrationTests.IncludeAddons is set to 'false'.`)
			}
		})
		Describe(description, callback)
	})
}
