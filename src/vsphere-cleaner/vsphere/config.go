package vsphere

import (
	"net/url"
)

//go:generate counterfeiter ./ Config
type Config interface {
	BuildUrl() *url.URL
	UsedIPs() ([]string, error)
	DirectorIP() string
}
