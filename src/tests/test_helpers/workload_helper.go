package test_helpers

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeploySmorgasbord(kubectl *KubectlRunner, iaas string) {
	storageClassSpec := PathFromRoot(fmt.Sprintf("specs/storage-class-%s.yml", iaas))
	smorgasbordSpec := PathFromRoot("specs/smorgasbord.yml")

	Eventually(kubectl.RunKubectlCommand("apply", "-f", storageClassSpec), "120s").Should(gexec.Exit(0))
	Eventually(kubectl.RunKubectlCommand("apply", "-f", smorgasbordSpec), "120s").Should(gexec.Exit(0))
	Eventually(kubectl.RunKubectlCommand("rollout", "status", "daemonset/fluentd-elasticsearch", "-w"), "900s").Should(gexec.Exit(0))
	WaitForPodsToRun(kubectl, "5m")
}

func WaitForPodsToRun(kubectl *KubectlRunner, timeout string) {
	Eventually(func() bool {
		clientset, err := NewKubeClient()
		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
			return false
		}
		pods, err := clientset.CoreV1().Pods(kubectl.Namespace()).List(v1.ListOptions{
			FieldSelector: "status.phase!=Running,status.phase!=Succeeded",
		})
		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
			return false
		}
		return len(pods.Items) == 0
	}, timeout, "5s").Should(BeTrue())
}

func DeleteSmorgasbord(kubectl *KubectlRunner, iaas string) {
	storageClassSpec := PathFromRoot(fmt.Sprintf("specs/storage-class-%s.yml", iaas))
	smorgasbordSpec := PathFromRoot("specs/smorgasbord.yml")

	Eventually(kubectl.RunKubectlCommand("delete", "-f", smorgasbordSpec), "120s").Should(gexec.Exit(0))
	Eventually(kubectl.RunKubectlCommand("delete", "-f", storageClassSpec), "120s").Should(gexec.Exit(0))
}
