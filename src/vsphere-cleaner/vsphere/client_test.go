package vsphere_test

import (
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"vsphere-cleaner/parser"
	"vsphere-cleaner/vsphere"
	"vsphere-cleaner/vsphere/vspherefakes"
)

var _ = Describe("Client", func() {
	It("should build url from config", func() {
		vSphereUrl := vsphere.BuildUrl(parser.VSphereConfig{IP: "host", User: "user", Password: "password"})
		expectedUrl, _ := url.Parse("https://user:password@host/sdk")
		Expect(vSphereUrl).To(Equal(expectedUrl))
	})

	It("should not return error if the vm is not found", func() {
		client := vsphere.NewClientWithFinder(vspherefakes.FailingFinder())

		Expect(client.DeleteVM("some ip")).ToNot(HaveOccurred())
	})
})
