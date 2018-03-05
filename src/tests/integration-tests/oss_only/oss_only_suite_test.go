package oss_only_test

import (
	"tests/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestOssOnly(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OssOnly Suite")
}

var testconfig *config.Config

var _ = BeforeSuite(func() {
	var err error
	testconfig, err = config.InitConfig()
	Expect(err).NotTo(HaveOccurred())
})

func OSSOnlyDescribe(description string, callback func()) bool {
	return Describe("[oss_only]", func() {
		BeforeEach(func() {
			if !testconfig.TestSuites.IncludeOSSOnly {
				Skip(`Skipping this test suite because Config.TestSuites.IncludeOSSOnly is set to 'false'.`)
			}
		})
		Describe(description, callback)
	})
}
