package cloudfoundry_test

import (
	"os"
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
)

var _ = BeforeSuite(func() {
	appsDomain = os.Getenv("CF_APPS_DOMAIN")
	if appsDomain == "" {
		Fail("CF_APPS_DOMAIN is not set")
	}
	tcpRouterDNSName = os.Getenv("TCP_ROUTER_DNS_NAME")
	if tcpRouterDNSName == "" {
		Fail("TCP_ROUTER_DNS_NAME is not set")
	}
})
