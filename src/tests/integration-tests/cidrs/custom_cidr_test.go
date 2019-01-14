package cidrs_test

import (
	"io/ioutil"
	"net"
	. "tests/test_helpers"

	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	client_v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type CIDRConfig struct {
	ClusterIPRange      string `yaml:"service_cluster_cidr"`
	KubeDNSIP           string `yaml:"kubedns_service_ip"`
	KubernetesServiceIP string `yaml:"first_ip_of_service_cluster_cidr"`
	PodIPRange          string `yaml:"pod_network_cidr"`
}

var _ = Describe("Custom CIDRs", func() {
	var (
		k8s               kubernetes.Interface
		testNamespaceName string
		svcController     client_v1.ServiceInterface
		cidrConfig        CIDRConfig
	)

	BeforeEach(func() {
		var err error
		k8s, err = NewKubeClient()
		Expect(err).NotTo(HaveOccurred())

		testNamespace, err := CreateTestNamespace(k8s, "test")
		Expect(err).NotTo(HaveOccurred())
		testNamespaceName = testNamespace.Name

		svcController = k8s.CoreV1().Services(testNamespaceName)

		cidrVarsFile := os.Getenv("CIDR_VARS_FILE")
		b, err := ioutil.ReadFile(PathFromRoot(cidrVarsFile))
		Expect(err).ToNot(HaveOccurred())
		err = yaml.Unmarshal(b, &cidrConfig)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		DeleteTestNamespace(k8s, testNamespaceName)
	})

	Context("Services", func() {
		It("configures Kubernetes API server to the provided IP", func() {
			service, err := k8s.CoreV1().Services("default").Get("kubernetes", meta_v1.GetOptions{})

			Expect(err).NotTo(HaveOccurred())
			Expect(service.Spec.ClusterIP).To(Equal(cidrConfig.KubernetesServiceIP))
		})

		It("creates service in the specified CIDR", func() {
			svcName := "test-svc-cidr-" + GenerateRandomUUID()
			svcSpec := v1.Service{
				ObjectMeta: meta_v1.ObjectMeta{Name: svcName},
				Spec:       v1.ServiceSpec{Ports: []v1.ServicePort{{Protocol: v1.ProtocolTCP, Port: 80}}},
			}
			svc, err := svcController.Create(&svcSpec)
			defer svcController.Delete(svcName, &meta_v1.DeleteOptions{})

			Expect(err).NotTo(HaveOccurred())
			_, subnet, _ := net.ParseCIDR(cidrConfig.ClusterIPRange)
			Expect(subnet.Contains(net.ParseIP(svc.Spec.ClusterIP))).To(BeTrue())
		})

		It("configures DNS to the provided IP", func() {
			service, err := k8s.CoreV1().Services("kube-system").Get("kube-dns", meta_v1.GetOptions{})

			Expect(err).NotTo(HaveOccurred())
			Expect(service.Spec.ClusterIP).To(Equal(cidrConfig.KubeDNSIP))
		})
	})

	Context("Pods", func() {
		It("creates pod in the specified CIDR", func() {
			podName := "test-pod-cidr-" + GenerateRandomUUID()
			podSpec := v1.Pod{
				ObjectMeta: meta_v1.ObjectMeta{Name: podName},
				Spec: v1.PodSpec{
					Containers: []v1.Container{{
						Name:  "nginx",
						Image: "nginx",
						Ports: []v1.ContainerPort{{ContainerPort: 80}},
					}},
				},
			}

			pod, err := k8s.CoreV1().Pods(testNamespaceName).Create(&podSpec)

			defer k8s.CoreV1().Pods(testNamespaceName).Delete(podName, &meta_v1.DeleteOptions{})
			Expect(err).NotTo(HaveOccurred())

			_, subnet, _ := net.ParseCIDR(cidrConfig.PodIPRange)
			Eventually(func() bool {
				pod, err = k8s.CoreV1().Pods(testNamespaceName).Get(podName, meta_v1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())

				return subnet.Contains(net.ParseIP(pod.Status.PodIP))
			}, "1m", "5s").Should(BeTrue())
		})
	})
})
