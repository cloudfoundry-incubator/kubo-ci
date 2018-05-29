package api_extensions_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestApiExtensions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ApiExtensions Suite")
}

var _ = BeforeSuite(func() {
	SetDefaultEventuallyTimeout(60 * time.Second)
})
