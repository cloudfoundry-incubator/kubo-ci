package generic_test

import (
	"fmt"
	"strconv"
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var defaultHPATimeout = "210s"

var _ = Describe("Horizontal Pod Autoscaling", func() {
	It("scales the pods accordingly", func() {
		HPATimeout := MustHaveEnv("HPA_TIMEOUT")
		if HPATimeout == "" {
			HPATimeout = defaultHPATimeout
		}

		runHPAPod()
		createHPA()

		By("creating more pods when the CPU load increases")

		increaseCPULoad()
		Eventually(func() int {
			session := runner.RunKubectlCommand("get", "hpa/php-apache", "-o", "jsonpath={.status.currentReplicas}")
			Eventually(session, "10s").Should(gexec.Exit(0))
			replicas, _ := strconv.Atoi(string(session.Out.Contents()))
			return replicas
		}, HPATimeout).Should(BeNumerically(">", 1))

		By("decreasing the number of pods when the CPU load decreases")

		session := runner.RunKubectlCommand("delete", "deployment/load-generator")
		Eventually(session, "10s").Should(gexec.Exit(0))

		Eventually(func() int {
			session := runner.RunKubectlCommand("get", "hpa/php-apache", "-o", "jsonpath={.status.currentReplicas}")
			Eventually(session, "10s").Should(gexec.Exit(0))
			replicas, _ := strconv.Atoi(string(session.Out.Contents()))
			return replicas
		}, HPATimeout).Should(BeNumerically("==", 1))
	})
})

func runHPAPod() {
	session := runner.RunKubectlCommand("run", "php-apache", "--image=k8s.gcr.io/hpa-example", "--requests=cpu=200m", "--expose", "--port=80")
	Eventually(session, "10s").Should(gexec.Exit(0))

	Eventually(func() string {
		return runner.GetPodStatusBySelector(runner.Namespace(), "run=php-apache")
	}, "120s").Should(Equal("Running"))
}

func createHPA() {
	session := runner.RunKubectlCommand("autoscale", "deployment/php-apache", "--cpu-percent=25", "--min=1", "--max=2")
	Eventually(session, "10s").Should(gexec.Exit(0))
}

func increaseCPULoad() {
	remoteCommand := fmt.Sprintf("while true; do wget -q -O- http://php-apache.%s.svc.cluster.local; done", runner.Namespace())

	session := runner.RunKubectlCommand("run", "-i", "--tty", "load-generator", "--image=busybox", "--", "/bin/sh", "-c", remoteCommand)
	Eventually(session, "10s").Should(gexec.Exit(0))

	Eventually(func() string {
		return runner.GetPodStatusBySelector(runner.Namespace(), "run=load-generator")
	}, "120s").Should(Equal("Running"))
}
