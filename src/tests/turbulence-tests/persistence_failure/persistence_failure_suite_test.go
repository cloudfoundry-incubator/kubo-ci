package persistence_failure_test

import (
	"tests/config"
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	testconfig *config.Config
)

func TestPersistenceFailure(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PersistenceFailure Suite")
}

var _ = BeforeSuite(func() {
	var err error
	testconfig, err = config.InitConfig()
	Expect(err).NotTo(HaveOccurred())
	director := NewDirector(testconfig.Bosh)
	director.EnableResurrection(true)
})

func PersistenceFailureDescribe(description string, callback func()) bool {
	return Describe("[persistence_failure]", func() {
		BeforeEach(func() {
			if !testconfig.TurbulenceTests.IncludePersistenceFailure {
				Skip(`Skipping this test suite because Config.TurbulenceTests.IncludePersistenceFailure is set to 'false'.`)
			}
		})
		Describe(description, callback)
	})
}
