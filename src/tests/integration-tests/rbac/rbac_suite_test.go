package rbac_test

import (
	"tests/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGeneric(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RBAC Suite")
}

var testconfig *config.Config

var _ = BeforeSuite(func() {
	var err error
	testconfig, err = config.InitConfig()
	Expect(err).NotTo(HaveOccurred())
})

func RBACDescribe(description string, callback func()) bool {
	return Describe("[rbac]", func() {
		Describe(description, callback)
	})
}
