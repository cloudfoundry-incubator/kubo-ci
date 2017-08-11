package test_helpers

import (
	"bytes"
	"io"
	"strings"

	"github.com/cloudfoundry/bosh-cli/cmd"
	"github.com/cloudfoundry/bosh-cli/ui"
	"github.com/cloudfoundry/bosh-utils/logger"

	"github.com/onsi/ginkgo"
)

func BoshCmdFactory(out, err *bytes.Buffer) cmd.Factory {
	l := logger.NewLogger(logger.LevelNone)
	output := io.MultiWriter(ginkgo.GinkgoWriter, out)
	errors := io.MultiWriter(ginkgo.GinkgoWriter, err)
	boshUI := ui.NewWriterUI(output, errors, l)
	someUI := ui.NewPaddingUI(boshUI)
	confUI := ui.NewWrappingConfUI(someUI, l)
	confUI.EnableNonInteractive()
	return cmd.NewFactory(cmd.NewBasicDeps(confUI, l))
}

func RunningVmList(vmType string) func() ([]string, error) {
	return func() ([]string, error) {
		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})

		cmdFactory := BoshCmdFactory(stdout, stderr)
		boshCommand, err := cmdFactory.New([]string{"vms"})
		if err != nil {
			return nil, err
		}
		boshCommand.Execute()
		vmTable := stdout.String()

		return FilterArrayOfStrings(strings.Split(vmTable, "\n"), func(line string) bool {
			return strings.Contains(line, vmType) && strings.Contains(line, "running")
		}), nil

	}
}
