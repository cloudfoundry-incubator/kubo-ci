package bbr_tests_test

import (
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gexec"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	tappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

var _ = Describe("BackupAndRestore", func() {
	var (
		bbrDir         string
		boshCACertFile string
		err            error
		k8s            kubernetes.Interface
		namespace      string
		deploymentApi  tappsv1.DeploymentInterface

		deploymentName string
		bbrArgs        []string
	)

	BeforeEach(func() {
		bbrDir, err = ioutil.TempDir("", "")
		Expect(err).ToNot(HaveOccurred())
		boshCACertFile = bbrDir + "/bosh-ca.pem"
		err = ioutil.WriteFile(boshCACertFile, []byte(MustHaveEnv("BOSH_CA_CERT")), 0600)
		Expect(err).ToNot(HaveOccurred())

		k8s, err = NewKubeClient()
		Expect(err).ToNot(HaveOccurred())

		nsObject, err := CreateTestNamespace(k8s, "bbr")
		Expect(err).ToNot(HaveOccurred())
		namespace = nsObject.Name

		deploymentApi = k8s.AppsV1().Deployments(namespace)

		deploymentName = MustHaveEnv("BOSH_DEPLOYMENT")
		bbrArgs = []string{"deployment",
			"--target", MustHaveEnv("BOSH_ENVIRONMENT"),
			"--username", MustHaveEnv("BOSH_CLIENT"),
			"--password", MustHaveEnv("BOSH_CLIENT_SECRET"),
			"--deployment", deploymentName,
			"--ca-cert", boshCACertFile}
	})

	AfterEach(func() {
		k8s.CoreV1().Namespaces().Delete(namespace, &metav1.DeleteOptions{})
	})

	It("should backup and restore successfully", func() {
		var (
			nginx1Deployment *v1.Deployment
			nginx2Deployment *v1.Deployment
			nginx3Deployment *v1.Deployment
		)

		By("Deploying workload 1 and 2", func() {
			nginx1Deployment = NewDeployment("nginx-1", GetNginxDeploymentSpec())
			nginx1Deployment, err = deploymentApi.Create(nginx1Deployment)
			Expect(err).ToNot(HaveOccurred())

			nginx2Deployment = NewDeployment("nginx-2", GetNginxDeploymentSpec())
			nginx2Deployment, err = deploymentApi.Create(nginx2Deployment)
			Expect(err).ToNot(HaveOccurred())

			err = WaitForDeployment(deploymentApi, namespace, nginx1Deployment.Name)
			Expect(err).NotTo(HaveOccurred())
			err = WaitForDeployment(deploymentApi, namespace, nginx2Deployment.Name)
			Expect(err).NotTo(HaveOccurred())
		})

		By("Backing up the cluster", func() {
			backupCmd := exec.Command("bbr", append(bbrArgs, "backup")...)
			backupCmd.Dir = bbrDir
			session, err := gexec.Start(backupCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, "1m").Should(gexec.Exit(0))
		})

		By("Deleting workload 2 and installing workload 3", func() {
			err = deploymentApi.Delete(nginx2Deployment.Name, &metav1.DeleteOptions{})
			Expect(err).ToNot(HaveOccurred())

			nginx3Deployment = NewDeployment("nginx-3", GetNginxDeploymentSpec())
			nginx3Deployment, err = k8s.AppsV1().Deployments(namespace).Create(nginx3Deployment)
			Expect(err).ToNot(HaveOccurred())
			err = WaitForDeployment(deploymentApi, namespace, nginx3Deployment.Name)
			Expect(err).NotTo(HaveOccurred())
		})

		By("Restoring the backup", func() {
			globbedFiles, err := filepath.Glob(bbrDir + "/" + deploymentName + "*")
			Expect(err).ToNot(HaveOccurred())
			Expect(globbedFiles).To(HaveLen(1))

			restoreCmd := exec.Command("bbr", append(bbrArgs, "restore", "--artifact-path", globbedFiles[0])...)
			session, err := gexec.Start(restoreCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, "1m").Should(gexec.Exit(0))
		})

		By("Waiting for API to be available", func() {
			Eventually(func() bool {
				var status int
				k8s.CoreV1().RESTClient().Get().RequestURI("/healthz").Do().StatusCode(&status)
				if status == http.StatusOK {
					return true
				}
				return false
			}, "60s", "5s").Should(BeTrue())

		})

		By("Waiting for workloads 1 and 2 to be available", func() {
			err = WaitForDeployment(deploymentApi, namespace, nginx1Deployment.Name)
			Expect(err).NotTo(HaveOccurred())
			err = WaitForDeployment(deploymentApi, namespace, nginx2Deployment.Name)
			Expect(err).NotTo(HaveOccurred())
		})

		By("Asserting that workload 3 is gone", func() {
			_, err = deploymentApi.Get(nginx3Deployment.Name, metav1.GetOptions{})
			Expect(err).To(HaveOccurred())
			statusErr, ok := err.(*errors.StatusError)
			Expect(ok).To(BeTrue())
			Expect(statusErr.ErrStatus.Code).To(Equal(int32(404)))
		})

		By("Waiting for system workloads", func() {
			expectedSelector := []string{"kube-dns", "heapster", "kubernetes-dashboard", "influxdb"}
			runner := NewKubectlRunner()

			systemDeploymentApi := k8s.AppsV1().Deployments("kube-system")
			for _, selector := range expectedSelector {
				deployment := runner.GetResourceNameBySelector("kube-system", "deployment", "k8s-app="+selector)
				err = WaitForDeployment(systemDeploymentApi, "kube-system", deployment)
				Expect(err).NotTo(HaveOccurred())
			}
		})
	})

})
