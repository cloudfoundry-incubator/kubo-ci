package cleaner_test

import (
	"vsphere-cleaner/cleaner"
	"vsphere-cleaner/parser/parserfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cleaner", func() {
	It("should parse the lock", func() {
		fakeparser := new(parserfakes.FakeParser)
		cleanerObj := cleaner.NewCleaner("lock", fakeparser)
		cleanerObj.Clean()

		Expect(fakeparser.ParseCallCount()).To(Equal(1))
		Expect(fakeparser.ParseArgsForCall(0)).To(Equal("lock"))
	})
})
