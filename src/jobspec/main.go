package main

import (
	"fmt"

	"github.com/spf13/pflag"
	"k8s.io/kubernetes/cmd/kube-apiserver/app/options"

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
		"etcd-servers",
		"apiserver-count",
		"cloud-provider",
		"cloud-config",
	}
	k := kflags{}
	specPath := os.Args[1]
	file, _ := os.OpenFile(specPath, os.O_RDWR|os.O_CREATE, 0644)
	defer file.Close()
	jobSpec := &JobSpec{}
	c, _ := ioutil.ReadAll(file)
	yaml.Unmarshal(c, jobSpec)
	flags := pflag.NewFlagSet("all", pflag.ContinueOnError)
	apiserverFlags := options.NewServerRunOptions()
	apiserverFlags.AddFlags(flags)
	flags.VisitAll(func(f *pflag.Flag) {
		if Contains(f.Name, blacklistedFlags) {
			delete(jobSpec.Properties, "args."+f.Name)
		} else {
			k = append(k, "args."+f.Name)
			jobSpec.Properties["args."+f.Name] = Property{Description: f.Usage}
		}
	})
	jobSpecBytes, _ := yaml.Marshal(jobSpec)
	file.WriteAt(jobSpecBytes, 0)

	for s := range jobSpec.Properties {
		if !Contains(s, k) {
			fmt.Fprintf(os.Stderr, "We have the following flag, %s, in our spec that was not provided by kubernetes (or was blacklisted by this tool)\n", s)
		}
	}
}
