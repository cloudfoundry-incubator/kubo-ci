package generic_test

import (
	"tests/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGeneric(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Generic Suite")
}

var testconfig *config.Config

var _ = BeforeSuite(func() {
	var err error
	testconfig, err = config.InitConfig()
	Expect(err).NotTo(HaveOccurred())
})

func GenericDescribe(description string, callback func()) bool {
	return Describe("[generic]", func() {
		BeforeEach(func() {
			if !testconfig.IntegrationTests.IncludeGeneric {
				Skip(`Skipping this test suite because Config.IntegrationTests.IncludeGeneric is set to 'false'.`)
			}
		})
		Describe(description, callback)
	})
}
