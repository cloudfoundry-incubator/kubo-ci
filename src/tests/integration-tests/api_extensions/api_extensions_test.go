package api_extensions_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Api Extensions", func() {

	const apiExtensionsNamespace = "wardle"
	const defaultNamespace = "kube-system"
	const sampleApiEndpoint = "v1alpha1.wardle.k8s.io"

	var (
		kubectl                   *KubectlRunner
		session                   *gexec.Session
		serviceAccountSpec        = PathFromRoot("specs/sample-apiserver/sa.yaml")
		authDelegatorSpec         = PathFromRoot("specs/sample-apiserver/auth-delegator.yaml")
		authReaderSpec            = PathFromRoot("specs/sample-apiserver/auth-reader.yaml")
		replicationControllerSpec = PathFromRoot("specs/sample-apiserver/rc.yaml")
		serviceSpec               = PathFromRoot("specs/sample-apiserver/service.yaml")
		apiServiceSpec            = PathFromRoot("specs/sample-apiserver/apiservice.yaml")
	)

	BeforeEach(func() {
		kubectl = NewKubectlRunner()
		session = kubectl.RunKubectlCommand("create", "namespace", apiExtensionsNamespace)
		Eventually(session, "60s").Should(gexec.Exit(0))
	})

	AfterEach(func() {
		kubectl.RunKubectlCommandInNamespace(apiExtensionsNamespace, "delete", "-f", serviceAccountSpec)
		kubectl.RunKubectlCommandInNamespace(defaultNamespace, "delete", "-f", authDelegatorSpec)
		kubectl.RunKubectlCommandInNamespace(defaultNamespace, "delete", "-f", authReaderSpec)
		kubectl.RunKubectlCommandInNamespace(apiExtensionsNamespace, "delete", "-f", replicationControllerSpec)
		kubectl.RunKubectlCommandInNamespace(apiExtensionsNamespace, "delete", "-f", serviceSpec)
		kubectl.RunKubectlCommandInNamespace(apiExtensionsNamespace, "delete", "-f", apiServiceSpec)
		kubectl.RunKubectlCommand("delete", "namespace", apiExtensionsNamespace)
	})

	It("successfully deploys an api service", func() {

		By("creating the associated service account")
		session = kubectl.RunKubectlCommandInNamespace(apiExtensionsNamespace, "create", "-f", serviceAccountSpec)
		Eventually(session, "60s").Should(gexec.Exit(0))

		By("creating the rolebindings for authentication delegation")
		session = kubectl.RunKubectlCommandInNamespace(defaultNamespace, "create", "-f", authDelegatorSpec)
		Eventually(session, "60s").Should(gexec.Exit(0))
		session = kubectl.RunKubectlCommandInNamespace(defaultNamespace, "create", "-f", authReaderSpec)
		Eventually(session, "60s").Should(gexec.Exit(0))

		By("creating the service and replication and replication controller")
		session = kubectl.RunKubectlCommandInNamespace(apiExtensionsNamespace, "create", "-f", replicationControllerSpec)
		Eventually(session, "60s").Should(gexec.Exit(0))
		session = kubectl.RunKubectlCommandInNamespace(apiExtensionsNamespace, "create", "-f", serviceSpec)
		Eventually(session, "60s").Should(gexec.Exit(0))

		By("creating the api service extension")
		session = kubectl.RunKubectlCommandInNamespace(apiExtensionsNamespace, "create", "-f", apiServiceSpec)
		Eventually(session, "60s").Should(gexec.Exit(0))

		session = kubectl.RunKubectlCommandInNamespace(apiExtensionsNamespace, "get", "pods")
		Eventually(session, "120s").Should(gexec.Exit(0))
		re := regexp.MustCompile(`(wardle-server-\w+)`)
		matches := re.FindStringSubmatch(string(session.Out.Contents()))
		podName := matches[1]
		Eventually(func() string {
			podStatus := kubectl.GetPodStatus(apiExtensionsNamespace, podName)
			return podStatus
		}, "120s", "2s").Should(Equal("Running"))

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
		jsonResp, err := ioutil.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())
		json.Unmarshal(jsonResp, &apiServiceResp)
		Expect(apiServiceResp.Metadata.Name).To(Equal(sampleApiEndpoint))
	})
})
