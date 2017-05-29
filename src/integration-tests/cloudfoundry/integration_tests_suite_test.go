package cloudfoundry_test

import (
	"os"
	"strconv"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"integration-tests/test_helpers"
)

func TestIntegrationTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IntegrationTests Suite")
}

var (
	appsDomain       string
	tcpRouterDNSName string
	tcpPort          int64
	runner           *test_helpers.KubectlRunner
	nginxSpec        = test_helpers.PathFromRoot("specs/nginx.yml")
)

var _ = BeforeSuite(func() {
	var err error

	tcpPort, err = strconv.ParseInt(os.Getenv("WORKLOAD_TCP_PORT"), 10, 64)
	if err != nil || tcpPort <= 0 {
		Fail("Correct WORKLOAD_TCP_PORT has to be set")
	}

	appsDomain = os.Getenv("CF_APPS_DOMAIN")
	if appsDomain == "" {
		Fail("CF_APPS_DOMAIN is not set")
	}

	tcpRouterDNSName = os.Getenv("TCP_ROUTER_DNS_NAME")
	if tcpRouterDNSName == "" {
		Fail("TCP_ROUTER_DNS_NAME is not set")
	}

	runner = test_helpers.NewKubectlRunner()
	runner.RunKubectlCommand("create", "namespace", runner.Namespace()).Wait("60s")
})

var _ = AfterSuite(func() {
	if runner != nil {
		runner.RunKubectlCommand("delete", "namespace", runner.Namespace()).Wait("60s")
	}
})
