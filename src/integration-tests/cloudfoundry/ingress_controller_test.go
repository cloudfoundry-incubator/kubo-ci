package cloudfoundry_test

import (
  "os"
  "fmt"
	"net/http"
  "strconv"
  "io/ioutil"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"integration-tests/test_helpers"
)

var _ = Describe("Testing Ingress Controller", func ()  {

  var (
    tlsKubernetesCert string
    tlsKubernetesPrivateKey string
    kubernetesServiceHost string
    kubernetesServicePort int64

    ingressSpec = test_helpers.PathFromRoot("specs/ingress.yml")
    runner           *test_helpers.KubectlRunner
  )

  BeforeEach(func() {
    tlsKubernetesCert = os.Getenv("TLS_KUBERNETES_CERT")
    if tlsKubernetesCert == "" {
      Fail("Correct TLS_KUBERNETES_CERT has to be set")
    }
    tlsKubernetesPrivateKey = os.Getenv("TLS_KUBERNETES_PRIVATE_KEY")
    if tlsKubernetesPrivateKey == "" {
      Fail("Correct TLS_KUBERNETES_PRIVATE_KEY has to be set")
    }
    kubernetesServiceHost = os.Getenv("KUBERNETES_SERVICE_HOST")
    if kubernetesServiceHost == "" {
      Fail("Correct KUBERNETES_SERVICE_HOST has to be set")
    }

    var portErr error
    kubernetesServicePort, portErr = strconv.ParseInt(os.Getenv("KUBERNETES_SERVICE_PORT"), 10, 64)
  	if portErr != nil || kubernetesServicePort <= 0 {
  		Fail("Correct KUBERNETES_SERVICE_PORT has to be set")
  	}

    runner = test_helpers.NewKubectlRunner()
		runner.RunKubectlCommand("create", "namespace", runner.Namespace()).Wait("60s")

    certFile, _ := ioutil.TempFile(os.TempDir(), "cert")
    certFile.WriteString(tlsKubernetesCert)
    defer os.Remove(certFile.Name())

    keyFile, _ := ioutil.TempFile(os.TempDir(), "key")
    keyFile.WriteString(tlsKubernetesPrivateKey)
    defer os.Remove(keyFile.Name())

    Eventually(runner.RunKubectlCommand(
      "create", "secret", "tls", "tls-kubernetes", "--cert", certFile.Name(), "--key", keyFile.Name())).Should(gexec.Exit(0))

    Eventually(runner.RunKubectlCommand(
      "create", "secret", "generic", "kubernetes-service", fmt.Sprintf("--from-literal=host=%s", kubernetesServiceHost), fmt.Sprintf("--from-literal=port=%d", kubernetesServicePort))).Should(gexec.Exit(0))

    Eventually(runner.RunKubectlCommand(
      "create", "-f", ingressSpec), "60s").Should(gexec.Exit(0))
  })

  AfterEach(func() {
    Eventually(runner.RunKubectlCommand(
      "delete", "-f", ingressSpec), "60s").Should(gexec.Exit())

    Eventually(runner.RunKubectlCommand(
      "delete", "secret", "tls-kubernetes")).Should(gexec.Exit())

    Eventually(runner.RunKubectlCommand(
      "delete", "secret", "kubernetes-service")).Should(gexec.Exit())

    runner.RunKubectlCommand("delete", "namespace", runner.Namespace()).Wait("60s")
	})

  It("Allows routing via Ingress Controller", func() {
    serviceName := runner.Namespace()
    appUrl := fmt.Sprintf("http://%s.%s", serviceName, appsDomain)

    By("exposing it via HTTP")
    result, err := http.Get(appUrl)
    Expect(err).NotTo(HaveOccurred())
    Expect(result.StatusCode).NotTo(Equal(200))

    httpLabel := fmt.Sprintf("http-route-sync=%s", serviceName)
    Eventually(runner.RunKubectlCommand("label", "services", "nginx-ingress-controller", httpLabel), "10s").Should(gexec.Exit(0))

    timeout := time.Duration(5 * time.Second)
    httpClient := http.Client{
      Timeout: timeout,
    }
    Eventually(func() int {
      result, err := httpClient.Get(appUrl+ "/simple-http-server")
      if err != nil {
        return -1
      }
      return result.StatusCode
    }, "120s", "5s").Should(Equal(200))

    result, err = http.Get(appUrl)
    Expect(err).NotTo(HaveOccurred())
    Expect(result.StatusCode).To(Equal(404))

    By("exposing it via TCP")
    appUrl = fmt.Sprintf("http://%s:%d", tcpRouterDNSName, tcpPort)

    result, err = http.Get(appUrl)
    Expect(err).To(HaveOccurred())

    tcpLabel := fmt.Sprintf("tcp-route-sync=%d", tcpPort)
    Eventually(runner.RunKubectlCommand("label", "services", "nginx-ingress-controller", tcpLabel), "10s").Should(gexec.Exit(0))
    Eventually(func() error {
      _, err := http.Get(appUrl)
      return err
    }, "120s", "5s").ShouldNot(HaveOccurred())

  })
})
