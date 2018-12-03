package test_helpers

import (
	"os"

	. "github.com/onsi/gomega"
)

func MustHaveEnv(keyname string) string {
	val := os.Getenv(keyname)
	Expect(val).NotTo(BeEmpty(), "Environment variable '"+keyname+"' must be set")
	return val
}
