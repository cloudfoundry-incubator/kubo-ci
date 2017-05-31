package cleaner

import (
	"vsphere-cleaner/parser"
	"vsphere-cleaner/vsphere"
)

type Cleaner struct {
	lockPath      string
	yamlParser    parser.Parser
	vSphereClient vsphere.Client
}

func NewCleaner(lockPath string, yamlParser parser.Parser, client vsphere.Client) Cleaner {
	return Cleaner{lockPath: lockPath, yamlParser: yamlParser, vSphereClient: client}
}

func (c Cleaner) Clean() error {
	config, err := c.yamlParser.Parse(c.lockPath)
	if err != nil {
		return err
	}
	err = c.vSphereClient.DeleteVM(config.InternalIP)
	if err != nil {
		return err
	}

	ips, _ := config.UsedIPs()
	for _, ip := range ips {
		err = c.vSphereClient.DeleteVM(ip)
		if err != nil {
			return err
		}
	}
	return err
}
