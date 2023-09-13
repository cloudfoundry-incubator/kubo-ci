package test_helpers

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"path/filepath"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var (
	PKS_CLI_PATH_ENV string
)

type PksCliRunner struct {
	CliPath          string
	TimeoutInSeconds float64
}

func SetupPksCli() (*PksCliRunner, error) {
	currentPath, _ := os.Getwd()
	rootPath := filepath.Join(currentPath, "../../../../..")
	cliPathDir := rootPath + "/pks-cli/pks-*"
	matchingPaths, err := filepath.Glob(cliPathDir)
	if err != nil {
		return nil, err
	}
	pksCliPath := matchingPaths[0]

	pks_cli := PksCliRunner{
		CliPath:          pksCliPath,
		TimeoutInSeconds: 3600,
	}
	return &pks_cli, nil
}

func (pkscli PksCliRunner) CreateKubernetesProfile(k8sProfile string) (string, error) {
	args := []string{"create-k8s-profile", k8sProfile}
	return pkscli.executeCommand(args...)
}

func (pkscli PksCliRunner) UpdateClusterWithProfile(name string, profileName string) (string, error) {
	args := []string{"update-cluster", name, "--kubernetes-profile", profileName, "--non-interactive", "--wait"}
	return pkscli.executeCommand(args...)
}

func (pkscli PksCliRunner) executeCommand(args ...string) (string, error) {
	command := exec.Command(pkscli.CliPath, args...)
	fmt.Printf("Executing command : %s", command.Args)
	session, err := gexec.Start(command, ginkgo.GinkgoWriter, ginkgo.GinkgoWriter)
	if err != nil {
		return "", err
	}

	gomega.Eventually(session, pkscli.TimeoutInSeconds).Should(gexec.Exit())
	if session.ExitCode() != 0 {
		return string(session.Out.Contents()), errors.New(string(session.Err.Contents()))
	}

	output := session.Out.Contents()
	return string(output), nil
}
