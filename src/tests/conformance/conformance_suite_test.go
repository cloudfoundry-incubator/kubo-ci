package conformance_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestConformance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Conformance suite")
}
