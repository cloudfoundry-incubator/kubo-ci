package main

import (
	"fmt"
	"os"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

type VarsFile struct {
	EtcdReleaseVersion       string `yaml:"etcd_release_version"`
	LatestEtcdReleaseVersion string `yaml:"latest_etcd_release_version"`
	StemcellVersion          string `yaml:"stemcell_version"`
	BOSHEnvironment          string `yaml:"bosh_environment"`
	BOSHClient               string `yaml:"bosh_client"`
	BOSHClientSecret         string `yaml:"bosh_client_secret"`
	BOSHDirectorCACert       string `yaml:"bosh_director_ca_cert"`
	EnableTurbulenceTests    bool   `yaml:"enable_turbulence_tests"`
	ParallelNodes            int    `yaml:"parallel_nodes"`
}

func main() {
	output, err := Generate()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, string(output))
}

func Generate() ([]byte, error) {
	varsFile := VarsFile{
		EtcdReleaseVersion:       os.Getenv("ETCD_RELEASE_VERSION"),
		LatestEtcdReleaseVersion: os.Getenv("LATEST_ETCD_RELEASE_VERSION"),
		StemcellVersion:          os.Getenv("STEMCELL_VERSION"),
		BOSHEnvironment:          os.Getenv("BOSH_ENVIRONMENT"),
		BOSHClient:               os.Getenv("BOSH_CLIENT"),
		BOSHClientSecret:         os.Getenv("BOSH_CLIENT_SECRET"),
		BOSHDirectorCACert:       os.Getenv("BOSH_CA_CERT"),
		EnableTurbulenceTests:    (os.Getenv("ENABLE_TURBULENCE_TESTS") == "true"),
	}

	parallelNodes, err := strconv.Atoi(os.Getenv("PARALLEL_NODES"))
	if err != nil {
		return nil, err
	}

	varsFile.ParallelNodes = parallelNodes

	contents, err := yaml.Marshal(varsFile)
	if err != nil {
		// not tested
		return nil, err
	}

	return contents, nil
}
