package rbac_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGeneric(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RBAC Suite")
}

func RBACDescribe(description string, callback func()) bool {
	return Describe("[rbac]", func() {
		Describe(description, callback)
	})
}
