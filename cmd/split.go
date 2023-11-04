package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"golangci-test/splitter"
)

var splitType string
var outputFilePath string
var numberOfSplits int

// splitCmd represents the split command
var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "Split packages into multiple groups of packages",
	Long: `
Split packages into multiple groups of packages. 
Afterwards, each group of packages can be tested 
in parallel on a separate machine.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("split called")
		splits, err := splitter.GenerateSplits("./...", numberOfSplits)
		if err != nil {
			return err
		}
		if outputFilePath != "" {
			err = splitter.StoreSplits(splits, outputFilePath)
			if err != nil {
				return err
			}
		} else {
			for i, split := range splits {
				fmt.Printf("Split %d:\n", i+1)
				for _, pkg := range split {
					fmt.Printf("- %s\n", pkg)
				}
				fmt.Println()
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)

	fs := splitCmd.Flags()
	cfg.splitFS = fs
	fs.StringVarP(&splitType, "type", "t", "time",
		`Choose how to split packages into groups. Possible values are:
		time - running time of each package group trends to be equal
		count - number of packages in each group trends to be equal
`)
	fs.StringVarP(&outputFilePath, "out", "o", "",
		`Convert split output to JSON suitable for further processing. 
This file can be used as input for "run" command.`)
	fs.IntVarP(&numberOfSplits, "groups-count", "g", 1, "Number of package groups to generate")
	err := splitCmd.MarkFlagRequired("groups-count")
	if err != nil {
		panic(err)
	}
}
