package cidrs_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "tests/test_helpers"
)

func TestCidrs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cidrs Suite")
}

var _ = BeforeSuite(func() {
	CheckRequiredEnvs([]string{"CIDR_VARS_FILE"})
})
