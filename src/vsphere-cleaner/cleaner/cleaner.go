package cleaner

import "vsphere-cleaner/parser"

type Cleaner struct {
	lockPath   string
	yamlParser parser.Parser
}

func NewCleaner(lockPath string, yamlParser parser.Parser) Cleaner {
	return Cleaner{lockPath: lockPath, yamlParser: yamlParser}
}

func (c Cleaner) Clean() {
	c.yamlParser.Parse(c.lockPath)
	//read file
	//parse file

	//get ips
	//get vms
	//delete each vm
}
