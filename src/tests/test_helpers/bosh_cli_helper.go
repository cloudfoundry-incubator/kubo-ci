package test_helpers

import (
	"fmt"

	. "github.com/onsi/gomega"

	testconfig "tests/config"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
	boshuaa "github.com/cloudfoundry/bosh-cli/uaa"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

const (
	WorkerVmType   = "worker"
	MasterVmType   = "master"
	VmRunningState = "running"
)

func CountDeploymentVmsOfType(deployment boshdir.Deployment, jobName, processState string) func() int {
	return func() int {
		return len(DeploymentVmsOfType(deployment, jobName, processState))
	}
}

func CountProcessesOnVmsOfType(deployment boshdir.Deployment, jobName, processName, processState string) func() int {
	return func() int {
		return len(ProcessesOnVmsOfType(deployment, jobName, processName, processState))
	}
}

func DeploymentVmsOfType(deployment boshdir.Deployment, jobName, processState string) []boshdir.VMInfo {
	vms, err := deployment.VMInfos()
	Expect(err).NotTo(HaveOccurred())
	return VmsMatchingPredicate(vms, func(vmInfo boshdir.VMInfo) bool {
		return vmInfo.JobName == jobName && vmInfo.ProcessState == processState
	})
}

func ProcessesOnVmsOfType(deployment boshdir.Deployment, jobName, processName, processState string) []boshdir.VMInfo {
	vms := DeploymentVmsOfType(deployment, jobName, VmRunningState)
	return VmsMatchingPredicate(vms, func(vmInfo boshdir.VMInfo) bool {
		for _, process := range vmInfo.Processes {
			if process.Name == processName && process.State == processState {
				return true
			}
		}
		return false
	})
}

func VmsMatchingPredicate(vms []boshdir.VMInfo, f func(boshdir.VMInfo) bool) []boshdir.VMInfo {
	result := make([]boshdir.VMInfo, 0)
	for _, vmInfo := range vms {
		if f(vmInfo) {
			result = append(result, vmInfo)
		}
	}
	return result
}

func GetWorkerIP(deployment boshdir.Deployment) string {
	vms := DeploymentVmsOfType(deployment, WorkerVmType, VmRunningState)
	return vms[0].IPs[0]
}

func GetMasterIP(deployment boshdir.Deployment) string {
	vms := DeploymentVmsOfType(deployment, MasterVmType, VmRunningState)
	return vms[0].IPs[0]
}

func NewDirector(testconfig testconfig.Bosh) boshdir.Director {
	uaa, err := buildUAA(testconfig)
	if err != nil {
		panic(err)
	}

	director, err := buildDirector(uaa, testconfig)
	if err != nil {
		panic(err)
	}

	return director
}

func buildUaaUrl(testconfig testconfig.Bosh) string {
	return buildDirectorUrl(testconfig) + ":8443"
}

func buildDirectorUrl(testconfig testconfig.Bosh) string {
	return fmt.Sprintf("https://%s", testconfig.Environment)
}

func buildUAA(testconfig testconfig.Bosh) (boshuaa.UAA, error) {
	logger := boshlog.NewLogger(boshlog.LevelError)
	factory := boshuaa.NewFactory(logger)

	config, err := boshuaa.NewConfigFromURL(buildUaaUrl(testconfig))
	if err != nil {
		return nil, err
	}

	config.Client = testconfig.Client
	config.ClientSecret = testconfig.ClientSecret

	config.CACert = testconfig.CaCert

	return factory.New(config)
}

func buildDirector(uaa boshuaa.UAA, testconfig testconfig.Bosh) (boshdir.Director, error) {
	logger := boshlog.NewLogger(boshlog.LevelError)
	factory := boshdir.NewFactory(logger)

	config, err := boshdir.NewConfigFromURL(buildDirectorUrl(testconfig))
	if err != nil {
		return nil, err
	}

	config.CACert = testconfig.CaCert

	config.TokenFunc = boshuaa.NewClientTokenSession(uaa).TokenFunc

	return factory.New(config, boshdir.NewNoopTaskReporter(), boshdir.NewNoopFileReporter())
}
