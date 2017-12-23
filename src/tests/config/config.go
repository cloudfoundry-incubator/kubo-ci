package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Config struct {
	Bosh       Bosh       `json:"bosh"`
	Turbulence Turbulence `json:"turbulence"`
	Cf         Cf         `json:"cf"`
	Kubernetes Kubernetes `json:"kubernetes"`
}

type Bosh struct {
	Iaas         string `json:"iaas"`
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

type Cf struct {
	AppsDomain string `json:"apps_domain"`
}

type Kubernetes struct {
	AuthorizationMode string `json:"authorization_mode"`
	MasterHost        string `json:"master_host"`
	MasterPort        int    `json:"master_port"`
	PathToKubeConfig  string `json:"path_to_kube_config"`
	TLSCert           string `json:"tls_cert"`
	TLSPrivateKey     string `json:"tls_private_key"`
}

func InitConfig() (*Config, error) {
	var config Config
	var configPath = os.Getenv("CONFIG")

	if configPath == "" {
		return nil, errors.New("config path must be defined")
	}

	configJSON, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(configJSON, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
