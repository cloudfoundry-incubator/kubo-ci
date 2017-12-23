package test_helpers

import (
	"errors"
	"fmt"
	testconfig "tests/config"

	"encoding/json"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cppforlife/turbulence/client"
	. "github.com/onsi/gomega"
)

func TurbulenceClient(testconfig testconfig.Turbulence) client.Turbulence {
	config := client.Config{
		Host: testconfig.Host,
		Port: testconfig.Port,

		Username: testconfig.Username,
		Password: testconfig.Password,

		CACert: testconfig.CaCert,
	}
	clientLogger := logger.NewLogger(logger.LevelNone)
	return client.NewFactory(clientLogger).New(config)
}

func AllBoshWorkersHaveJoinedK8s(deployment director.Deployment, kubectl *KubectlRunner) bool {
	Eventually(func() []director.VMInfo {
		return DeploymentVmsOfType(deployment, WorkerVmType, VmRunningState)
	}, "600s", "30s").Should(HaveLen(3))

	Eventually(func() []string { return GetReadyNodes(GetNodes(kubectl)) }, "240s", "5s").Should(HaveLen(3))
	return true
}

func BoshIdByIp(deployment director.Deployment, externalIp string) (string, error) {
	vms, err := deployment.VMInfos()
	Expect(err).NotTo(HaveOccurred())
	for _, vm := range vms {
		for _, ip := range vm.IPs {
			if ip == externalIp {
				return vm.ID, nil
			}
		}
	}
	return "", errors.New(fmt.Sprintf("Can't find vm id with ip %s", externalIp))
}

func GetReadyNodes(nodes []Node) []string {
	readyNodes := []string{}
	for _, node := range nodes {
		for _, condition := range node.Status.Conditions {
			if condition.ConditionType == "Ready" && condition.Status == "True" {
				readyNodes = append(readyNodes, node.Metadata.Name)
				break
			}
		}
	}
	return readyNodes
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

func GetNodes(kubectl *KubectlRunner) []Node {
	nodes := NodesArray{}
	bytes := kubectl.GetOutputBytes("get", "nodes", "-o", "json")
	err := json.Unmarshal(bytes, &nodes)
	Expect(err).ToNot(HaveOccurred())
	return nodes.Items
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
