package main

import (
	"fmt"
	"os"
	"vsphere-cleaner/cleaner"
	"vsphere-cleaner/parser"
	"vsphere-cleaner/vsphere"
)

func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}
	err := cleaner.NewCleaner(os.Args[1], parser.NewParser(), vsphere.NewClient).Clean()
	if err != nil {
		fmt.Printf("Failed to clean environment %#v \n", err)
		panic(err)
	}
	fmt.Println("Environment cleaned")
}
