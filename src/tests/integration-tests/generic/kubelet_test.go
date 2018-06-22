package generic_test

import (
	"crypto/tls"
	"fmt"
	"net/http"

	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		req, err := http.NewRequest("GET", endpoint, nil)
		Expect(err).ToNot(HaveOccurred())

		req.Header.Add("Authorization", "Bearer "+bearerToken)
		resp, err := client.Do(req)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(200))
	})

	It("Should fail when requests are made to kubelet with invalid Bearer Token", func() {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		req, err := http.NewRequest("GET", endpoint, nil)
		Expect(err).ToNot(HaveOccurred())

		req.Header.Add("Authorization", "Bearer IMFAKEBEAR")
		resp, err := client.Do(req)
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
