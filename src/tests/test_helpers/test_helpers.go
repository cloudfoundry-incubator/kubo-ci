package test_helpers

import (
	"os"

	"fmt"
	. "github.com/onsi/ginkgo"
)

func CheckRequiredEnvs(envs []string) {
	for _, env := range envs {
		_, present := os.LookupEnv(env)

		if present == false {
			fmt.Fprintf(GinkgoWriter, "Environment Variable %s must be set", envs)
			os.Exit(1)
		}
	}
}
