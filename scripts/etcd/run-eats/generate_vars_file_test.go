package main_test

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/onsi/gomega/gexec"
	"github.com/pivotal-cf-experimental/gomegamatchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Generate vars.yml", func() {
	var variables map[string]string

	BeforeEach(func() {
		variables = map[string]string{
			"BOSH_ENVIRONMENT":            "some-bosh-target",
			"BOSH_CLIENT":                 "some-bosh-username",
			"BOSH_CLIENT_SECRET":          "some-bosh-password",
			"BOSH_CA_CERT":                "some-bosh-director-ca-cert",
			"ENABLE_TURBULENCE_TESTS":     "true",
			"PARALLEL_NODES":              "10",
			"ETCD_RELEASE_VERSION":        "some-etcd-release-version",
			"STEMCELL_VERSION":            "some-stemcell-version",
			"LATEST_ETCD_RELEASE_VERSION": "some-latest-etcd-release-version",
		}

		for name, value := range variables {
			variables[name] = os.Getenv(name)
			os.Setenv(name, value)
		}
	})

	AfterEach(func() {
		for name, value := range variables {
			os.Setenv(name, value)
		}
	})

	It("generates a vars.yml", func() {
		expectedVars, err := ioutil.ReadFile("fixtures/expected-vars.yml")
		Expect(err).NotTo(HaveOccurred())

		session, err := gexec.Start(exec.Command(pathToMain), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(gexec.Exit(0))
		Expect(session.Out.Contents()).To(gomegamatchers.MatchYAML(expectedVars))
	})

	Context("failure cases", func() {
		Context("when the PARALLEL_NODES env var is not an integer", func() {
			It("exits with a non-zero return value and outputs an error", func() {
				os.Setenv("PARALLEL_NODES", "foo")

				session, err := gexec.Start(exec.Command(pathToMain), GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				Expect(string(session.Err.Contents())).To(ContainSubstring("strconv.Atoi: parsing \"foo\": invalid syntax"))
			})
		})
	})
})
