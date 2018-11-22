package certificates_test

import (
	"io/ioutil"
	"os"

	"k8s.io/api/certificates/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "tests/test_helpers"

	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Certificate Signing Requests", func() {

	var (
		csrSpec     string
		keyFile     string
		certFile    string
		csrUsername string
		kubectl     *KubectlRunner
	)

	BeforeEach(func() {
		csrSpec = PathFromRoot("specs/csr/csr.yml")
		keyFile = PathFromRoot("specs/csr/key.pem")
		csrUsername = "test-user-" + GenerateRandomUUID()
		kubectl = NewKubectlRunner()
	})

	AfterEach(func() {
		Eventually(kubectl.StartKubectlCommand("delete", "-f", csrSpec), "30s").Should(gexec.Exit(0))
		Eventually(kubectl.StartKubectlCommand("config", "unset", "users."+csrUsername), "30s").Should(gexec.Exit(0))

		if certFile != "" {
			os.Remove(certFile)
		}
	})

	Context("When a user creates a CSR within the 'system:master' group", func() {
		It("should create a client certificate that can talk to Kube API Server", func() {
			Eventually(kubectl.StartKubectlCommand("apply", "-f", csrSpec), "30s").Should(gexec.Exit(0))

			k8s, err := NewKubeClient()
			Expect(err).NotTo(HaveOccurred())

			CSRClient := k8s.CertificatesV1beta1().CertificateSigningRequests()
			pendingCSR, err := CSRClient.Get("test-csr", v1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())

			pendingCSR.Status.Conditions = append(pendingCSR.Status.Conditions, v1beta1.CertificateSigningRequestCondition{
				Type:    v1beta1.CertificateApproved,
				Reason:  "because I said so",
				Message: "just do it",
			})

			_, err = CSRClient.UpdateApproval(pendingCSR)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() []byte {
				clientCert, err := CSRClient.Get("test-csr", v1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				return clientCert.Status.Certificate
			}, "30s").ShouldNot(BeEmpty())

			clientCert, err := CSRClient.Get("test-csr", v1.GetOptions{})
			certFile := writeCertToFile(clientCert.Status.Certificate)

			Eventually(kubectl.StartKubectlCommand("config", "set-credentials", csrUsername,
				"--client-certificate", certFile, "--client-key", keyFile), "30s").Should(gexec.Exit(0))

			Eventually(kubectl.StartKubectlCommand("--user", csrUsername, "get", "nodes"), "30s").Should(gexec.Exit(0))
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
