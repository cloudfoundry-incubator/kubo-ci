package cleaner

import "vsphere-cleaner/parser"

type Cleaner struct {
	lockPath   string
	yamlParser parser.Parser
}

func NewCleaner(lockPath string, yamlParser parser.Parser) Cleaner {
	return Cleaner{lockPath: lockPath, yamlParser: yamlParser}
}

func (c Cleaner) Clean() error {
	_, err := c.yamlParser.Parse(c.lockPath)

	// err = vspherDeleteVM(config, "ip")

	return err
	//vsphere.getVMs(config.InternalCIDR) = 128
	//possibleVMs := config.InternalCIDR - config.ReservedIPs //30
	//get ips
	//get vms
	//delete each vm
}
