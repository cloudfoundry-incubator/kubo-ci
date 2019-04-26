package test_helpers

import (
	"fmt"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"

	cmdconf "github.com/cloudfoundry/bosh-cli/cmd/config"
	boshdir "github.com/cloudfoundry/bosh-cli/director"
	boshuaa "github.com/cloudfoundry/bosh-cli/uaa"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

const (
	WorkerVMType   = "worker"
	MasterVMType   = "master"
	VMRunningState = "running"
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
	environmentURL, err := url.Parse(os.Getenv("BOSH_ENVIRONMENT"))
	Expect(err).NotTo(HaveOccurred())
	if environmentURL.Scheme == "" {
		environmentURL.Scheme = "https"
	}
	hostname := environmentURL.Hostname()
	if hostname == "" {
		hostname = os.Getenv("BOSH_ENVIRONMENT")
	}
	return fmt.Sprintf("%s://%s:8443", environmentURL.Scheme, hostname)
}

func buildUAA() (boshuaa.UAA, error) {
	logger := boshlog.NewLogger(boshlog.LevelError)
	factory := boshuaa.NewFactory(logger)

	config, err := boshuaa.NewConfigFromURL(getUAAUrl())
	if err != nil {
		return nil, err
	}

	config.Client = os.Getenv("BOSH_CLIENT")
	config.ClientSecret = os.Getenv("BOSH_CLIENT_SECRET")
	config.CACert = os.Getenv("BOSH_CA_CERT")

	return factory.New(config)
}

func buildDirector(uaa boshuaa.UAA) (boshdir.Director, error) {
	logger := boshlog.NewWriterLogger(boshlog.LevelInfo, GinkgoWriter)
	factory := boshdir.NewFactory(logger)

	config, err := boshdir.NewConfigFromURL(os.Getenv("BOSH_ENVIRONMENT"))
	if err != nil {
		return nil, err
	}

	config.CACert = os.Getenv("BOSH_CA_CERT")
	config.TokenFunc = boshuaa.NewClientTokenSession(uaa).TokenFunc

	return factory.New(config, cmdconf.FSConfig{}, boshdir.NewNoopTaskReporter(), boshdir.NewNoopFileReporter())
}
