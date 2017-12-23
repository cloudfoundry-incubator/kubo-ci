package oss_only_test

import (
	"tests/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestOssOnly(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OssOnly Suite")
}

var testconfig *config.Config

var _ = BeforeSuite(func() {
	var err error
	testconfig, err = config.InitConfig()
	Expect(err).NotTo(HaveOccurred())
})
