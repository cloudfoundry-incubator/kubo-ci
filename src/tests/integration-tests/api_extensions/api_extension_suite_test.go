package api_extensions_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestApiExtensions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ApiExtensions Suite")
}
