package cri_failure_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	iaas string
)

func TestCRIFailure(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CRI Failure Suite")
}
