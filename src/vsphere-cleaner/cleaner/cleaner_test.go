package cleaner_test

import (
	"errors"
	"net/url"

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
	var fakeConfig *vspherefakes.FakeConfig

	BeforeEach(func() {
		fakeParser = new(parserfakes.FakeParser)
		fakeVSphereClient = new(vspherefakes.FakeClient)
		fakeConfig = new(vspherefakes.FakeConfig)
		builder := func(*url.URL) (vsphere.Client, error) {
			return fakeVSphereClient, nil
		}
		cleanerObj = cleaner.NewCleaner("lock", fakeParser, builder)
		fakeConfig.DirectorIPReturns("10.2.2.1")
		fakeConfig.UsedIPsReturns([]string{"10.2.3.1", "10.2.3.10"}, nil)
		fakeParser.ParseReturns(fakeConfig, nil)
	})

	It("should parse the lock", func() {
		Expect(cleanerObj.Clean()).To(Succeed())

		Expect(fakeParser.ParseCallCount()).To(Equal(1))
		Expect(fakeParser.ParseArgsForCall(0)).To(Equal("lock"))
	})

	Context("when parsing lock fails", func() {
		It("should fail if parsing the lock fails", func() {
			fakeParser.ParseReturns(new(vspherefakes.FakeConfig), errors.New("I c4n haz eRr0rz"))

			err := cleanerObj.Clean()

			Expect(err).To(HaveOccurred())
		})
	})

	Context("when vsphere client creation fails", func() {
		It("should return error", func() {
			builder := func(*url.URL) (vsphere.Client, error) {
				return fakeVSphereClient, errors.New("error")
			}
			cleanerObj = cleaner.NewCleaner("lock", fakeParser, builder)
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

	It("should destroy the all vms from UsedIPs", func() {
		err := cleanerObj.Clean()

		Expect(err).NotTo(HaveOccurred())
		Expect(fakeVSphereClient.DeleteVMArgsForCall(1)).To(Equal("10.2.3.1"))
		Expect(fakeVSphereClient.DeleteVMArgsForCall(2)).To(Equal("10.2.3.10"))
	})

	It("should return error if deleting the vms in CIDR fails", func() {
		fakeVSphereClient.DeleteVMReturnsOnCall(1, errors.New("Some error"))

		err := cleanerObj.Clean()

		Expect(err).To(HaveOccurred())
	})

	Context("when getting used ips from config fails", func() {
		It("should return error", func() {
			fakeConfig.UsedIPsReturns([]string{}, errors.New("error"))
			err := cleanerObj.Clean()

			Expect(err).To(HaveOccurred())

		})
	})
})
