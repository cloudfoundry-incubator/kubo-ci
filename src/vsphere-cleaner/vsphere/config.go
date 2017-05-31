package vsphere

import (
	"net/url"
	"vsphere-cleaner/ipcalc"
)

type Config struct {
	IP       string `yaml:"vcenter_ip"`
	User     string `yaml:"vcenter_user"`
	Password string `yaml:"vcenter_password"`

	InternalCIDR string   `yaml:"internal_cidr"`
	InternalIP   string   `yaml:"internal_ip"`
	ReservedIPs  []string `yaml:"reserved_ips"`
}

type IConfig interface {
	BuildUrl() *url.URL
	UsedIPs() ([]string, error)
}

func (c Config) UsedIPs() ([]string, error) {
	ips, err := ipcalc.GetIPsFromCIDR(c.InternalCIDR)
	if err != nil {
		return []string{}, err
	}
	for _, reservedRange := range c.ReservedIPs {
		reserved, err := ipcalc.GetIPsFromRange(reservedRange)
		if err != nil {
			return []string{}, err
		}
		ips = ipcalc.Difference(ips, reserved)
	}
	return ips, nil
}

func (c Config) BuildUrl() *url.URL {
	parsedUrl := url.URL{
		Scheme: "https",
		Host:   c.IP,
		Path:   "/sdk",
		User:   url.UserPassword(c.User, c.Password),
	}
	return &parsedUrl
}
