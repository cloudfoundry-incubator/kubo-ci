package generic_test

import (
	"crypto/tls"
	"fmt"
	"net/http"

	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var _ = Describe("Kubelet", func() {

	var (
		firstWorkerIP string
		err           error
		endpoint      string
	)
	BeforeEach(func() {
		firstWorkerIP, err = GetNodeIP()
		Expect(err).NotTo(HaveOccurred())
		endpoint = fmt.Sprintf("https://%s:10250/pods", firstWorkerIP)
	})

	It("Should fail when unauthenticated requests are made to kubelet", func() {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		invalidRequest(tr, endpoint)
	})

	It("Should respond successful with valid Bearer Token", func() {
		bearerToken, err := BearerToken()
		Expect(err).ToNot(HaveOccurred())

		resp, err := CurlInsecureWithToken(endpoint, bearerToken)

		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(200))
	})

	Context("When using Service Accounts", func() {
		var kubeclient kubernetes.Interface
		var sa *v1.ServiceAccount
		var err error

		BeforeEach(func() {
			kubeclient, err = NewKubeClient()
			Expect(err).NotTo(HaveOccurred())

			sa, err = kubeclient.Core().ServiceAccounts("default").Create(&v1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: "robot-beep-bop"}})
			Expect(err).NotTo(HaveOccurred())

			// Wait for kube-controller-manager to create a token
			Eventually(func() bool {
				sa, _ = kubeclient.Core().ServiceAccounts("default").Get("robot-beep-bop", metav1.GetOptions{})
				return len(sa.Secrets) != 0
			}).Should(BeTrue())
		})

		AfterEach(func() {
			kubeclient.Core().ServiceAccounts("default").Delete("robot-beep-bop", &metav1.DeleteOptions{})
		})

		It("Should reject unauthorized Service Account curl", func() {
			secret, err := kubeclient.Core().Secrets("default").Get(sa.Secrets[0].Name, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())

			resp, err := CurlInsecureWithToken(endpoint, string(secret.Data["token"]))
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(403))
		})
	})

	It("Should fail when requests are made to kubelet with invalid Bearer Token", func() {
		resp, err := CurlInsecureWithToken(endpoint, "IMAFAKEBEAR")
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(401))
	})

	It("Should fail when requests are made to kubelet with invalid cert", func() {
		cert, err := tls.LoadX509KeyPair(PathFromRoot("src/tests/integration-tests/fixtures/selfsigned-client.cert"), PathFromRoot("src/tests/integration-tests/fixtures/selfsigned-client.key"))
		Expect(err).NotTo(HaveOccurred())

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: true},
		}
		invalidRequest(tr, endpoint)
	})
})

func invalidRequest(tr *http.Transport, endpoint string) {
	client := &http.Client{Transport: tr}
	resp, err := client.Get(endpoint)
	Expect(err).ToNot(HaveOccurred())
	Expect(resp.StatusCode).To(Equal(401))
}
