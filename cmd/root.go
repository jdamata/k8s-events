package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "k8s-events",
		Short: "Logs kubernetes events",
		Long:  "Logs kubernetes events",
	}
)

// Execute executes the root command.
func Execute(version string) error {
	rootCmd.Version = version
	return rootCmd.Execute()
}

func logging() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
}
