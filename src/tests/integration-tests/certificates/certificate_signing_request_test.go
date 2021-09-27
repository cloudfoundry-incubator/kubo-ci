package certificates_test

import (
	"context"
	"io/ioutil"
	"os"

	certificates "k8s.io/api/certificates/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "tests/test_helpers"

	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Certificate Signing Requests", func() {

	var (
		csrWorkerClientSpec string
		keyWorkerClientFile string
		csrWorkerServerSpec string
		csrMasterSpec string
		keyMasterFile string
		csrUsername   string
		certFile      string
		specsFile     string
		kubectl       *KubectlRunner
	)

	BeforeEach(func() {
		specsFile = PathFromRoot("specs/csr/")
		csrWorkerClientSpec = PathFromRoot("specs/csr/csr-client-worker.yml")
		keyWorkerClientFile = PathFromRoot("specs/csr/key-client-worker.pem")
		csrWorkerServerSpec = PathFromRoot("specs/csr/csr-server-worker.yml")
		csrMasterSpec = PathFromRoot("specs/csr/csr-client-master.yml")
		keyMasterFile = PathFromRoot("specs/csr/key-client-master.pem")
		csrUsername = "test-user-" + GenerateRandomUUID()
		kubectl = NewKubectlRunner()
	})

	AfterEach(func() {
		Eventually(kubectl.StartKubectlCommand("delete", "-f", specsFile, "--ignore-not-found=true"), "30s").Should(gexec.Exit(0))

		Eventually(kubectl.StartKubectlCommand("config", "unset", "users."+csrUsername), "30s").Should(gexec.Exit(0))

		if certFile != "" {
			os.Remove(certFile)
		}
	})

	Context("When a user creates a CSR within the 'system:master' group", func() {
		// About this test case, we need to disable CertificateSubjectRestriction admission plugin in kube-apiserver side first
		// Becase CertificateSubjectRestriction is enabled by default, and it restrict signing cert in system:masters group
		It("should create a client certificate that can talk to Kube API Server", func() {
			Skip("CertificateSubjectRestriction should be disabled first, we don;t support it now")
			Eventually(kubectl.StartKubectlCommand("apply", "-f", csrMasterSpec), "30s").Should(gexec.Exit(0))

			k8s, err := NewKubeClient()
			Expect(err).NotTo(HaveOccurred())

			CSRClient := k8s.CertificatesV1().CertificateSigningRequests()
			pendingCSR, err := CSRClient.Get(context.TODO(), "test-csr-master", v1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())

			pendingCSR.Status.Conditions = append(pendingCSR.Status.Conditions, certificates.CertificateSigningRequestCondition{
				Type:    certificates.CertificateApproved,
				Reason:  "because I said so",
				Message: "just do it",
				Status:  "True",
			})

			_, err = CSRClient.UpdateApproval(context.TODO(), "test-csr-master", pendingCSR, v1.UpdateOptions{})
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() []byte {
				clientCert, err := CSRClient.Get(context.TODO(), "test-csr-master", v1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				return clientCert.Status.Certificate
			}, "30s").ShouldNot(BeEmpty())

			clientCert, err := CSRClient.Get(context.TODO(), "test-csr-master", v1.GetOptions{})
			certFile := writeCertToFile(clientCert.Status.Certificate)

			Eventually(kubectl.StartKubectlCommand("config", "set-credentials", csrUsername,
				"--client-certificate", certFile, "--client-key", keyMasterFile), "30s").Should(gexec.Exit(0))

			Eventually(kubectl.StartKubectlCommand("--user", csrUsername, "get", "nodes"), "30s").Should(gexec.Exit(0))
		})
	})

	Context("When a user creates a CSR within the 'system:nodes' group", func() {
		It("should create a client certificate that can talk to Kube API Server", func() {
			Eventually(kubectl.StartKubectlCommand("apply", "-f", csrWorkerClientSpec), "30s").Should(gexec.Exit(0))

			k8s, err := NewKubeClient()
			Expect(err).NotTo(HaveOccurred())

			CSRClient := k8s.CertificatesV1().CertificateSigningRequests()
			_, err = CSRClient.Get(context.TODO(), "test-csr-client-worker", v1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() []byte {
				clientCert, err := CSRClient.Get(context.TODO(), "test-csr-client-worker", v1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				return clientCert.Status.Certificate
			}, "30s").ShouldNot(BeEmpty())

			clientCert, err := CSRClient.Get(context.TODO(), "test-csr-client-worker", v1.GetOptions{})
			certFile := writeCertToFile(clientCert.Status.Certificate)

			Eventually(kubectl.StartKubectlCommand("config", "set-credentials", csrUsername,
				"--client-certificate", certFile, "--client-key", keyWorkerClientFile), "30s").Should(gexec.Exit(0))

			Eventually(kubectl.StartKubectlCommand("--user", csrUsername, "get", "nodes"), "30s").Should(gexec.Exit(0))
		})
	})

	Context("When a user creates a CSR within the 'system:nodes' group", func() {
		It("should create a serving certificate that can talk to Kube API Server", func() {
			Eventually(kubectl.StartKubectlCommand("apply", "-f", csrWorkerServerSpec), "30s").Should(gexec.Exit(0))

			k8s, err := NewKubeClient()
			Expect(err).NotTo(HaveOccurred())

			CSRClient := k8s.CertificatesV1().CertificateSigningRequests()
			pendingServerCSR, err := CSRClient.Get(context.TODO(), "test-csr-server-worker", v1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())

			pendingServerCSR.Status.Conditions = append(pendingServerCSR.Status.Conditions, certificates.CertificateSigningRequestCondition{
				Type:    certificates.CertificateApproved,
				Reason:  "because I said so",
				Message: "just do it",
				Status:  "True",
			})

			_, err = CSRClient.UpdateApproval(context.TODO(), "test-csr-server-worker", pendingServerCSR, v1.UpdateOptions{})
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() []byte {
				servingCert, err := CSRClient.Get(context.TODO(), "test-csr-server-worker", v1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				return servingCert.Status.Certificate
			}, "30s").ShouldNot(BeEmpty())

		})
	})
})

func writeCertToFile(cert []byte) string {
	tmpFile, err := ioutil.TempFile("/tmp", "client-cert")
	Expect(err).NotTo(HaveOccurred())

	_, err = tmpFile.Write(cert)
	Expect(err).NotTo(HaveOccurred())

	return tmpFile.Name()
}
