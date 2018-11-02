package cri_failure_test

import (
	"tests/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	iaas       string
	testconfig *config.Config
)

func TestCRIFailure(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CRI Failure Suite")
}

var _ = BeforeSuite(func() {
	var err error
	testconfig, err = config.InitConfig()
	Expect(err).NotTo(HaveOccurred())
})

func CRIFailureDescribe(description string, callback func()) bool {
	return Describe("[CRI_failure]", func() {
		Describe(description, callback)
	})
}
