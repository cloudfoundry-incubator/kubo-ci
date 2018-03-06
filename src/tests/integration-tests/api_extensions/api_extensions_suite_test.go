package api_extensions_test

import (
	"tests/config"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestApiExtensions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ApiExtensions Suite")
}

var testconfig *config.Config

var _ = BeforeSuite(func() {
	var err error

	SetDefaultEventuallyTimeout(60 * time.Second)

	testconfig, err = config.InitConfig()
	Expect(err).NotTo(HaveOccurred())

})

func APIExtensionsDescribe(description string, callback func()) bool {
	return Describe("[api_extensions]", func() {
		BeforeEach(func() {
			if !testconfig.IntegrationTests.IncludeAPIExtensions {
				Skip(`Skipping this test suite because Config.IntegrationTests.IncludAPIExtensions is set to 'false'.`)
			}
		})
		Describe(description, callback)
	})
}
