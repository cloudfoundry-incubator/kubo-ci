package test_helpers

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DeploySmorgasbord(kubectl *KubectlRunner, iaas string) {
	storageClassSpec := PathFromRoot(fmt.Sprintf("specs/storage-class-%s.yml", iaas))
	smorgasbordSpec := PathFromRoot("specs/smorgasbord.yml")

	Eventually(kubectl.RunKubectlCommand("apply", "-f", storageClassSpec), kubectl.TimeoutInSeconds*2).Should(gexec.Exit(0))
	Eventually(kubectl.RunKubectlCommand("apply", "-f", smorgasbordSpec), kubectl.TimeoutInSeconds*2).Should(gexec.Exit(0))
	Eventually(kubectl.RunKubectlCommand("rollout", "status", "daemonset/fluentd-elasticsearch", "-w"), "900s").Should(gexec.Exit(0))
	WaitForPodsToRun(kubectl, kubectl.TimeoutInSeconds*5)
}

func WaitForPodsToRun(kubectl *KubectlRunner, timeout float64) {
	waitForPods(kubectl, "status.phase!=Running,status.phase!=Succeeded", timeout)
}

func WaitForPodsToDie(kubectl *KubectlRunner, timeout float64) {
	waitForPods(kubectl, "status.phase!=Succeeded", timeout)
}

func waitForPods(kubectl *KubectlRunner, selector string, timeout float64) {
	Eventually(func() bool {
		clientset, err := NewKubeClient()
		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
			return false
		}
		pods, err := clientset.CoreV1().Pods(kubectl.Namespace()).List(meta_v1.ListOptions{
			FieldSelector: selector,
		})
		if err != nil {
			GinkgoWriter.Write([]byte(err.Error()))
			return false
		}
		for _, pod := range pods.Items {
			fmt.Fprintf(GinkgoWriter, "Pod name:%s, pod status: %s, Events:\n", pod.Name, pod.Status.Phase)
			events, err := clientset.CoreV1().Events(kubectl.Namespace()).List(meta_v1.ListOptions{
				FieldSelector: fmt.Sprintf("involvedObject.kind=Pod,involvedObject.name=%s", pod.Name),
			})
			if err != nil {
				fmt.Fprintf(GinkgoWriter, "\tFailed to list events for pod: %s\n\terr: %s\n", pod.Name, err.Error())
			} else {
				for _, event := range events.Items {
					fmt.Fprintf(GinkgoWriter, "\t%s: %s\n", event.Reason, event.Message)
				}
			}
		}
		return len(pods.Items) == 0
	}, timeout, "5s").Should(BeTrue())
}

func DeleteSmorgasbord(kubectl *KubectlRunner, iaas string) {
	storageClassSpec := PathFromRoot(fmt.Sprintf("specs/storage-class-%s.yml", iaas))
	smorgasbordSpec := PathFromRoot("specs/smorgasbord.yml")

	Eventually(kubectl.RunKubectlCommand("delete", "-f", smorgasbordSpec), kubectl.TimeoutInSeconds*2).Should(gexec.Exit(0))
	Eventually(kubectl.RunKubectlCommand("delete", "--all", "pvc"), kubectl.TimeoutInSeconds*2).Should(gexec.Exit(0))
	Eventually(kubectl.RunKubectlCommand("delete", "-f", storageClassSpec), kubectl.TimeoutInSeconds*2).Should(gexec.Exit(0))
	WaitForPodsToDie(kubectl, kubectl.TimeoutInSeconds*5)
}
