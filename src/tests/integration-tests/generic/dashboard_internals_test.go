package generic_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var _ = Describe("Dashboard Internals", func() {

	var (
		k8s     kubernetes.Interface
		kubectl *KubectlRunner
		err     error
	)

	type DashboardIntegrationStatus struct {
		Connected bool `json:"connected"`
	}

	type InfluxDBResults struct {
		Results []struct {
			Series []struct {
				Values [][]string `json:"values"`
			} `json:"series"`
		} `json:"results"`
	}

	BeforeEach(func() {
		k8s, err = NewKubeClient()
		Expect(err).ToNot(HaveOccurred())

		kubectl = NewKubectlRunner()
		kubectl.CreateNamespace()
	})

	AfterEach(func() {
		kubectl.RunKubectlCommand("delete", "namespace", kubectl.Namespace())
	})

	It("dashboard should be able to connect to heapster", func() {
		var nodeIP string
		nodeIP, err = GetNodeIP()

		svc, err := k8s.Core().Services("kube-system").Get("kubernetes-dashboard", metav1.GetOptions{})
		Expect(err).ToNot(HaveOccurred())
		nodePort := svc.Spec.Ports[0].NodePort

		endpoint := fmt.Sprintf("https://%s:%d/api/v1/integration/heapster/state", nodeIP, nodePort)

		resp, err := CurlInsecure(endpoint)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(200))

		respBytes, err := ioutil.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())

		heapsterIntegrationStatus := DashboardIntegrationStatus{}
		err = json.Unmarshal(respBytes, &heapsterIntegrationStatus)
		Expect(err).ToNot(HaveOccurred())
		Expect(heapsterIntegrationStatus.Connected).To(BeTrue())
	})

	It("heapster should be able to connect to influxdb", func() {
		url := "https://monitoring-influxdb.kube-system.svc.cluster.local:8086/query"

		session := kubectl.RunKubectlCommand("run", "influxdb-test", "--image=tutum/curl",
			"--restart=Never", "-it", "--rm", "--",
			"curl", "-k", url, "--data-urlencode", "q=SHOW DATABASES")
		Eventually(session, "30s").Should(gexec.Exit(0))

		influxDBStatus := InfluxDBResults{}
		err = json.Unmarshal(session.Out.Contents(), &influxDBStatus)
		Expect(err).ToNot(HaveOccurred())

		defaultHeapsterDB := []string{"k8s"}
		Expect(influxDBStatus.Results[0].Series[0].Values).To(ContainElement(defaultHeapsterDB))
	})
})
