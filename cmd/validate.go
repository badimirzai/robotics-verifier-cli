package cmd

import (
	"fmt"
	"os"

	"github.com/badimirzai/architon-cli/internal/model"
	"github.com/badimirzai/architon-cli/internal/output"
	"github.com/badimirzai/architon-cli/internal/parts"
	"github.com/badimirzai/architon-cli/internal/resolve"
	"github.com/badimirzai/architon-cli/internal/validate"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var checkCmd = &cobra.Command{
	Use:     "check <spec.yaml>",
	Aliases: []string{"validate"},
	Args:    cobra.MaximumNArgs(1),
	Short:   "Validate a robot spec against deterministic electrical rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("file")
		if path == "" && len(args) > 0 {
			path = args[0]
		}
		if path == "" {
			return fmt.Errorf("missing spec file (arg or -f/--file)")
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read spec: %w", err)
		}

		var raw model.RobotSpec
		if err := yaml.Unmarshal(b, &raw); err != nil {
			return fmt.Errorf("parse yaml: %w", err)
		}

		store := parts.NewStore("parts")
		resolved, err := resolve.ResolveAll(raw, store)
		if err != nil {
			return fmt.Errorf("resolve spec with parts: %w", err)
		}

		rep := validate.RunAll(resolved)
		fmt.Println(output.RenderReport(rep))
		if rep.HasErrors() {
			os.Exit(2) // deterministic non-zero for CI
		}
		return nil
	},
}

func init() {
	checkCmd.Flags().StringP("file", "f", "", "Path to YAML spec")
	rootCmd.AddCommand(checkCmd)
}
