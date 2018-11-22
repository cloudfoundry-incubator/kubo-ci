package api_extensions_test

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Api Extensions", func() {
	const (
		systemNamespace                   = "kube-system"
		serviceAccountSpecTemplate        = "fixtures/sa.yml"
		authDelegatorSpecTemplate         = "fixtures/auth-delegator.yml"
		authReaderSpecTemplate            = "fixtures/auth-reader.yml"
		replicationControllerSpecTemplate = "fixtures/rc.yml"
		serviceSpecTemplate               = "fixtures/service.yml"
		apiServiceSpecTemplate            = "fixtures/apiservice.yml"
	)

	var (
		kubectl *KubectlRunner

		tmpDir                    string
		apiExtensionsNamespace    string
		sampleApiEndpoint         string
		serviceAccountSpec        string
		authDelegatorSpec         string
		authReaderSpec            string
		replicationControllerSpec string
		serviceSpec               string
		apiServiceSpec            string
	)

	templateNamespaceIntoFile := func(tmpDir, path, namespace string) string {
		t, err := template.ParseFiles(path)
		Expect(err).NotTo(HaveOccurred())

		f, err := ioutil.TempFile(tmpDir, filepath.Base(path))
		Expect(err).NotTo(HaveOccurred())
		defer f.Close()

		type templateInfo struct{ Namespace string }
		Expect(t.Execute(f, templateInfo{Namespace: namespace})).To(Succeed())

		return f.Name()
	}

	BeforeEach(func() {
		var err error
		kubectl = NewKubectlRunner()
		kubectl.Setup()
		apiExtensionsNamespace = kubectl.Namespace()

		tmpDir, err = ioutil.TempDir("", "api-extensions")
		Expect(err).NotTo(HaveOccurred())

		serviceAccountSpec = templateNamespaceIntoFile(tmpDir, serviceAccountSpecTemplate, apiExtensionsNamespace)
		authDelegatorSpec = templateNamespaceIntoFile(tmpDir, authDelegatorSpecTemplate, apiExtensionsNamespace)
		authReaderSpec = templateNamespaceIntoFile(tmpDir, authReaderSpecTemplate, apiExtensionsNamespace)
		replicationControllerSpec = templateNamespaceIntoFile(tmpDir, replicationControllerSpecTemplate, apiExtensionsNamespace)
		serviceSpec = templateNamespaceIntoFile(tmpDir, serviceSpecTemplate, apiExtensionsNamespace)
		apiServiceSpec = templateNamespaceIntoFile(tmpDir, apiServiceSpecTemplate, apiExtensionsNamespace)

		sampleApiEndpoint = "v1alpha1." + apiExtensionsNamespace + ".k8s.io"
	})

	AfterEach(func() {
		kubectl.RunKubectlCommandWithTimeout("delete", "-f", apiServiceSpec)
		kubectl.RunKubectlCommandWithTimeout("delete", "-f", serviceAccountSpec)
		session := kubectl.RunKubectlCommandInNamespace(systemNamespace, "delete", "-f", authDelegatorSpec)
		session.Wait("5s")
		fmt.Fprintf(GinkgoWriter, "AuthDelegatorSpec delete exit code %d\n", session.ExitCode())
		session = kubectl.RunKubectlCommandInNamespace(systemNamespace, "delete", "-f", authReaderSpec)
		session.Wait("5s")
		fmt.Fprintf(GinkgoWriter, "AuthReaderSpec delete exit code %d\n", session.ExitCode())
		kubectl.RunKubectlCommandWithTimeout("delete", "-f", replicationControllerSpec)
		kubectl.RunKubectlCommandWithTimeout("delete", "-f", serviceSpec)
		kubectl.Teardown()
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
	})

	It("successfully deploys an api service", func() {
		By("creating the associated service account")
		kubectl.RunKubectlCommandWithTimeout("create", "-f", serviceAccountSpec)

		By("creating the rolebindings for authentication delegation")
		session := kubectl.RunKubectlCommandInNamespace(systemNamespace, "create", "-f", authDelegatorSpec)
		Eventually(session).Should(gexec.Exit(0))
		session = kubectl.RunKubectlCommandInNamespace(systemNamespace, "create", "-f", authReaderSpec)
		Eventually(session).Should(gexec.Exit(0))

		By("creating the service and replication and replication controller")
		kubectl.RunKubectlCommandWithTimeout("create", "-f", replicationControllerSpec)
		kubectl.RunKubectlCommandWithTimeout("create", "-f", serviceSpec)

		By("creating the api service extension")
		kubectl.RunKubectlCommandWithTimeout("create", "-f", apiServiceSpec)

		WaitForPodsToRun(kubectl, kubectl.TimeoutInSeconds*2)

		By("verifying the api extension has been registered to the cluster")
		var apiServiceResp struct {
			Metadata struct {
				Name string `json:"name"`
			} `json:"metadata"`
		}

		session = kubectl.RunKubectlCommand("proxy", "--port=8000")
		Eventually(session).Should(gbytes.Say("Starting to serve on"))
		defer session.Kill()
		resp, err := http.Get(fmt.Sprintf("http://localhost:8000/apis/apiregistration.k8s.io/v1beta1/apiservices/%s", sampleApiEndpoint))
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()
		Expect(resp.StatusCode).To(Equal(200))
		jsonResp, err := ioutil.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())
		err = json.Unmarshal(jsonResp, &apiServiceResp)
		Expect(err).NotTo(HaveOccurred())
		Expect(apiServiceResp.Metadata.Name).To(Equal(sampleApiEndpoint))
	})
})
