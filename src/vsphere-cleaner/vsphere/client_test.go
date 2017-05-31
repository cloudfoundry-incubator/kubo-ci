package vsphere_test

import (
	"errors"
	"net/url"

	"vsphere-cleaner/vsphere"
	"vsphere-cleaner/vsphere/vspherefakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	It("should build url from config", func() {
		vSphereUrl := vsphere.BuildUrl(vsphere.Config{IP: "host", User: "user", Password: "password"})
		expectedUrl, _ := url.Parse("https://user:password@host/sdk")
		Expect(vSphereUrl).To(Equal(expectedUrl))
	})

	It("should not return error if the vm is not found", func() {
		fakeVmFinder := &vspherefakes.FakeVmFinder{}
		client := vsphere.NewClientWithFinder(fakeVmFinder)
		fakeVmFinder.FindByIpReturns(nil, errors.New("Some error"))
		err := client.DeleteVM("some ip")
		Expect(err).ToNot(HaveOccurred())
	})
})
