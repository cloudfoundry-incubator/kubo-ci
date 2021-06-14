package generic_test

import (
	"os"
	"strconv"
	"tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

const defaultHPATimeout = "10m"

var (
	hpaDeployment = test_helpers.PathFromRoot("specs/hpa-php-apache.yml")
	loadGenerator = test_helpers.PathFromRoot("specs/load-generator.yml")
)

var _ = Describe("Horizontal Pod Autoscaling", func() {
	BeforeEach(func() {
		createHPADeployment()
	})

	AfterEach(func() {
		kubectl.StartKubectlCommand("delete", "-f", hpaDeployment).Wait(kubectl.TimeoutInSeconds)
		kubectl.StartKubectlCommand("delete", "pods", "--all").Wait(kubectl.TimeoutInSeconds)
	})

	It("scales the pods accordingly", func() {
		HPATimeout := os.Getenv("HPA_TIMEOUT")
		if HPATimeout == "" {
			HPATimeout = defaultHPATimeout
		}

		Eventually(getNumberOfPods, HPATimeout, "5s").Should(Equal(1))
		By("creating more pods when the CPU load increases")

		increaseCPULoad()
		Eventually(getNumberOfPods, HPATimeout, "5s").Should(BeNumerically(">", 1))

		By("decreasing the number of pods when the CPU load decreases")

		kubectl.StartKubectlCommand("delete", "pod/load-generator", "--now").Wait(kubectl.TimeoutInSeconds / 2)

		Eventually(getNumberOfPods, HPATimeout, "5s").Should(BeNumerically("==", 1))
	})
})

func getNumberOfPods() int {
	session := kubectl.StartKubectlCommand("get", "hpa/php-apache", "-o", "jsonpath={.status.currentReplicas}")
	Eventually(session, "20s").Should(gexec.Exit())
	if session.ExitCode() != 0 {
		return 0
	}
	replicas, _ := strconv.Atoi(string(session.Out.Contents()))
	return replicas
}

func createHPADeployment() {
	session := kubectl.StartKubectlCommand("apply", "-f", hpaDeployment)
	Eventually(session, "10s").Should(gexec.Exit(0))

	Eventually(func() string {
		return kubectl.GetPodStatusBySelector(kubectl.Namespace(), "app=php-apache")
	}, kubectl.TimeoutInSeconds*2).Should(Equal("Running"))
}

func increaseCPULoad() {
	session := kubectl.StartKubectlCommand("create", "-f", loadGenerator)
	Eventually(session, "10s").Should(gexec.Exit(0))

	Eventually(func() string {
		return kubectl.GetPodStatus(kubectl.Namespace(), "load-generator")
	}, kubectl.TimeoutInSeconds*2).Should(Equal("Running"))
}
