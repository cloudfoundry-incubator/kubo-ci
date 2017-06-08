package persistent_volume_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPersistentVolume(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PersistentVolume Suite")
}
