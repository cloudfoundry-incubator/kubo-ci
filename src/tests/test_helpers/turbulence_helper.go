package test_helpers

import (
	"errors"
	"fmt"
	"os/exec"

	"encoding/json"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cppforlife/turbulence/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
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

	Eventually(func() []string { return GetNodes(kubectl) }).Should(HaveLen(3))
	return true
}

func ExpectAllComponentsToBeHealthy(kubectl *KubectlRunner) {
	components := GetComponentStatus(kubectl)
	Expect(components).ToNot(BeEmpty())
	for _, component := range components {
		Expect(component.Conditions[0].Status).To(Equal("True"))
	}
}

type Condition struct {
	ConditionType string `json:"type"`
	Status        string `json:"status"`
}

type Node struct {
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
	Status struct {
		Conditions []Condition `json:"conditions"`
	} `json:"status"`
}

type NodesArray struct {
	Items []Node `json:"items"`
}

type ComponentStatus struct {
	Conditions []Condition `json:"conditions"`
}

type ComponentStatusResponse struct {
	Items []ComponentStatus `json:"items"`
}

func GetNodes(kubectl *KubectlRunner) []string {
	nodes := NodesArray{}
	bytes := kubectl.GetOutputBytes("get", "nodes", "-o", "json")
	err := json.Unmarshal(bytes, &nodes)
	Expect(err).ToNot(HaveOccurred())

	names := []string{}
	for _, item := range nodes.Items {
		names = append(names, item.Metadata.Name)
	}
	return names
}


func GetComponentStatus(kubectl *KubectlRunner) []ComponentStatus {
	response := ComponentStatusResponse{}
	bytes := kubectl.GetOutputBytes("get", "componentstatus", "-o", "json")
	fmt.Println(string(bytes))
	err := json.Unmarshal(bytes, &response)
	Expect(err).ToNot(HaveOccurred())
	return response.Items
}

func GetNodeNamesForRunningPods(kubectl *KubectlRunner) []string {
	return kubectl.GetOutput("get", "pods", "-o", "jsonpath='{.items[*].spec.nodeName}'")
}

func NewVmId(oldVms []director.VMInfo, newVmIds []string) (string, error) {
	oldVmIds := []string{oldVms[1].VMID, oldVms[2].VMID}
	fmt.Printf("oldVmIds: %s", oldVmIds)
	for _, vmId := range newVmIds {
		if !contains(oldVmIds, vmId) {
			return vmId, nil
		}
	}
	return "", errors.New("no new VM found")
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
	if iaas == "vsphere" {
		cid := vms[0].IPs[0]
		KillVMById(cid, iaas)
	} else {
		cid := vms[0].VMID
		KillVMById(cid, iaas)
	}

}

func KillVMById(iaasSpecificVmIdentifier string, iaas string) {
	var cmd *exec.Cmd

	switch iaas {
	case "gcp":
		cmd = exec.Command("gcloud", "-q", "compute", "instances", "delete", iaasSpecificVmIdentifier)
		break
	case "aws":
		cmd = exec.Command("aws", "ec2", "terminate-instances", "--instance-ids", iaasSpecificVmIdentifier)
		break
	case "vsphere":
		cmd = exec.Command("govc", "vm.power", "-off=true", "-vm.ip", iaasSpecificVmIdentifier)
		break
	case "openstack":
		cmd = exec.Command("openstack", "server", "delete", iaasSpecificVmIdentifier)
	default:
		Fail(fmt.Sprintf("Unsupported IaaS: %s", iaas))
	}

	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 300, 20).Should(gexec.Exit(0))
}
