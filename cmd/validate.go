package cmd

import (
	"fmt"
	"os"

	"github.com/badimirzai/robotics-verifier-cli/internal/model"
	"github.com/badimirzai/robotics-verifier-cli/internal/output"
	"github.com/badimirzai/robotics-verifier-cli/internal/parts"
	"github.com/badimirzai/robotics-verifier-cli/internal/resolve"
	"github.com/badimirzai/robotics-verifier-cli/internal/validate"
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
		var doc yaml.Node
		if err := yaml.Unmarshal(b, &doc); err != nil {
			return fmt.Errorf("parse yaml: %w", err)
		}
		if err := doc.Decode(&raw); err != nil {
			return fmt.Errorf("decode yaml: %w", err)
		}

		store := parts.NewStore("parts")
		resolved, err := resolve.ResolveAll(raw, store)
		if err != nil {
			return fmt.Errorf("resolve spec with parts: %w", err)
		}

		locs := buildLocationMap(path, &doc)
		rep := validate.RunAll(resolved, locs)
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

func buildLocationMap(path string, doc *yaml.Node) map[string]validate.Location {
	locs := make(map[string]validate.Location)

	var walk func(n *yaml.Node, prefix string)
	walk = func(n *yaml.Node, prefix string) {
		switch n.Kind {
		case yaml.DocumentNode:
			for _, child := range n.Content {
				walk(child, prefix)
			}
		case yaml.MappingNode:
			for i := 0; i+1 < len(n.Content); i += 2 {
				key := n.Content[i]
				val := n.Content[i+1]
				if key.Kind != yaml.ScalarNode {
					continue
				}
				next := key.Value
				if prefix != "" {
					next = prefix + "." + key.Value
				}
				locs[next] = validate.Location{File: path, Line: key.Line}
				if val.Kind == yaml.ScalarNode {
					locs[next] = validate.Location{File: path, Line: val.Line}
				}
				walk(val, next)
			}
		case yaml.SequenceNode:
			for i, item := range n.Content {
				next := fmt.Sprintf("%s[%d]", prefix, i)
				locs[next] = validate.Location{File: path, Line: item.Line}
				if item.Kind == yaml.ScalarNode {
					locs[next] = validate.Location{File: path, Line: item.Line}
				}
				walk(item, next)
			}
		}
	}

	if doc != nil {
		walk(doc, "")
	}

	return locs
}
