package flag_generator

import (
	"github.com/spf13/pflag"
)

type K8sFlags interface {
	AddFlags(fs *pflag.FlagSet)
}
