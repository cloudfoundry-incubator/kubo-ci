package etcd_test

import (
	"fmt"
	"tests/test_helpers"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Etcd cert on worker", func() {
	var (
		deployment boshdir.Deployment
		workers    []boshdir.VMInfo
		masters    []boshdir.VMInfo
		master     boshdir.VMInfo
		directory  string
		err        error
		director   boshdir.Director
	)

	director = test_helpers.NewDirector()
	deployment, err = director.FindDeployment(deploymentName)
	Expect(err).NotTo(HaveOccurred())
	workers = test_helpers.DeploymentVmsOfType(deployment, test_helpers.WorkerVMType, "")

	masters = test_helpers.DeploymentVmsOfType(deployment, test_helpers.MasterVMType, "")
	Expect(len(masters) > 0).To(BeTrue())
	master = masters[0]

	Context("For directorys under /coreos.com/network/", func() {
		directory = "/coreos.com/network/"

		AfterEach(func() {
			for _, vm := range workers {
				args := []string{"rm", fmt.Sprintf("%s%s", directory, vm.ID)}
				value := test_helpers.RunEtcdCommandFromMasterWithFullPrivilege(deploymentName, master.ID, args...)
				Expect(value).NotTo(ContainSubstring("Insufficient credentials"))
			}
		})
		It("should have read access ", func() {
			args := []string{"ls", directory}
			for _, vm := range workers {
				value := test_helpers.RunEtcdCommandFromWorker(deploymentName, vm.ID, args...)
				Expect(value).NotTo(ContainSubstring("Insufficient credentials"))
			}

		})

		It("should have write access", func() {
			for _, vm := range workers {
				args := []string{"set", fmt.Sprintf("%s%s", directory, vm.ID)}
				value := test_helpers.RunEtcdCommandFromWorker(deploymentName, vm.ID, args...)
				Expect(value).NotTo(ContainSubstring("Insufficient credentials"))
			}

		})
	})

	Context("For directorys under /", func() {
		directory = "/"

		AfterEach(func() {
			for _, vm := range workers {
				args := []string{"rm", fmt.Sprintf("%s%s", directory, vm.ID)}
				value := test_helpers.RunEtcdCommandFromMasterWithFullPrivilege(deploymentName, master.ID, args...)
				Expect(value).NotTo(ContainSubstring("Insufficient credentials"))
			}
		})

		It("should not have read access", func() {
			for _, vm := range workers {
				args := []string{"ls", directory}
				value := test_helpers.RunEtcdCommandFromWorker(deploymentName, vm.ID, args...)
				Expect(value).To(ContainSubstring("Insufficient credentials"))
			}
		})

		It("should not have write access", func() {
			for _, vm := range workers {
				args := []string{"set", fmt.Sprintf("%s%s", directory, vm.ID)}
				value := test_helpers.RunEtcdCommandFromWorker(deploymentName, vm.ID, args...)
				Expect(value).To(ContainSubstring("Insufficient credentials"))
			}
		})
	})
})
