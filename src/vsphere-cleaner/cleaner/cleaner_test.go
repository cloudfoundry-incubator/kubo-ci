package cleaner_test

import (
	"errors"
	"vsphere-cleaner/cleaner"
	"vsphere-cleaner/parser"
	"vsphere-cleaner/parser/parserfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cleaner", func() {
	It("should parse the lock", func() {
		fakeparser := new(parserfakes.FakeParser)
		cleanerObj := cleaner.NewCleaner("lock", fakeparser)
		Expect(cleanerObj.Clean()).To(Succeed())

		Expect(fakeparser.ParseCallCount()).To(Equal(1))
		Expect(fakeparser.ParseArgsForCall(0)).To(Equal("lock"))
	})

	It("should fail if parsing the lock fails", func() {
		fakeparser := new(parserfakes.FakeParser)
		cleanerObj := cleaner.NewCleaner("lock", fakeparser)
		fakeparser.ParseReturns(parser.VMWareConfig{}, errors.New("I am error"))

		err := cleanerObj.Clean()
		Expect(err).To(HaveOccurred())
	})

	// It("should calculate IPs to be terminated", func() {
	// 	fakeparser := new(parserfakes.FakeParser)
	// 	cleanerObj := cleaner.NewCleaner("lock", fakeparser, fakeVSphereClient)
	// 	fakeparser.ParseReturns(parser.VMWareConfig{InternalCIDR: "10.2.2.0/30", ReservedIPs: []string{"10.2.2.1"}}, nil)

	// 	cleanerObj.Clean()

	// 	Expect(fakeVSphereClient.DeleteVMCallCount()).To(Equal(3))
	// })
})
