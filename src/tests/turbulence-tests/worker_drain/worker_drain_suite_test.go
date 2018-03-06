package worker_drain

import (
	"tests/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	iaas       string
	testconfig *config.Config
)

func TestWorkerDrain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WorkerDrain Suite")
}

var _ = BeforeSuite(func() {
	var err error
	testconfig, err = config.InitConfig()
	Expect(err).NotTo(HaveOccurred())
})

func WorkerDrainDescribe(description string, callback func()) bool {
	return Describe("[persistence_failure]", func() {
		BeforeEach(func() {
			if !testconfig.TurbulenceTests.IncludeWorkerDrain {
				Skip(`Skipping this test suite because Config.TurbulenceTests.IncludeWorkerDrain is set to 'false'.`)
			}
		})
		Describe(description, callback)
	})
}
