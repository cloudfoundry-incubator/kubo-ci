package persistent_volume_test

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

			Eventually(kubectl.RunKubectlCommand("create", "-f", storageClassSpec), "60s").Should(gexec.Exit(0))
			Eventually(kubectl.RunKubectlCommand("create", "-f", nfsServerSpec), "60s").Should(gexec.Exit(0))
			Eventually(kubectl.RunKubectlCommand("create", "-f", nfsServerServiceSpec), "60s").Should(gexec.Exit(0))

			k8s, err := NewKubeClient()
			Expect(err).NotTo(HaveOccurred())

			NFSService, err := k8s.CoreV1().Services(kubectl.Namespace()).Get("nfs", meta_v1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())

			nfsPvSpec = templateNFSSpec(NFSService.Spec.ClusterIP, PathFromRoot("specs/nfs-pv.yml"))
			Eventually(kubectl.RunKubectlCommand("create", "-f", nfsPvSpec), "60s").Should(gexec.Exit(0))

			Eventually(kubectl.RunKubectlCommand("create", "-f", nfsPvcSpec), "60s").Should(gexec.Exit(0))
			Eventually(kubectl.RunKubectlCommand("create", "-f", nfsPodRcSpec), "60s").Should(gexec.Exit(0))
		})

		AfterEach(func() {
			Eventually(kubectl.RunKubectlCommand("delete", "-f", nfsPodRcSpec), "60s").Should(gexec.Exit())
			Eventually(kubectl.RunKubectlCommand("delete", "-f", nfsPvcSpec), "60s").Should(gexec.Exit())
			Eventually(kubectl.RunKubectlCommand("delete", "-f", nfsPvSpec), "60s").Should(gexec.Exit())
			Eventually(kubectl.RunKubectlCommand("delete", "-f", nfsServerServiceSpec), "60s").Should(gexec.Exit())
			Eventually(kubectl.RunKubectlCommand("delete", "-f", nfsServerSpec), "60s").Should(gexec.Exit())

			// Some pv(c)s aren't being cleaned
			Eventually(kubectl.RunKubectlCommand("delete", "pvc", "--all"), "60s").Should(gexec.Exit())
			Eventually(kubectl.RunKubectlCommand("delete", "pv", "--all"), "60s").Should(gexec.Exit())
			Eventually(kubectl.RunKubectlCommand("delete", "-f", storageClassSpec), "60s").Should(gexec.Exit())
		})

		It("should mount an NFS PV to a workload", func() {
			WaitForPodsToRun(kubectl, "120s")
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
