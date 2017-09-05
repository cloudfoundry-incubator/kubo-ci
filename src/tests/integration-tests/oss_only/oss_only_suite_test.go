package oss_only_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestOssOnly(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OssOnly Suite")
}
