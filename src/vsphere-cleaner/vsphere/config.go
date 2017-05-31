package vsphere

import (
	"errors"
	"net"
)

type Config struct {
	IP       string `yaml:"vcenter_ip"`
	User     string `yaml:"vcenter_user"`
	Password string `yaml:"vcenter_password"`

	InternalCIDR string   `yaml:"internal_cidr"`
	InternalIP   IP       `yaml:"internal_ip"`
	ReservedIPs  []string `yaml:"reserved_ips"`
}

func (config Config) UsedIPs() []string {
	return []string{}
}

type IP string

func (ip *IP) UnmarshalYAML(u func(interface{}) error) error {
	var s string
	err := u(&s)
	if err != nil {
		return err
	}
	parsedIP := net.ParseIP(s)
	if parsedIP == nil {
		return errors.New("Invalid IP")
	}
	*ip = IP(s)
	return nil
}
