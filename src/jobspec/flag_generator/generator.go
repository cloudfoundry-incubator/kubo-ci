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

func ReadExistingSpec(specPath string) *JobSpec {
	file, _ := os.OpenFile(specPath, os.O_RDWR|os.O_CREATE, 0644)
	defer file.Close()
	jobSpec := &JobSpec{}
	c, _ := ioutil.ReadAll(file)
	yaml.Unmarshal(c, jobSpec)
	return jobSpec
}

func WriteNewSpec(specPath string, jobSpec *JobSpec) {
	jobSpecBytes, _ := yaml.Marshal(jobSpec)
	file, _ := os.OpenFile(specPath, os.O_RDWR|os.O_CREATE, 0644)
	if err := file.Truncate(0); err != nil {
		panic(err)
	}
	file.WriteAt(jobSpecBytes, 0)
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
