package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Config struct {
	Iaas             string           `json:"iaas"`
	TimeoutScale     float64          `json:"timeout_scale"`
	AWS              AWS              `json:"aws"`
	Bosh             Bosh             `json:"bosh"`
	Turbulence       Turbulence       `json:"turbulence"`
	TurbulenceTests  TurbulenceTests  `json:"turbulence_tests"`
	Kubernetes       Kubernetes       `json:"kubernetes"`
	CFCR             CFCR             `json:"cfcr"`
	UpgradeTests     UpgradeTests     `json:"upgrade_tests"`
}

type AWS struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	Region          string `json:"region"`
	IngressGroupID  string `json:"ingress_group_id"`
}

type Bosh struct {
	Environment  string `json:"environment"`
	CaCert       string `json:"ca_cert"`
	Client       string `json:"client"`
	ClientSecret string `json:"client_secret"`
	Deployment   string `json:"deployment"`
}

type Turbulence struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	CaCert   string `json:"ca_cert"`
}

type TurbulenceTests struct {
	IncludeWorkerDrain        bool `json:"include_worker_drain"`
	IncludeWorkerFailure      bool `json:"include_worker_failure"`
	IncludeMasterFailure      bool `json:"include_master_failure"`
	IncludePersistenceFailure bool `json:"include_persistence_failure"`
	IsMultiAZ                 bool `json:"is_multiaz"`
}

type UpgradeTests struct {
	IncludeMultiAZ bool `json:"include_multiaz"`
}

type Kubernetes struct {
	TLSCert             string `json:"tls_cert"`
	TLSPrivateKey       string `json:"tls_private_key"`
	KubernetesServiceIP string `json:"kubernetes_service_ip"`
	ClusterIPRange      string `json:"cluster_ip_range"`
	KubeDNSIP           string `json:"kube_dns_ip"`
	PodIPRange          string `json:"pod_ip_range"`
}

type CFCR struct {
	DeploymentPath           string `json:"deployment_path"`
	UpgradeToStemcellVersion string `json:"upgrade_to_stemcell_version"`
}

func InitConfig() (*Config, error) {
	var config Config
	var configPath = os.Getenv("CONFIG")

	if configPath == "" {
		return nil, errors.New("config path must be defined using CONFIG variable")
	}

	configJSON, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(configJSON, &config)
	if err != nil {
		return nil, err
	}

	// Do not allow zero for timeout scale as it would fail all the time.
	if config.TimeoutScale == 0 {
		config.TimeoutScale = 1
	}

	return &config, nil
}
