package main

import (
	"github.com/spf13/pflag"
	"k8s.io/kubernetes/cmd/kube-apiserver/app/options"

	goflag "flag"
	"os"

	"io/ioutil"

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

type kflags []string

func Contains(key string, arr []string) bool {
	for _, v := range arr {
		if v == key {
			return true
		}
	}
	return false
}

func main() {
	blacklistedFlags := []string{
		"apiserver-count",
		"cloud-provider",
		"cloud-config",
	}
	specPath := os.Args[1]
	file, _ := os.OpenFile(specPath, os.O_RDWR|os.O_CREATE, 0644)
	defer file.Close()
	jobSpec := &JobSpec{}
	c, _ := ioutil.ReadAll(file)
	yaml.Unmarshal(c, jobSpec)
	jobSpec.Properties["k8s-args"] = Property{Properties: map[string]Property{}}
	flags := pflag.NewFlagSet("all", pflag.ContinueOnError)
	flags.AddGoFlagSet(goflag.CommandLine)
	apiserverFlags := options.NewServerRunOptions()
	apiserverFlags.AddFlags(flags)
	flags.VisitAll(func(f *pflag.Flag) {
		if Contains(f.Name, blacklistedFlags) {
			delete(jobSpec.Properties["k8s-args"].Properties, f.Name)
		} else {
			jobSpec.Properties["k8s-args"].Properties[f.Name] = Property{Description: f.Usage}
		}
	})
	jobSpecBytes, _ := yaml.Marshal(jobSpec)
	if err := file.Truncate(0); err != nil {
		panic(err)
	}
	file.WriteAt(jobSpecBytes, 0)
}
