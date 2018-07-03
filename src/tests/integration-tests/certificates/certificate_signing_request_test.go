package certificates_test

import (
	"encoding/base64"
	"io/ioutil"
	"os"

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
		Eventually(kubectl.RunKubectlCommand("delete", "-f", csrSpec), "30s").Should(gexec.Exit(0))
		Eventually(kubectl.RunKubectlCommand("config", "unset", "users."+csrUsername), "30s").Should(gexec.Exit(0))

		if certFile != "" {
			os.Remove(certFile)
		}
	})

	Context("When a user creates a CSR within the 'system:master' group", func() {
		It("should create a client certificate that can talk to Kube API Server", func() {
			Eventually(kubectl.RunKubectlCommand("apply", "-f", csrSpec), "30s").Should(gexec.Exit(0))
			Eventually(kubectl.RunKubectlCommand("certificate", "approve", "test-csr"), "30s").Should(gexec.Exit(0))
			Eventually(kubectl.GetOutput("get", "csr", "test-csr", "-o", "jsonpath={.status}"), "30s").ShouldNot(BeEmpty())

			clientCert := kubectl.GetOutput("get", "csr", "test-csr", "-o", "jsonpath={.status.certificate}")
			decodedCert := decodeCert(clientCert[0])
			certFile = writeCertToFile(decodedCert)

			Eventually(kubectl.RunKubectlCommand("config", "set-credentials", csrUsername,
				"--client-certificate", certFile, "--client-key", keyFile), "30s").Should(gexec.Exit(0))

			Eventually(kubectl.RunKubectlCommand("--user", csrUsername, "get", "nodes"), "30s").Should(gexec.Exit(0))
		})
	})
})

func decodeCert(cert string) []byte {
	decodedCert, err := base64.StdEncoding.DecodeString(cert)
	Expect(err).NotTo(HaveOccurred())

	return decodedCert
}

func writeCertToFile(cert []byte) string {
	tmpFile, err := ioutil.TempFile("/tmp", "client-cert")
	Expect(err).NotTo(HaveOccurred())

	_, err = tmpFile.Write(cert)
	Expect(err).NotTo(HaveOccurred())

	return tmpFile.Name()
}
