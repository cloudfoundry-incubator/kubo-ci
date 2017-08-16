package test_helpers

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cppforlife/turbulence/client"
	"github.com/onsi/gomega/gexec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	WorkerVmType   = "worker"
	VmRunningState = "running"
)

func TurbulenceClient() client.Turbulence {
	config := client.NewConfigFromEnv()
	clientLogger := logger.NewLogger(logger.LevelNone)
	return client.NewFactory(clientLogger).New(config)
}

func AllBoshWorkersHaveJoinedK8s(deployment director.Deployment, kubectl *KubectlRunner) bool {
	Eventually(func() []director.VMInfo {
		return DeploymentVmsOfType(deployment, WorkerVmType, VmRunningState)
	}, "600s", "30s").Should(HaveLen(3))

	Eventually(func() []string { return GetNodes(kubectl) }, "120s", "5s").Should(HaveLen(3))
	return true
}

func GetNodes(kubectl *KubectlRunner) []string {
	getNodesSession := kubectl.RunKubectlCommand("get", "nodes", "-o", "name")
	Eventually(getNodesSession, "10s").Should(gexec.Exit(0))
	output := getNodesSession.Out.Contents()
	return strings.Fields(string(output))
}

func GetNodeNamesForRunningPods(kubectl *KubectlRunner) []string {
	getPodsSession := kubectl.RunKubectlCommand("get", "pods", "-o", "jsonpath='{.items[*].spec.nodeName}'")
	Eventually(getPodsSession).Should(gexec.Exit(0))
	var output []byte
	output = getPodsSession.Out.Contents()
	return strings.Fields(string(output))
}

func NewVmId(oldVms []director.VMInfo, newVmIds []string) (string, error) {
	oldVmIds := []string{oldVms[1].VMID, oldVms[2].VMID}
	fmt.Printf("oldVmIds:", oldVmIds)
	for _, vmId := range newVmIds {
		if !contains(oldVmIds, vmId) {
			return vmId, nil
		}
	}
	return "", errors.New("No new VM found!")
}

func contains(vmNames []string, vmName string) bool {
	for _, element := range vmNames {
		if element == vmName {
			return true
		}
	}
	return false
}

func KillVM(vms []director.VMInfo, iaas string) {
	cid := vms[0].VMID
	var cmd *exec.Cmd

	switch iaas {
	case "gcp":
		cmd = exec.Command("gcloud", "-q", "compute", "instances", "delete", cid)
		break
	case "aws":
		cmd = exec.Command("aws", "ec2", "terminate-instances", "--instance-ids", cid)
		break
	default:
		Fail(fmt.Sprintf("Unsupported IaaS: %s", iaas))
	}

	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 300, 20).Should(gexec.Exit(0))
}
