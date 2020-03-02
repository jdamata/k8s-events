package main

import (
	cmd "github.com/jdamata/k8s-event/cmd"
)

var version = "dev"

func main() {
	cmd.Execute(version)
}
