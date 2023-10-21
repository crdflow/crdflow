package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AddStringFlag to cobra.Command.
func AddStringFlag(cmd *cobra.Command, name, value, usage string, required bool) {
	cmd.Flags().String(name, value, usage)
	_ = viper.BindPFlag(name, cmd.Flags().Lookup(name))

	if required {
		_ = cmd.MarkFlagRequired(name)

		u := cmd.Flag(name).Usage
		cmd.Flag(name).Usage = fmt.Sprintf("%s (required)", u)
	}
}
