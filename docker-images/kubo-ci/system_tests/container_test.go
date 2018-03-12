package systemTests

import (
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Docker Image", func() {

	commands := map[string][]string{
		"bosh-cli": []string{"-v"},
		"bosh":     []string{"-v"},
		"cf":       []string{"version"},
		"credhub":  []string{"--version"},
		"gcloud":   []string{"--version"},
		"ginkgo":   []string{"help"},
		"make":     []string{"-v"},
		"ruby":     []string{"-v"},
		"kubectl":  []string{"help"},
		"go":       []string{"doc", "cmd/vet"},
		"aws":      []string{"--version"},
		"govc":     []string{"version"},
		"ipcalc":   []string{},
		"java":     []string{"-version"},
		"om":       []string{"--version"},
		"vara":     []string{"help"},
		"pry":      []string{"--version"},
		"rspec":    []string{"--version"},
		"rake":     []string{"--version"},
		"zip":      []string{"--version"},
		"jq":       []string{"--version"},
		"which":    []string{"sshuttle"},
		"spruce":   []string{"--version"},
		"vegeta":   []string{"-version"},
		"haproxy":  []string{"-v"},
	}

	for executable, args := range commands {
		executable := executable
		args := args

		It(fmt.Sprintf("has %v installed", executable), func() {
			command := exec.Command(executable, args...)
			session, err := Start(command, GinkgoWriter, GinkgoWriter)

			Expect(err).ToNot(HaveOccurred())
			Eventually(session, "5s").Should(Exit(0))
		})
	}
})
