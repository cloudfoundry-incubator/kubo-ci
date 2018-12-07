package test_helpers

import (
	"errors"
	"fmt"

	"encoding/json"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cppforlife/turbulence/client"
	. "github.com/onsi/gomega"
)

func TurbulenceClient() client.Turbulence {
	config := client.NewConfigFromEnv()
	clientLogger := logger.NewLogger(logger.LevelNone)
	return client.NewFactory(clientLogger).New(config)
}

func AllBoshWorkersHaveJoinedK8s(deployment director.Deployment, kubectl *KubectlRunner) bool {
	workerCount := len(DeploymentVmsOfType(deployment, WorkerVMType, ""))
	Eventually(func() []director.VMInfo {
		return DeploymentVmsOfType(deployment, WorkerVMType, VMRunningState)
	}, "600s", "30s").Should(HaveLen(workerCount))

	Eventually(GetReadyNodes, "240s", "5s").Should(HaveLen(workerCount))
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

func AllComponentsAreHealthy(kubectl *KubectlRunner) bool {
	components, err := GetComponentStatusOrError(kubectl)
	if err != nil {
		return false
	}

	if len(components) == 0 {
		return false
	}

	for _, component := range components {
		if component.Conditions[0].Status != "True" {
			return false
		}
	}
	return true
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

func GetComponentStatusOrError(kubectl *KubectlRunner) ([]ComponentStatus, error) {
	response := ComponentStatusResponse{}
	bytes, err := kubectl.GetOutputBytesOrError("get", "componentstatus", "-o", "json")
	if err != nil {
		return []ComponentStatus{}, err
	}
	err = json.Unmarshal(bytes, &response)
	return response.Items, err
}

func GetNodeNamesForRunningPods(kubectl *KubectlRunner) []string {
	output, err := kubectl.GetOutput("get", "pods", "-o", "jsonpath='{.items[*].spec.nodeName}'")
	Expect(err).ToNot(HaveOccurred())
	return output
}

func GetNewVmId(oldVms []director.VMInfo, newVmIds []string) (string, error) {
	Expect(len(oldVms)).NotTo(BeNumerically("<", 3))
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
