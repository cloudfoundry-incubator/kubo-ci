package generic_test

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/cloudfoundry/bosh-cli/director"

	"tests/config"
	. "tests/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kubelet", func() {

	var (
		deployment director.Deployment
		testconfig *config.Config
	)

	BeforeEach(func() {
		var err error
		testconfig, err = config.InitConfig()
		Expect(err).NotTo(HaveOccurred())
		director := NewDirector(testconfig.Bosh)
		deployment, err = director.FindDeployment(testconfig.Bosh.Deployment)
		Expect(err).NotTo(HaveOccurred())
	})

	FIt("Should fail when unauthenticated requests are made to kubelet", func() {
		firstWorkerIP := GetWorkerIP(deployment)
		endpoint := fmt.Sprintf("https://%s:10250/pods", firstWorkerIP)
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Get(endpoint)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(401))
	})
})
