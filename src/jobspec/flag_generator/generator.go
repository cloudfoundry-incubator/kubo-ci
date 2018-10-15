package flag_generator

import (
	goflag "flag"
	"io/ioutil"
	"os"

	"github.com/spf13/pflag"
	yaml "gopkg.in/yaml.v2"
)

type JobSpec struct {
	Name       string                   `yaml:"name"`
	Templates  map[string]string        `yaml:"templates"`
	Packages   []string                 `yaml:"packages"`
	Properties map[string]Property      `yaml:"properties"`
	Consumes   []map[string]interface{} `yaml:"consumes"`
	Provides   []map[string]interface{} `yaml:"provides"`
}

type Property struct {
	Properties  map[string]Property `yaml:"properties,omitempty,inline"`
	Description string              `yaml:"description,omitempty"`
	Default     interface{}         `yaml:"default,omitempty"`
}

func Contains(key string, arr []string) bool {
	for _, v := range arr {
		if v == key {
			return true
		}
	}
	return false
}

func ReadSpecFile(specPath string) (*JobSpec, error) {
	file, err := os.OpenFile(specPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return ReadSpec(contents)
}

func ReadSpec(spec []byte) (*JobSpec, error) {
	jobSpec := &JobSpec{}
	err := yaml.Unmarshal(spec, jobSpec)
	return jobSpec, err
}

func WriteSpecFile(specPath string, jobSpec *JobSpec) error {
	jobSpecBytes, err := WriteSpec(jobSpec)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(specPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	err = file.Truncate(0)
	if err != nil {
		return err
	}
	file.WriteAt(jobSpecBytes, 0)
	return nil
}

func WriteSpec(jobSpec *JobSpec) ([]byte, error) {
	return yaml.Marshal(jobSpec)
}

func GenerateArgsFromFlags(apiserverFlags K8sFlags, blacklistedFlags []string) Property {
	newProperties := Property{Properties: map[string]Property{}}
	flags := pflag.NewFlagSet("all", pflag.ContinueOnError)
	flags.AddGoFlagSet(goflag.CommandLine)

	apiserverFlags.AddFlags(flags)
	flags.VisitAll(func(f *pflag.Flag) {
		if !Contains(f.Name, blacklistedFlags) {
			newProperties.Properties[f.Name] = Property{Description: f.Usage}
		}
	})

	return newProperties
}
