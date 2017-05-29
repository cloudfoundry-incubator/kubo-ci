package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestVsphereCleaner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VsphereCleaner Suite")
}

var pathToCLI string

var _ = BeforeSuite(func() {
	var err error
	pathToCLI, err = gexec.Build("vsphere-cleaner")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
