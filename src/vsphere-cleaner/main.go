package main

import "os"

func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}
}
