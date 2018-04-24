package certificates_test

import (
	"tests/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCertificates(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Certificates Suite")
}

var testconfig *config.Config

var _ = BeforeSuite(func() {
	var err error
	testconfig, err = config.InitConfig()
	Expect(err).NotTo(HaveOccurred())
})

func CertificatesDescribe(description string, callback func()) bool {
	return Describe("[certificates]", func() {
		BeforeEach(func() {
			if !testconfig.IntegrationTests.IncludeCertificates {
				Skip(`Skipping this test suite because Config.IntegrationTests.IncludeCertificates is set to 'false'.`)
			}
		})
		Describe(description, callback)
	})
}
