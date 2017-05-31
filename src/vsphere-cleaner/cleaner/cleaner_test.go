package cleaner_test

import (
	"errors"

	"vsphere-cleaner/cleaner"
	"vsphere-cleaner/parser/parserfakes"
	"vsphere-cleaner/vsphere"
	"vsphere-cleaner/vsphere/vspherefakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cleaner", func() {
	var fakeVSphereClient *vspherefakes.FakeClient
	var fakeParser *parserfakes.FakeParser
	var cleanerObj cleaner.Cleaner

	BeforeEach(func() {
		fakeParser = new(parserfakes.FakeParser)
		fakeVSphereClient = new(vspherefakes.FakeClient)
		cleanerObj = cleaner.NewCleaner("lock", fakeParser, fakeVSphereClient)
		fakeParser.ParseReturns(vsphere.Config{InternalIP: "10.2.2.1", InternalCIDR: "10.2.3.0/29", ReservedIPs: []string{"10.2.3.2-10.2.3.3", "10.2.3.4"}}, nil)
	})

	It("should parse the lock", func() {
		Expect(cleanerObj.Clean()).To(Succeed())

		Expect(fakeParser.ParseCallCount()).To(Equal(1))
		Expect(fakeParser.ParseArgsForCall(0)).To(Equal("lock"))
	})

	Context("when parsing lock fails", func() {
		It("should fail if parsing the lock fails", func() {
			fakeParser.ParseReturns(vsphere.Config{}, errors.New("I c4n haz eRr0rz"))

			err := cleanerObj.Clean()

			Expect(err).To(HaveOccurred())
		})
	})

	It("should destroy the bosh vm", func() {

		err := cleanerObj.Clean()

		Expect(err).NotTo(HaveOccurred())
		Expect(fakeVSphereClient.DeleteVMArgsForCall(0)).To(Equal("10.2.2.1"))
	})

	It("should return error if Deleting BOSH vm fails", func() {
		fakeVSphereClient.DeleteVMReturnsOnCall(0, errors.New("Some Error"))

		err := cleanerObj.Clean()

		Expect(err).To(HaveOccurred())
	})

	It("should destroy the vms in CIDR", func() {
		err := cleanerObj.Clean()

		Expect(err).NotTo(HaveOccurred())
		Expect(fakeVSphereClient.DeleteVMArgsForCall(1)).To(Equal("10.2.3.1"))
	})

	It("should return error if deleting the vms in CIDR fails", func() {
		fakeVSphereClient.DeleteVMReturnsOnCall(1, errors.New("Some error"))

		err := cleanerObj.Clean()

		Expect(err).To(HaveOccurred())
	})

	It("should not destroy VMs in ReservedIPs", func() {
		err := cleanerObj.Clean()

		Expect(err).ToNot(HaveOccurred())
		Expect(fakeVSphereClient.Invocations()["DeleteVM"]).NotTo(ContainElement([]interface{}{"10.2.3.2"}))
		Expect(fakeVSphereClient.Invocations()["DeleteVM"]).NotTo(ContainElement([]interface{}{"10.2.3.3"}))
		Expect(fakeVSphereClient.Invocations()["DeleteVM"]).NotTo(ContainElement([]interface{}{"10.2.3.4"}))
	})
})
