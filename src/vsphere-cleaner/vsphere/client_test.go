package vsphere_test

import (
	"errors"

	"vsphere-cleaner/vsphere"
	"vsphere-cleaner/vsphere/vspherefakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware/govmomi/object"
)

var _ = Describe("Client", func() {
	var (
		fakeVM    *vspherefakes.FakeVM
		converter = func(object.Reference) (vsphere.VM, error) {
			return fakeVM, nil
		}
		fakeVmFinder *vspherefakes.FakeVmFinder
		client       vsphere.Client
	)
	BeforeEach(func() {
		fakeVM = &vspherefakes.FakeVM{}
		fakeVmFinder = &vspherefakes.FakeVmFinder{}
		fakeVmFinder.FindByIpReturns(&object.VirtualMachine{}, nil)
		client = vsphere.NewClientWithFinder(fakeVmFinder, converter)
	})

	It("destroys VM", func() {
		err := client.DeleteVM("some ip")
		Expect(err).ToNot(HaveOccurred())
		Expect(fakeVM.PowerOffCallCount()).To(Equal(1))
		Expect(fakeVM.DestroyCallCount()).To(Equal(1))
	})

	Context("when VM search fails with error", func() {
		It("should return error", func() {
			fakeVmFinder.FindByIpReturns(nil, errors.New("Some error"))
			err := client.DeleteVM("some ip")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when VM is not found without error", func() {
		It("should not return error", func() {
			fakeVmFinder.FindByIpReturns(nil, nil)
			err := client.DeleteVM("some ip")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when VM can not be converted", func() {
		BeforeEach(func() {
			var converter = func(object.Reference) (vsphere.VM, error) {
				return nil, errors.New("error")
			}
			client = vsphere.NewClientWithFinder(fakeVmFinder, converter)
		})

		It("returns an error", func() {
			err := client.DeleteVM("some ip")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when PowerOff fails", func() {
		BeforeEach(func() {
			fakeVM.PowerOffReturns(errors.New("PowerOff Failed"))
		})
		It("returns an error", func() {
			err := client.DeleteVM("some ip")
			Expect(err).To(HaveOccurred())
		})

		It("does not tries to destroy VM", func() {
			client.DeleteVM("some ip")
			Expect(fakeVM.DestroyCallCount()).To(Equal(0))
		})
	})

	Context("when Destroy fails", func() {
		It("returns an error", func() {
			fakeVM.DestroyReturns(errors.New("Destroy Failed"))
			err := client.DeleteVM("some ip")
			Expect(err).To(HaveOccurred())
		})
	})
})
