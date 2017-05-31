package parser

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
		ips = difference(ips, reserved)
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

func (c Config) DirectorIP() string {
	return c.InternalIP
}

func difference(slice1 []string, slice2 []string) []string {
	result := []string{}
	for _, e := range slice1 {
		if !contains(slice2, e) {
			result = append(result, e)
		}
	}
	return result
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
