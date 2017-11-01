package etcd_failure_test

import (
	. "tests/test_helpers"

	"github.com/cloudfoundry/bosh-cli/director"
	"github.com/cppforlife/turbulence/incident"
	"github.com/cppforlife/turbulence/incident/selector"
	"github.com/cppforlife/turbulence/tasks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const testKey = "foo"
const testValue = "bar"

var _ = Describe("Etcd failure scenarios", func() {
	var deployment director.Deployment
	var countRunningEtcd func() int
	var kubectl *KubectlRunner
	var etcdNodeIP string

	BeforeEach(func() {
		var err error

		director := NewDirector()
		deployment, err = director.FindDeployment(DeploymentName)
		Expect(err).NotTo(HaveOccurred())
		countRunningEtcd = CountDeploymentVmsOfType(deployment, EtcdVmType, VmRunningState)

		kubectl = NewKubectlRunner()
		kubectl.CreateNamespace()

		Expect(countRunningEtcd()).To(Equal(3))
		Expect(AllEtcdHaveJoinedK8s(deployment, kubectl)).To(BeTrue())
	})

	AfterEach(func() {
		kubectl.RunKubectlCommand("delete", "namespace", kubectl.Namespace())
		DeleteKeyFromEtcd(etcdNodeIP, testKey)
	})

	Specify("Etcd nodes rejoin the cluster and contain up-to-date data", func() {

		By("Writing data to the Etcd leader")
		etcdNodeIP = GetEtcdIP(deployment)
		PutKeyToEtcd(etcdNodeIP, testKey, testValue)

		By("Deleting the Etcd VM")
		turbulenceClient := TurbulenceClient()
		killOneEtcd := incident.Request{
			Selector: selector.Request{
				Deployment: &selector.NameRequest{
					Name: DeploymentName,
				},
				Group: &selector.NameRequest{
					Name: EtcdVmType,
				},
				ID: &selector.IDRequest{
					Limit: selector.MustNewLimitFromString("1"),
				},
			},
			Tasks: tasks.OptionsSlice{
				tasks.KillOptions{},
			},
		}
		incident := turbulenceClient.CreateIncident(killOneEtcd)

		By("Killing VM")
		incident.Wait()

		By("Waiting for Bosh to recognize dead VMs")
		Eventually(countRunningEtcd, 600, 20).Should(Equal(2))

		By("Waiting for resurrection")
		Eventually(countRunningEtcd, 600, 20).Should(Equal(3))

		By("Verifying that the Etcd VM has joined the K8s cluster")
		Eventually(func() bool { return AllEtcdHaveJoinedK8s(deployment, kubectl) }, 600, 20).Should(BeTrue())

		By("Reading the data from the Etcd cluster")
		Expect(GetKeyFromEtcd(etcdNodeIP, testKey)).To(Equal(testValue))
	})

})
