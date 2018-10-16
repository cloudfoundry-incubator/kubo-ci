package flag_generator_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFlagGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FlagGenerator Suite")
}
