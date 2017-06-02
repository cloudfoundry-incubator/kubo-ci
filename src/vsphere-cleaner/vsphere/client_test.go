package vsphere_test

import (
	"errors"

	"vsphere-cleaner/vsphere"
	"vsphere-cleaner/vsphere/vspherefakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	It("should not return error if the vm is not found", func() {
		fakeVmFinder := &vspherefakes.FakeVmFinder{}
		client := vsphere.NewClientWithFinder(fakeVmFinder)
		fakeVmFinder.FindByIpReturns(nil, errors.New("Some error"))
		err := client.DeleteVM("some ip")
		Expect(err).ToNot(HaveOccurred())
	})

	It("should not return error if the vm is not found without error", func() {
		fakeVmFinder := &vspherefakes.FakeVmFinder{}
		client := vsphere.NewClientWithFinder(fakeVmFinder)
		fakeVmFinder.FindByIpReturns(nil, nil)
		err := client.DeleteVM("some ip")
		Expect(err).ToNot(HaveOccurred())
	})
})
