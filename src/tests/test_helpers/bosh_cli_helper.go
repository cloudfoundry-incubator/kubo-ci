package test_helpers

import (
	"fmt"
	"os"

	. "github.com/onsi/gomega"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
	boshuaa "github.com/cloudfoundry/bosh-cli/uaa"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

const (
	WorkerVmType   = "worker"
	MasterVmType   = "master"
	EtcdVmType     = "etcd"
	VmRunningState = "running"
	DeploymentName = "ci-service"
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

func GetEtcdIP(deployment boshdir.Deployment) string {
	vms := DeploymentVmsOfType(deployment, EtcdVmType, VmRunningState)
	return vms[0].IPs[0]
}

func NewDirector() boshdir.Director {
	uaa, err := buildUAA()
	if err != nil {
		panic(err)
	}

	director, err := buildDirector(uaa)
	if err != nil {
		panic(err)
	}

	return director
}

func buildUaaUrl() string {
	return buildDirectorUrl() + ":8443"
}

func buildDirectorUrl() string {
	return fmt.Sprintf("https://%s", os.Getenv("BOSH_ENVIRONMENT"))
}

func buildUAA() (boshuaa.UAA, error) {
	logger := boshlog.NewLogger(boshlog.LevelError)
	factory := boshuaa.NewFactory(logger)

	config, err := boshuaa.NewConfigFromURL(buildUaaUrl())
	if err != nil {
		return nil, err
	}

	config.Client = "bosh_admin"
	config.ClientSecret = os.Getenv("BOSH_CLIENT_SECRET")

	config.CACert = os.Getenv("BOSH_CA_CERT")

	return factory.New(config)
}

func buildDirector(uaa boshuaa.UAA) (boshdir.Director, error) {
	logger := boshlog.NewLogger(boshlog.LevelError)
	factory := boshdir.NewFactory(logger)

	config, err := boshdir.NewConfigFromURL(buildDirectorUrl())
	if err != nil {
		return nil, err
	}

	config.CACert = os.Getenv("BOSH_CA_CERT")

	config.TokenFunc = boshuaa.NewClientTokenSession(uaa).TokenFunc

	return factory.New(config, boshdir.NewNoopTaskReporter(), boshdir.NewNoopFileReporter())
}
