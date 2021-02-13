package cmd

import (
	"fmt"
	"github.com/karuppiah7890/go-jsonschema-generator"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

func uncommentYAML(node *yaml.Node) *yaml.Node {
	switch node.Kind {
	case yaml.DocumentNode, yaml.SequenceNode:
		for i, child := range node.Content {
			node.Content[i] = uncommentYAML(child)
		}
		return node
	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i = i + 2 {
			key := node.Content[i]
			value := node.Content[i+1]

			if len(key.FootComment) > 0 {
				if len(value.Content) == 0 && (value.Kind == yaml.MappingNode || value.Kind == yaml.SequenceNode) {
					comment := strings.ReplaceAll(key.FootComment, "#", strings.Repeat(" ", key.Column-1))
					root := yaml.Node{}
					if err := yaml.Unmarshal([]byte(comment), &root); err == nil {
						if root.Content[0].Kind == value.Kind {
							node.Content[i].FootComment = ""
							node.Content[i+1].Content = append(value.Content, root.Content[0].Content...)
						}
					}
				}
			}

			if i > 0 && len(key.HeadComment) > 0 {
				prev := node.Content[i-1]
				if len(prev.Content) == 0 && (prev.Kind == yaml.MappingNode || prev.Kind == yaml.SequenceNode) {
					comment := strings.ReplaceAll(key.HeadComment, "#", strings.Repeat(" ", key.Column-1))
					root := yaml.Node{}
					if err := yaml.Unmarshal([]byte(comment), &root); err == nil {
						if root.Content[0].Kind == prev.Kind {
							node.Content[i].HeadComment = ""
							node.Content[i-1].Content = append(prev.Content, root.Content[0].Content...)
						}
					}
				}
			}

			node.Content[i+1] = uncommentYAML(node.Content[i+1])
		}

		return node
	default:
		return node
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "helm schema-gen <values-yaml-file>",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "Helm plugin to generate json schema for values yaml",
	Long: `Helm plugin to generate json schema for values yaml

Examples:
  $ helm schema-gen values.yaml    # generate schema json
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("pass one values yaml file")
		}
		if len(args) != 1 {
			return fmt.Errorf("schema can be generated only for one values yaml at once")
		}

		valuesFilePath := args[0]
		values := make(map[string]interface{})
		valuesFileData, err := ioutil.ReadFile(valuesFilePath)
		if err != nil {
			return fmt.Errorf("error when reading file '%s': %v", valuesFilePath, err)
		}
		root := yaml.Node{}
		err = yaml.Unmarshal(valuesFileData, &root)
		if err != nil {
			return fmt.Errorf("error when unmarshaling file '%s': %v", valuesFilePath, err)
		}
		if ok, _ := cmd.Flags().GetBool("yaml-comment"); ok {
			root = *uncommentYAML(&root)
		}
		err = root.Decode(&values)
		if err != nil {
			fmt.Println(err)
		}
		s := &jsonschema.Document{}
		s.ReadDeep(&values)
		fmt.Println(s)

		return nil
	},
}

// Execute executes the root command
func Execute() {
	rootCmd.Flags().Bool("yaml-comment", false, "Generate a schema with YAML format comments as YAML when it is an empty array or map")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
