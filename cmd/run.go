package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"golangci-test/splitrunner"
	"golangci-test/splitter"
)

var splitFile string
var jsonOut string
var machineInstance int
var cfg struct {
	runFS   *pflag.FlagSet
	splitFS *pflag.FlagSet
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run tests with ability to run specific split",
	RunE: func(cmd *cobra.Command, args []string) error {
		splits, err := splitter.LoadSplits(splitFile)
		if err != nil {
			return err
		}
		if machineInstance >= len(splits) {
			fmt.Printf(`Warning: There are %d groups of packages to run, but passed
machine instance index is %d. Please consider reducing
number of runners or increase number of package groups.

Nothing left to test. Exiting quietly.
`, len(splits), machineInstance)
			return nil
		}
		fmt.Printf("Running group of packages number %d\n", machineInstance)
		return splitrunner.RunSplit(splits[machineInstance], jsonOut)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	fs := runCmd.Flags()
	cfg.runFS = fs
	fs.StringVarP(&splitFile, "split-file", "s", "",
		"Path to file containing package groups")
	fs.StringVar(&jsonOut, "json", "",
		"Converts output to JSON suitable for further processing. E.g. for generating new splits.")
	fs.IntVarP(&machineInstance, "index", "i", 0,
		"Index of package group to run")

	err := runCmd.MarkFlagRequired("index")
	if err != nil {
		panic(err)
	}
}
