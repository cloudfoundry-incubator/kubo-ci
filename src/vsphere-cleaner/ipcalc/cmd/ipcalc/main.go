package main

import (
	"fmt"
	"os"

	"vsphere-cleaner/parser"
)

func main() {
	config, _ := parser.NewParser().Parse(os.Args[1])
	ips, _ := config.UsedIPs()
	fmt.Println(ips)
}