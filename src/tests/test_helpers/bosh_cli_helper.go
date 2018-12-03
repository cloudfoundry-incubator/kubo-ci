package test_helpers

import (
	"fmt"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
	boshuaa "github.com/cloudfoundry/bosh-cli/uaa"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

const (
	WorkerVMType    = "worker"
	MasterVMType    = "master"
	VMRunningState  = "running"
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
	vms, err := deployment.InstanceInfos()
	Expect(err).NotTo(HaveOccurred())
	fmt.Fprintf(GinkgoWriter, "Bosh vms for deployment %s: \n", deployment.Name())
	return VmsMatchingPredicate(vms, func(vmInfo boshdir.VMInfo) bool {
		fmt.Fprintf(GinkgoWriter, "%s/%s - %s\n", vmInfo.JobName, vmInfo.ID, vmInfo.ProcessState)
		return vmInfo.JobName == jobName && (processState == "" || vmInfo.ProcessState == processState)
	})
}

func ProcessesOnVmsOfType(deployment boshdir.Deployment, jobName, processName, processState string) []boshdir.VMInfo {
	vms := DeploymentVmsOfType(deployment, jobName, VMRunningState)
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

func getUAAUrl() string {
	environmentURL, err := url.Parse(getBoshEnvironment())
	Expect(err).NotTo(HaveOccurred())

	return fmt.Sprintf("%s://%s:8443", environmentURL.Scheme, environmentURL.Hostname())
}

func getBoshEnvironment() string {
	return MustHaveEnv("BOSH_ENVIRONMENT")
}

func getBoshClientSecret() string {
	return MustHaveEnv("BOSH_CLIENT_SECRET")
}

func getBoshClient() string {
	return MustHaveEnv("BOSH_CLIENT")
}

func getBoshCACert() string {
	return MustHaveEnv("BOSH_CA_CERT")
}

func GetBoshDeployment() string {
	return MustHaveEnv("BOSH_DEPLOYMENT")
}

func GetIaas() string {
	return MustHaveEnv("IAAS")
}

func buildUAA() (boshuaa.UAA, error) {
	logger := boshlog.NewLogger(boshlog.LevelError)
	factory := boshuaa.NewFactory(logger)

	config, err := boshuaa.NewConfigFromURL(getUAAUrl())
	if err != nil {
		return nil, err
	}

	config.Client = getBoshClient()
	config.ClientSecret = getBoshClientSecret()

	config.CACert = getBoshCACert()

	return factory.New(config)
}

func buildDirector(uaa boshuaa.UAA) (boshdir.Director, error) {
	logger := boshlog.NewWriterLogger(boshlog.LevelInfo, GinkgoWriter, GinkgoWriter)
	factory := boshdir.NewFactory(logger)

	config, err := boshdir.NewConfigFromURL(getBoshEnvironment())
	if err != nil {
		return nil, err
	}

	config.CACert = getBoshCACert()

	config.TokenFunc = boshuaa.NewClientTokenSession(uaa).TokenFunc

	return factory.New(config, boshdir.NewNoopTaskReporter(), boshdir.NewNoopFileReporter())
}
