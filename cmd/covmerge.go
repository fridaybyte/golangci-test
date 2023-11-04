package cmd

import (
	"fmt"
	"os"

	"github.com/shabbyrobe/gocovmerge"
	"github.com/spf13/cobra"
	"golang.org/x/tools/cover"
)

// covmergeCmd represents the covmerge command
var covmergeCmd = &cobra.Command{
	Use:   "covmerge",
	Short: "Merge multiple coverage profiles into one.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("covmerge called")

		var merged []*cover.Profile

		for _, file := range args {
			profiles, err := cover.ParseProfiles(file)
			if err != nil {
				panic(fmt.Errorf("failed to parse profiles: %w", err))
			}
			for _, p := range profiles {
				merged = gocovmerge.AddProfile(merged, p)
			}
		}

		err := gocovmerge.DumpProfiles(merged, os.Stdout)
		if err != nil {
			panic(fmt.Errorf("failed to dump profiles: %w", err))
		}
	},
}

func init() {
	rootCmd.AddCommand(covmergeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// covmergeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// covmergeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
