package ipcalc_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIpcalc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ipcalc Suite")
}
