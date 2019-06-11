package volume_test

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gexec"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("NFS", func() {
	var (
		kubectl *KubectlRunner
		iaas    string
	)

	BeforeEach(func() {
		kubectl = NewKubectlRunner()
		kubectl.Setup()

		var err error
		iaas, err = IaaS()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		kubectl.Teardown()
	})

	Context("when creating an NFS PV", func() {
		var (
			storageClassSpec     string
			nfsServerSpec        string
			nfsServerServiceSpec string
			nfsPvSpec            string
			nfsPvcSpec           string
			nfsPodRcSpec         string
		)

		BeforeEach(func() {
			storageClassSpec = PathFromRoot(fmt.Sprintf("specs/storage-class-%s.yml", iaas))
			nfsServerSpec = PathFromRoot("specs/nfs-server-statefulset.yml")
			nfsServerServiceSpec = PathFromRoot("specs/nfs-server-service.yml")
			nfsPvcSpec = PathFromRoot("specs/nfs-pvc.yml")
			nfsPodRcSpec = PathFromRoot("specs/nfs-pod-rc.yml")

			Eventually(kubectl.StartKubectlCommand("create", "-f", storageClassSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
			Eventually(kubectl.StartKubectlCommand("create", "-f", nfsServerSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
			Eventually(kubectl.StartKubectlCommand("create", "-f", nfsServerServiceSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit(0))

			k8s, err := NewKubeClient()
			Expect(err).NotTo(HaveOccurred())

			NFSService, err := k8s.CoreV1().Services(kubectl.Namespace()).Get("nfs", meta_v1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())

			nfsPvSpec = templateNFSSpec(NFSService.Spec.ClusterIP, PathFromRoot("specs/nfs-pv.yml"))
			Eventually(kubectl.StartKubectlCommand("create", "-f", nfsPvSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit(0))

			Eventually(kubectl.StartKubectlCommand("create", "-f", nfsPvcSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
			Eventually(kubectl.StartKubectlCommand("create", "-f", nfsPodRcSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit(0))
		})

		AfterEach(func() {
			Eventually(kubectl.StartKubectlCommand("delete", "-f", nfsPodRcSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit())
			Eventually(kubectl.StartKubectlCommand("delete", "-f", nfsPvcSpec), kubectl.TimeoutInSeconds*2).Should(gexec.Exit())
			Eventually(kubectl.StartKubectlCommand("delete", "-f", nfsPvSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit())
			Eventually(kubectl.StartKubectlCommand("delete", "-f", nfsServerServiceSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit())
			Eventually(kubectl.StartKubectlCommand("delete", "-f", nfsServerSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit())

			// Some pv(c)s aren't being cleaned
			Eventually(kubectl.StartKubectlCommand("delete", "pvc", "--all"), kubectl.TimeoutInSeconds).Should(gexec.Exit())
			Eventually(kubectl.StartKubectlCommand("delete", "pv", "--all"), kubectl.TimeoutInSeconds).Should(gexec.Exit())
			Eventually(kubectl.StartKubectlCommand("delete", "-f", storageClassSpec), kubectl.TimeoutInSeconds).Should(gexec.Exit())
		})

		It("should mount an NFS PV to a workload", func() {
			WaitForPodsToRun(kubectl, kubectl.TimeoutInSeconds*10)
		})
	})
})

func templateNFSSpec(serviceIP string, spec string) string {
	t, err := template.ParseFiles(spec)
	Expect(err).NotTo(HaveOccurred())

	f, err := ioutil.TempFile("", filepath.Base(spec))
	Expect(err).NotTo(HaveOccurred())
	defer f.Close()

	type templateInfo struct{ NFSServerIP string }
	Expect(t.Execute(f, templateInfo{NFSServerIP: serviceIP})).To(Succeed())

	return f.Name()
}
