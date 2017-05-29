package vsphere_test

import (
	"net/url"
	"vsphere-cleaner/parser"
	"vsphere-cleaner/vsphere"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	It("should build url from config", func() {
		vsphereUrl := vsphere.BuildUrl(parser.VMWareConfig{IP: "host", User: "user", Password: "password"})
		expectedUrl, _ := url.Parse("https://user:password@host/sdk")
		Expect(vsphereUrl).To(Equal(expectedUrl))
	})
})
