package parser

import (
	"io/ioutil"
	"vsphere-cleaner/vsphere"

	yaml "gopkg.in/yaml.v2"
)

//go:generate counterfeiter ./ Parser
type Parser interface {
	Parse(string) (vsphere.Config, error)
}

type parserImpl struct{}

func NewParser() Parser {
	return parserImpl{}
}

func (parserImpl) Parse(lockPath string) (vsphere.Config, error) {
	config := vsphere.Config{}
	dat, err := ioutil.ReadFile(lockPath)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(dat, &config)
	return config, err
}
