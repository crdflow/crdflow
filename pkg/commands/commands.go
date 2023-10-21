package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "crdflow",
		Short: "Generate gRPC server from Kubernetes CRD",
	}
)

// Execute CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	addCommands()
}

func addCommands() {
	rootCmd.AddCommand(InitCommand())
}
