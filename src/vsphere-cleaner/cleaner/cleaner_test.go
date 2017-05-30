package cleaner_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"vsphere-cleaner/cleaner"
	"vsphere-cleaner/parser"
	"vsphere-cleaner/parser/parserfakes"
	"vsphere-cleaner/vsphere/vspherefakes"
)

var _ = Describe("Cleaner", func() {
	var fakeVSphereClient *vspherefakes.FakeClient
	var fakeparser *parserfakes.FakeParser

	BeforeEach(func() {
		fakeparser = new(parserfakes.FakeParser)
		fakeVSphereClient = new(vspherefakes.FakeClient)
	})

	It("should parse the lock", func() {
		cleanerObj := cleaner.NewCleaner("lock", fakeparser, fakeVSphereClient)
		Expect(cleanerObj.Clean()).To(Succeed())

		Expect(fakeparser.ParseCallCount()).To(Equal(1))
		Expect(fakeparser.ParseArgsForCall(0)).To(Equal("lock"))
	})

	It("should fail if parsing the lock fails", func() {
		fakeparser := new(parserfakes.FakeParser)
		cleanerObj := cleaner.NewCleaner("lock", fakeparser, fakeVSphereClient)
		fakeparser.ParseReturns(parser.VSphereConfig{}, errors.New("I c4n haz eRr0rz"))

		err := cleanerObj.Clean()
		Expect(err).To(HaveOccurred())
	})

	It("should destroy the bosh vm", func(){
		fakeParser := new(parserfakes.FakeParser)
		envCleaner := cleaner.NewCleaner("lock", fakeParser, fakeVSphereClient)
		fakeParser.ParseReturns(parser.VSphereConfig{InternalIP: "10.2.2.1"}, nil)
		envCleaner.Clean()
		Expect(fakeVSphereClient.DeleteVMCallCount()).To(Equal(1))
		Expect(fakeVSphereClient.DeleteVMArgsForCall(0)).To(Equal("10.2.2.1"))
	})



})
