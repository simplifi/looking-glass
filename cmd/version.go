package cmd

import (
	"fmt"
	"github.com/simplifi/looking-glass/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of looking-glass",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("looking-glass Version %s", version.Version)
	},
}