package cri_failure_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	. "tests/test_helpers"
)

func TestCRIFailure(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CRI Failure Suite")
}

var _ = BeforeSuite(func() {
	CheckRequiredEnvs([]string{
		"BOSH_DEPLOYMENT",
		"BOSH_ENVIRONMENT",
		"BOSH_CLIENT",
		"BOSH_CLIENT_SECRET",
		"BOSH_CA_CERT",
		"TURBULENCE_HOST",
		"TURBULENCE_PORT",
		"TURBULENCE_USERNAME",
		"TURBULENCE_PASSWORD",
		"TURBULENCE_CA_CERT",
	})
})
