package main

import (
	cmd "github.com/jdamata/k8s-events/cmd"
)

var version = "dev"

func main() {
	cmd.Execute(version)
}
