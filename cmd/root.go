package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "rv",
	Short:   "Robotics Verifier CLI",
	Long:    "Robotics Verifier CLI – early-stage electrical architecture checks for robotics projects.",
	Aliases: []string{"robotics-verifier-cli"},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// subcommands register themselves in init()
}
