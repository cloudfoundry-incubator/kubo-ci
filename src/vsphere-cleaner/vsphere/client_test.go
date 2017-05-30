package vsphere_test

import (
	"errors"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"vsphere-cleaner/parser"
	"vsphere-cleaner/vsphere"
	"vsphere-cleaner/vsphere/vspherefakes"
)

var _ = Describe("Client", func() {
	It("should build url from config", func() {
		vsphereUrl := vsphere.BuildUrl(parser.VSphereConfig{IP: "host", User: "user", Password: "password"})
		expectedUrl, _ := url.Parse("https://user:password@host/sdk")
		Expect(vsphereUrl).To(Equal(expectedUrl))
	})

	It("should not return error if the vm is not found", func() {
		fakeVMFinder := vspherefakes.FakeVmFinder{Err: errors.New("Sunt lanistaes convertam domesticus, fidelis adgiumes.")}
		client := vsphere.NewClientWithFinder(fakeVMFinder)

		Expect(client.DeleteVM("some ip")).ToNot(HaveOccurred())
	})
})
