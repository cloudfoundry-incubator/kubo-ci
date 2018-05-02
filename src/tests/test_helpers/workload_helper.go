package test_helpers

import (
	"fmt"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func DeploySmorgasbord(kubectl *KubectlRunner, iaas string) {
	storageClassSpec := PathFromRoot(fmt.Sprintf("specs/storage-class-%s.yml", iaas))
	smorgasbordSpec := PathFromRoot("specs/smorgasbord.yml")

	Eventually(kubectl.RunKubectlCommand("apply", "-f", storageClassSpec), "120s").Should(gexec.Exit(0))
	Eventually(kubectl.RunKubectlCommand("apply", "-f", smorgasbordSpec), "120s").Should(gexec.Exit(0))
	Eventually(kubectl.RunKubectlCommand("rollout", "status", "daemonset/fluentd-elasticsearch", "-w"), "900s").Should(gexec.Exit(0))
	CheckSmorgasbord(kubectl, "5m")
}

func CheckSmorgasbord(kubectl *KubectlRunner, timeout string) {
	Eventually(func() int {
		output := kubectl.GetOutput("get", "pods", "--field-selector", "status.phase!=Running,status.phase!=Succeeded")
		return len(output)
	}, timeout, "5s").Should(Equal(0))

}

func DeleteSmorgasbord(kubectl *KubectlRunner, iaas string) {
	storageClassSpec := PathFromRoot(fmt.Sprintf("specs/storage-class-%s.yml", iaas))
	smorgasbordSpec := PathFromRoot("specs/smorgasbord.yml")

	Eventually(kubectl.RunKubectlCommand("delete", "-f", smorgasbordSpec), "120s").Should(gexec.Exit(0))
	Eventually(kubectl.RunKubectlCommand("delete", "-f", storageClassSpec), "120s").Should(gexec.Exit(0))
}
