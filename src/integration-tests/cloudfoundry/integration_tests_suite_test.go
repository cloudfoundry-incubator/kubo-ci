package cloudfoundry_test

import (
	"os"
	"strconv"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIntegrationTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IntegrationTests Suite")
}

var (
	appsDomain       string
	tcpRouterDNSName string
	tcpPort          int64
)

var _ = BeforeSuite(func() {
	tcpPort, err := strconv.ParseInt(os.Getenv("WORKLOAD_TCP_PORT"), 10, 64)
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
})
