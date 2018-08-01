package workload_test

import (
	"tests/config"
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	iaas       string
	testconfig *config.Config
)

func TestWorkerFailure(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WorkerFailure Suite")
}

var _ = BeforeSuite(func() {
	var err error
	testconfig, err = config.InitConfig()
	Expect(err).NotTo(HaveOccurred())
	director := NewDirector(testconfig.Bosh)
	director.EnableResurrection(true)

})

func WorkerFailureDescribe(description string, callback func()) bool {
	return Describe("[persistence_failure]", func() {
		BeforeEach(func() {
			if !testconfig.TurbulenceTests.IncludeWorkerFailure {
				Skip(`Skipping this test suite because Config.TurbulenceTests.IncludeWorkerFailure is set to 'false'.`)
			}
		})
		Describe(description, callback)
	})
}
