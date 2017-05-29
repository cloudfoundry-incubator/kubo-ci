package cleaner_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCleaner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cleaner Suite")
}
