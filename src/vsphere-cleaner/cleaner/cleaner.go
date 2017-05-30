package cleaner

import (
	"vsphere-cleaner/parser"
	"vsphere-cleaner/vsphere"
)

type Cleaner struct {
	lockPath   string
	yamlParser parser.Parser
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
	return err
	//vsphere.getVMs(config.InternalCIDR) = 128
	//possibleVMs := config.InternalCIDR - config.ReservedIPs //30
	//get ips
	//get vms
	//delete each vm
}
