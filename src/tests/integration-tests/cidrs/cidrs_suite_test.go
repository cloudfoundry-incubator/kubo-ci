package cidrs_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCidrs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cidrs Suite")
}
