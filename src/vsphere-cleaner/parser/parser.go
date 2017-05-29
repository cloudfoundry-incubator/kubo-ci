package parser

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

//go:generate counterfeiter ./ Parser
type Parser interface {
	Parse(string) (VMWareConfig, error)
}

type parserImpl struct{}

func NewParser() Parser {
	return parserImpl{}
}

type VMWareConfig struct {
	IP        string `yaml:"vcenter_ip"`
	User      string `yaml:"vcenter_user"`
	Password  string `yaml:"vcenter_password"`
	DC        string `yaml:"vcenter_dc"`
	Cluster   string `yaml:"vcenter_cluster"`
	DS        string `yaml:"vcenter_ds"`
	RP        string `yaml:"vcenter_rp"`
	VMs       string `yaml:"vcenter_vms"`
	Templates string `yaml:"vcenter_templates"`
	Disks     string `yaml:"vcenter_disks"`

	NetworkName  string   `yaml:"network_name"`
	InternalCIDR string   `yaml:"internal_cidr"`
	InternalGW   string   `yaml:"internal_gw"`
	InternalIP   string   `yaml:"internal_ip"`
	ReservedIPs  []string `yaml:"reserved_ips"`
}

func (parserImpl) Parse(lockPath string) (VMWareConfig, error) {
	config := VMWareConfig{}
	dat, err := ioutil.ReadFile(lockPath)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(dat, &config)
	return config, err
}
