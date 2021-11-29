package certificates_test

import (
	"context"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	certificates "k8s.io/api/certificates/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "tests/test_helpers"

	"github.com/davecgh/go-spew/spew"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Certificate Signing Requests", func() {
	var (
		csrWorkerClientSpec string
		keyWorkerClientFile string
		csrWorkerServerSpec string
		csrMasterSpec       string
		keyMasterFile       string
		csrUsername         string
		certFile            string
		specsFile           string
		kubectl             *KubectlRunner
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
			// Disable CertificateSubjectRestriction admission plugin
			clusterName := getClusterName()
			profileName := "test-profile-enable"
			admissionPluginArgument := map[string]string{"disable-admission-plugins": "CertificateSubjectRestriction"}
			err := updateClusterWithProfile(clusterName, profileName, admissionPluginArgument)
			Expect(err).NotTo(HaveOccurred())

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

			// Enable CertificateSubjectRestriction admission plugin
			profileName = "test-profile-disable"
			newAdmissionPluginArgument := map[string]string{"disable-admission-plugins": ""}
			err = updateClusterWithProfile(clusterName, profileName, newAdmissionPluginArgument)
			Expect(err).NotTo(HaveOccurred())
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

func updateClusterWithProfile(clusterName string, profileName string, admissionPluginArgument map[string]string) error {
	type Customizations struct {
		Component string
		Arguments map[string]string
	}

	type kubernetesProfile struct {
		Name                        string
		Description                 string
		Experimental_customizations []Customizations
	}

	k8sProfileData := kubernetesProfile{
		Name:        profileName,
		Description: "A kubeneter profile for CertificateSubjectRestriction plugin",
		Experimental_customizations: []Customizations{
			Customizations{
				Component: "kube-apiserver",
				Arguments: admissionPluginArgument,
			},
		},
	}
	k8sProfile, _ := json.MarshalIndent(k8sProfileData, "", " ")
	_ = ioutil.WriteFile("k8sProfileTmp.json", k8sProfile, 0644)
	// Update cluster with profle
	pks_cli, err := SetupPksCli()
	Expect(err).ShouldNot(HaveOccurred())
	path, _ := os.Getwd()
	profilePath := path + "/k8sProfileTmp.json"
	spew.Dump("ProfilePath", profilePath)
	_, err = pks_cli.CreateKubernetesProfile(profilePath)
	_, err = pks_cli.UpdateClusterWithProfile(clusterName, profileName)
	Expect(err).NotTo(HaveOccurred())
	return err
}

func getClusterName() string {
	type EnvInfo struct {
		Name string `yaml:"name"`
		Uuid string `yaml:"uuid"`
		Api  string `yaml:"api"`
	}
	currentPath, _ := os.Getwd()
	rootPath := filepath.Join(currentPath, "../../../../..")
	envInfoPathDir := rootPath + "/gcs-cluster-info/info.yml"
	spew.Dump("Environment info path", envInfoPathDir)
	yamlFile, err := ioutil.ReadFile(envInfoPathDir)
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}
	c := &EnvInfo{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	spew.Dump("Environment name", c.Name)
	return c.Name
}
